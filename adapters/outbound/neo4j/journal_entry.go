package neo4j

import (
	"fmt"
	"mnemosyne-api/domain"
	"mnemosyne-api/entities"
	"mnemosyne-api/entities/association"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"
	t "mnemosyne-api/templates"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	notationVariable     = "note"
	journalEntryVariable = "je"
	journalEntryLabel    = "JournalEntry"
	hasEntryLabel        = "HAS_ENTRY"
)

type JournalEntryRepository struct {
	Neo4JAdapter
}

func (r JournalEntryRepository) prepareCreateJournalEntryStmt(journalID sys.UUID) (string, statementParams, sys.Error) {
	now := time.Now().UTC()
	journalSpec := t.NewTemplateSpec(
		t.Parameters{
			journal.FieldUUID.String(): journalID,
		},
		journalLabel,
	)

	journalEntrySpec := t.NewTemplateSpec(
		t.Parameters{
			journal_entry.FieldUUID.String():      sys.NewUUID().String(),
			journal_entry.FieldCreatedAt.String(): now,
			journal_entry.FieldUpdatedAt.String(): now,
			journal_entry.FieldDeletedAt.String(): time.Time{},
			journal_entry.FieldTitle.String():     "",
			//TODO
			//journal_entry.FieldCaption.String():   caption,
		},
		journalEntryLabel,
	)

	args := t.CreateJournalEntryArgs{
		JournalMatch: t.MatchChain{
			Origin: &journalSpec,
		},
		JournalEntry:      journalEntrySpec,
		JournalEntryVar:   journalEntryVariable,
		RelationshipLabel: hasEntryLabel,
	}

	stmt, err := statement(t.Cypher, args)
	if err != nil {
		return "", nil, nil
	}
	return stmt, journalSpec.Parameters.CombinedWith(journalEntrySpec.Parameters), nil
}
func getColumnNode(ctx sys.ServiceContext, record *neo4j.Record, columnName string) (out neo4j.Node, found bool, err sys.Error) {
	c, ok := record.Get(columnName)
	if !ok {
		found = false
		return
	}
	out, ok = c.(neo4j.Node)
	if !ok {
		err = ctx.NewError(sys.ErrNeo4JUnmarshallingFailure,
			fmt.Errorf("could not conver column %s which is type %T to Node", columnName, c))
		return
	}
	return
}
func newAssociationFromNodes(ctx sys.ServiceContext, a, r neo4j.Node) (association.Type, sys.Error) {
	ctx = ctx.WithMetaData(sys.MetaData{"association": a, "relationship": r})
	var out association.Impl
	if len(r.Labels) == 0 {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, fmt.Errorf("relationship has no label"))
	}
	out.AssociationKind = r.Labels[0]
	if err := bindUUIDProperty(entities.FieldUUID.String(), a, &out.UUID); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}

	if err := bindNeo4JProperty[string](entities.FieldTitle.String(), a, &out.Title); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	if err := bindNeo4JProperty[string](entities.FieldCaption.String(), a, &out.Caption); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}

	if err := bindNeo4JProperty[string](entities.FieldEntityKind.String(), a, &out.EntityKind); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	return out, nil
}

func newTagFromRecord(ctx sys.ServiceContext, associationVar string, record *neo4j.Record) (out string, found bool, err sys.Error) {
	n, found, err := getColumnNode(ctx, record, associationVar)
	if err != nil {
		err = err.NextFrame()
		return
	}
	if !found {
		return
	}

	value, ok := n.Props[entities.FieldTagName.String()]
	if !ok {
		found = false
		ctx.GetLogger().LogWarning(fmt.Sprintf("%s was not a property on tag node %d", entities.FieldTagName.String(), n.GetId()))
	}

	out, ok = value.(string)
	if !ok {
		ctx.GetLogger().LogWarning(fmt.Sprintf("%s property on %d is not a string: %v", entities.FieldTagName.String(), n.GetId(), value))
	}

	return
}

func newAssociationFromRecord(ctx sys.ServiceContext, relationshipVar, associationVar string, record *neo4j.Record) (out association.Type, found bool, err sys.Error) {
	relNode, found, err := getColumnNode(ctx, record, relationshipVar)
	if err != nil {
		err = err.NextFrame()
		return
	}
	if !found {
		return
	}
	assNode, found, err := getColumnNode(ctx, record, associationVar)
	if err != nil {
		err = err.NextFrame()
		return
	}
	if !found {
		return
	}

	out, err = newAssociationFromNodes(ctx, assNode, relNode)
	return
}

func processJournalEntryAssociations(records []*neo4j.Record) (journal_entry.ImplWithAssociations, sys.Error) {

}

func prepareJournalEntryMapFn(ctx sys.ServiceContext) resultMapping[journal_entry.TypeWithAssociations] {

	return func(r neo4j.ResultWithContext) (journal_entry.TypeWithAssociations, error) {
		var records []*neo4j.Record
		if r, err := r.Collect(ctx.GetPortContext()); err == nil {
			records = r
		} else {
			return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
		}
		if len(records) == 0 {
			return nil, nil
		}
		journalColumn, found := records[0].Get(journalVariable)
		if !found {
			return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure)
		}
		var je journal_entry.ImplWithAssociations
		var serr sys.Error
		if je.Type, serr = nodeToJournalEntry(ctx, journalColumn.(neo4j.Node)); serr != nil {
			return nil, serr.NextFrame()
		}

	}
}

func (r JournalEntryRepository) CreateJournalEntry(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[journal_entry.TypeWithAssociations] {
	stmt, args, err := r.prepareCreateJournalEntryStmt(journalID)
	if err != nil {
		return sys.Result[journal_entry.TypeWithAssociations]{Error: err.NextFrame()}
	}
	mapFn := buildMapManagedResultsToSingleFn(ctx, journalEntryVariable, nodeToJournalEntry)
	return writeWithTransaction[journal_entry.TypeWithAssociations](
		TransactionArguments[journal_entry.TypeWithAssociations]{
			Context:    ctx,
			Statement:  stmt,
			Args:       args,
			Driver:     r.driver,
			MapResults: mapFn,
		})
}

func (r JournalEntryRepository) FetchJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[journal_entry.TypeWithAssociations] {
	//TODO implement me
	panic("implement me")
}

func (r JournalEntryRepository) TagJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r JournalEntryRepository) UntagJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r JournalEntryRepository) DeleteJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func NewJournalEntryRepository(base Neo4JAdapter) domain.JournalEntryRepositoryAdapter {
	return JournalEntryRepository{base}
}
