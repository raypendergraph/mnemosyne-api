package neo4j

import (
	"bytes"
	"fmt"
	sys "mnemosyne-api/system"
	"text/template"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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

type statementParams = map[string]any
type TransactionArguments[T any] struct {
	Context    sys.ServiceContext
	Statement  string
	Args       statementParams
	Driver     neo4j.DriverWithContext
	MapResults resultMapping[T]
}

func safeCloseSession(session neo4j.SessionWithContext, ctx sys.ServiceContext) {
	err := session.Close(ctx.GetPortContext())
	if err != nil {
		ctx.GetLogger().LogWarning("problem closing neo4j session")
	}
}

func writeWithTransaction[T any](ta TransactionArguments[T]) (out sys.Result[T]) {
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
		return identifyError[T](ctx, err)
	}

	return sys.Result[T]{
		Value: result,
	}
}

func identifyError[T any](ctx sys.ServiceContext, e error) sys.Result[T] {
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

func statement(t *template.Template, values any) (string, error) {
	var buff bytes.Buffer
	if err := t.Execute(&buff, values); err != nil {
		return "", err
	}
	return buff.String(), nil

}
