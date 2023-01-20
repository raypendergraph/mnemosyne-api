package templates

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
)

const (
	GenericCreate = "generic_create"
)

var funcMap = template.FuncMap{
	"JoinStrings": strings.Join,
}
var Cypher = template.Must(createTemplates(
	GenericCreate))

func createTemplates(names ...string) (*template.Template, error) {
	files := make([]string, len(names))
	for i, name := range names {
		files[i] = fmt.Sprintf("templates/%s", name)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return t.Funcs(funcMap), nil
}

type TemplateSpec struct {
	Parameters  Parameters
	stringValue string
}

func (r TemplateSpec) String() string {
	return r.stringValue
}

type Relationship struct {
	TemplateSpec
	IsDirected bool
}
type MatchChain struct {
	Origin   *TemplateSpec
	Relation *Relationship
	Target   *TemplateSpec
}
type CreateJournalEntryArgs struct {
	JournalMatch      MatchChain
	RelationshipLabel string
	JournalEntry      TemplateSpec
	JournalEntryVar   string
}

type GenericCreateArgs struct {
	Match               *MatchChain
	Relationship        *TemplateSpec
	Created             TemplateSpec
	MatchedIsOrigin     bool
	ReturnsCreated      bool
	ReturnsRelationship bool
}

func buildSpec(variable string, params Parameters, labels ...string) (Parameters, string) {
	const labelSeparator = ":"
	propertiesString, scopedParameters := buildMaps(variable, params)
	labelsString := ""
	if len(labels) > 0 {
		labelsString = labelSeparator + strings.Join(labels, labelSeparator)
	}
	return scopedParameters,
		fmt.Sprintf("%s %s", labelsString, propertiesString)
}

type Parameters = map[string]any

func buildMaps(scope string, parameters Parameters) (string, Parameters) {
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

func CombinedParameters(parms ...Parameters) Parameters {
	out := make(Parameters)
	for _, o := range parms {
		for k, v := range o {
			out[k] = v
		}
	}
	return out
}

func NewTemplateSpec(parameters Parameters, labels ...string) TemplateSpec {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, rand.Uint32())
	variable := base64.RawURLEncoding.EncodeToString(bs)
	parameters, stringValue := buildSpec(variable, parameters, labels...)
	return TemplateSpec{
		Parameters:  parameters,
		stringValue: stringValue,
	}
}
