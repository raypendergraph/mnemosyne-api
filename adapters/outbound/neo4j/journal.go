package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"mnemosyne-api/domain"
	e "mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal"
	sys "mnemosyne-api/system"
	"time"
)

const (
	tagsVariable    = "tags"
	journalVariable = "j"
	journalLabel    = "Journal"
)

func NewJournalRepository(base Neo4JAdapter) domain.JournalRepositoryAdapter {
	return JournalRepository{base}
}

type JournalRepository struct {
	Neo4JAdapter
}

func (r JournalRepository) CreateJournal(ctx sys.ServiceContext, title, caption string) sys.Result[journal.Type] {
	stmt, args := r.prepareCreateJournalRequest(title, caption)
	mapFn := r.prepareCreateJournalMapFn(ctx)
	return writeWithTransaction[journal.Type](
		TransactionArguments[journal.Type]{
			Context:    ctx,
			Statement:  stmt,
			Args:       args,
			Driver:     r.driver,
			MapResults: mapFn,
		})
}
func (r JournalRepository) prepareCreateJournalMapFn(ctx sys.ServiceContext) resultMapping[journal.Type] {
	return func(r neo4j.ResultWithContext) (journal.Type, error) {
		record, err := r.Single(ctx.GetPortContext())
		if err != nil {
			return nil, err
		}
		column, found := record.Get(journalVariable)
		if !found {
			return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure)
		}

		n := column.(neo4j.Node)
		return nodeToJournal(ctx, n)
	}
}
func (r JournalRepository) prepareCreateJournalRequest(title, caption string) (string, parameters) {
	now := time.Now().UTC()
	journalNode := spec{
		Variable: journalVariable,
		Parameters: parameters{
			journal.FieldUUID.String():      sys.NewUUID().String(),
			journal.FieldCreatedAt.String(): now,
			journal.FieldUpdatedAt.String(): now,
			journal.FieldDeletedAt.String(): time.Time{},
			journal.FieldTitle.String():     title,
			journal.FieldCaption.String():   caption,
		},
		Labels: []string{journalLabel},
	}.Build()

	return fmt.Sprintf("CREATE (%s) RETURN(%s)", journalNode, journalNode.Variable), journalNode.Parameters
}

func (r JournalRepository) FetchJournalWithAssociations(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[journal.TypeWithAssociations] {
	//TODO implement me
	panic("implement me")
}

func (r JournalRepository) ListJournalEntries(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[e.PagedList[facets.ListDisplayable]] {
	//TODO implement me
	panic("implement me")
}

func (r JournalRepository) TagJournal(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r JournalRepository) UntagJournal(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r JournalRepository) DeleteJournal(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r JournalRepository) ModifyJournalFields(ctx sys.ServiceContext, fields map[journal.Field]any) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}
