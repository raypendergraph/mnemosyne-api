package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"html/template"
	"mnemosyne-api/domain"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"
	"strings"
	"time"
)

const (
	notationVariable     = "note"
	journalEntryVariable = "je"
	journalEntryLabel    = "JournalEntry"
	entryOfLabel         = "ENTRY_OF"
)

type JournalEntryRepository struct {
	Neo4JAdapter
}

type createJournalEntryArgs struct {
	Journal      node
	EntryOfLabel string
	JournalEntry node
}

var createJournalEntryStmt = template.
	Must(template.New("createJournalEntry").
		Funcs(template.FuncMap{"JoinStrings": strings.Join}).
		Parse(`
	MATCH(j{{ .Journal }})<-[eo:{{ .EntryOfLabel }}]-(:{{ JoinStrings .JournalEntry.Labels ":" }})
	WITH eo.order AS o, j 
	ORDER BY eo.order DESC
	LIMIT 1
	WITH head(collect(o)) AS lastOrder, b
	CREATE(j)<-[:{{ JoinStrings JournalEntry.Labels ":" }} {order:lastOrder + 1}]-(je{{ .JournalEntry }}"})
	RETURN p
`))

func (r JournalEntryRepository) prepareCreateJournalEntryStmt(journalID sys.UUID) (string, parameters, sys.Error) {
	now := time.Now().UTC()
	journalNode := spec{
		Variable: journalVariable,
		Parameters: parameters{
			journal.FieldUUID.String(): journalID,
		},
		Labels: []string{journalLabel},
	}.Build()
	journalEntryNode := spec{
		Parameters: parameters{
			journal_entry.FieldUUID.String():      sys.NewUUID().String(),
			journal_entry.FieldCreatedAt.String(): now,
			journal_entry.FieldUpdatedAt.String(): now,
			journal_entry.FieldDeletedAt.String(): time.Time{},
			journal_entry.FieldTitle.String():     "",
			//TODO
			//journal_entry.FieldCaption.String():   caption,
		},
		Labels: []string{journalEntryLabel},
	}.Build()

	values := createJournalEntryArgs{
		Journal:      journalNode,
		EntryOfLabel: entryOfLabel,
		JournalEntry: journalEntryNode,
	}
	stmt, err := statement(createJournalEntryStmt, values)
	if err != nil {
		return "", nil, nil
	}
	return stmt, journalNode.Parameters.CombinedWith(journalEntryNode.Parameters), nil
}

func (r JournalEntryRepository) prepareCreateJournalEntryMapFn(ctx sys.ServiceContext) func(r neo4j.ResultWithContext) (journal_entry.TypeWithAssociations, error) {
	return func(r neo4j.ResultWithContext) (journal_entry.TypeWithAssociations, error) {
		records, err := r.Collect(ctx.GetPortContext())
		for _, record := range records{
				record.
		}

		//record, err := r.Single(ctx.GetPortContext())
		//if err != nil {
		//	return nil, err
		//}
		//je, err := recordToJournalEntry(ctx, record)
		//if err != nil {
		//	return nil, err
		//}
		//
		//return journal_entry.ImplWithAssociations{
		//	Type:      je,
		//	Taggable:  nil,
		//	Notations: nil,
		//}, nil
	}
}
func (r JournalEntryRepository) CreateJournalEntry(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[journal_entry.TypeWithAssociations] {
	stmt, args, err := r.prepareCreateJournalEntryStmt(journalID)
	if err != nil {
		return sys.Result[journal_entry.TypeWithAssociations]{Error: err.NextFrame()}
	}
	mapFn := r.prepareCreateJournalEntryMapFn(ctx)
	return WriteWithTransaction[journal_entry.TypeWithAssociations](
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
