package neo4j

import (
	"bytes"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	sys "mnemosyne-api/system"
	"strings"
	"time"
)

type parameters map[string]any

func (r parameters) CombinedWith(other ...parameters) parameters {
	for _, o := range other {
		for k, v := range o {
			r[k] = v
		}
	}
	return r
}

type Neo4JAdapter struct {
	driver neo4j.DriverWithContext
	errors sys.ErrorCatalog
}

func NewRepoBase(config sys.Neo4JConfiguration, errs sys.ErrorCatalog) (Neo4JAdapter, sys.Error) {
	driver, err := neo4j.NewDriverWithContext("neo4j://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		return Neo4JAdapter{}, errs.NewError(sys.ErrNeo4JConfigurationFailure, nil, err)
	}
	return Neo4JAdapter{
		driver: driver,
		errors: errs,
	}, nil
}

type ResultMapping[T any] func(result neo4j.ResultWithContext) (T, error)
type TransactionArguments[T any] struct {
	Context    sys.ServiceContext
	Statement  string
	Args       parameters
	Driver     neo4j.DriverWithContext
	MapResults ResultMapping[T]
}

func safeCloseSession(session neo4j.SessionWithContext, ctx sys.ServiceContext) {
	err := session.Close(ctx.GetPortContext())
	if err != nil {
		ctx.GetLogger().LogWarning("problem closing neo4j session")
	}
}

func WriteWithTransaction[T any](ta TransactionArguments[T]) (out sys.Result[T]) {
	ctx := ta.Context
	logger := ta.Context.GetLogger()
	logger.LogDebug("enter")
	defer logger.LogDebug("exit")
	session := ta.Driver.NewSession(ctx.GetPortContext(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer safeCloseSession(session, ctx)
	result, err := neo4j.ExecuteWrite[T](ctx.GetPortContext(), session, func(transaction neo4j.ManagedTransaction) (out T, err error) {
		var result neo4j.ResultWithContext
		if result, err = transaction.Run(ctx.GetPortContext(), ta.Statement, ta.Args); err != nil {
			return
		}

		return ta.MapResults(result)
	})
	if err != nil {
		return reidentifyError[T](ctx, err)
	}

	return sys.Result[T]{
		Value: result,
	}
}

func reidentifyError[T any](ctx sys.ServiceContext, e error) sys.Result[T] {
	if e, ok := e.(sys.Error); ok {
		return sys.Result[T]{Error: e.NextFrame()}
	}
	return sys.Result[T]{Error: ctx.NewError(sys.ErrNeo4JTransactedWriteFailure, e)}
}

func bindUUIDProperty(key string, node neo4j.Node, binding *sys.UUID) error {
	stringValue, err := neo4j.GetProperty[string](node, key)
	if err != nil {
		return err
	}
	uuidValue, err := sys.NewUUIDFromString(stringValue)
	if err != nil {
		return err
	}
	*binding = uuidValue
	return nil
}

func bindDateTimeProperty(key string, node neo4j.Node, binding *time.Time) error {
	rawValue, found := node.GetProperties()[key]
	if !found {
		return fmt.Errorf("%s is a mandatory property but was not found on the node", key)
	}
	value, ok := rawValue.(time.Time)
	if !ok {
		return fmt.Errorf("expected value to have type neo4j.Time but found type %T", rawValue)
	}
	*binding = value
	return nil
}

func bindNeo4JProperty[T neo4j.PropertyValue](key string, node neo4j.Node, binding *T) error {
	t, err := neo4j.GetProperty[T](node, key)
	if err != nil {
		return err
	}
	*binding = t
	return nil
}

func buildMaps(scope string, parameters parameters) (string, parameters) {
	if scope != "" {
		scope = fmt.Sprintf("%s_", scope)
	}
	parameterMap := make(map[string]any, len(parameters))
	pairStrings := make([]string, len(parameters))
	i := 0
	for field, value := range parameters {
		scopedField := fmt.Sprintf("%s%s", scope, field)
		pairString := fmt.Sprintf("%s: $%s", field, scopedField)
		parameterMap[scopedField] = value
		pairStrings[i] = pairString
		i += 1
	}
	mapString := fmt.Sprintf("{%s}", strings.Join(pairStrings, ", "))
	return mapString, parameterMap
}

func statement(t *template.Template, values any) (string, error) {
	var buff bytes.Buffer
	if err := t.Execute(&buff, values); err != nil {
		return "", err
	}
	return buff.String(), nil

}

type node struct {
	spec
	stringValue string
}

func (n node) String() string {
	return n.stringValue
}

type spec struct {
	Variable   string
	Parameters parameters
	Labels     []string
}

func (r spec) Build() node {
	const labelSeparator = ":"
	propertiesString, scopedParameters := buildMaps(r.Variable, r.Parameters)
	labelsString := ""
	if len(r.Labels) > 0 {
		labelsString = labelSeparator + strings.Join(r.Labels, labelSeparator)
	}
	return node{
		spec: spec{
			Variable:   r.Variable,
			Parameters: scopedParameters,
			Labels:     r.Labels,
		},
		stringValue: fmt.Sprintf("%s%s %s", r.Variable, labelsString, propertiesString),
	}
}
