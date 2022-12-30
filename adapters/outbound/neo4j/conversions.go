package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"
)

func nodeToGloballyIdentifiable(ctx sys.ServiceContext, node neo4j.Node) (facets.GloballyIdentifiable, sys.Error) {
	var out facets.GloballyIdentifiableImpl
	if err := bindUUIDProperty(journal.FieldUUID.String(), node, &out.UUID); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	return out, nil
}

func nodeToTimeTrackable(ctx sys.ServiceContext, node neo4j.Node) (facets.TimeTrackable, sys.Error) {
	var out facets.TimeTrackableImpl
	if err := bindDateTimeProperty(journal.FieldDeletedAt.String(), node, &out.DeletedAt); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}

	if err := bindDateTimeProperty(journal.FieldCreatedAt.String(), node, &out.CreatedAt); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	if err := bindDateTimeProperty(journal.FieldUpdatedAt.String(), node, &out.UpdatedAt); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	return out, nil
}

func nodeToListDisplayable(ctx sys.ServiceContext, node neo4j.Node) (facets.ListDisplayable, sys.Error) {
	var out facets.ListDisplayableImpl
	if err := bindNeo4JProperty(journal.FieldTitle.String(), node, &out.Title); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}

	if err := bindNeo4JProperty(journal.FieldCaption.String(), node, &out.Caption); err != nil {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure, err)
	}
	return out, nil
}

func nodeToJournalEntry(ctx sys.ServiceContext, node neo4j.Node) (journal_entry.Type, sys.Error) {
	var out journal_entry.Impl
	if ga, err := nodeToGloballyIdentifiable(ctx, node); err == nil {
		out.GloballyIdentifiable = ga
	} else {
		return nil, err
	}

	if tt, err := nodeToTimeTrackable(ctx, node); err == nil {
		out.TimeTrackable = tt
	} else {
		return nil, err
	}

	if ld, err := nodeToListDisplayable(ctx, node); err == nil {
		out.ListDisplayable = ld
	} else {
		return nil, err
	}

	return out, nil
}

func recordToJournal(ctx sys.ServiceContext, record *neo4j.Record) (journal.Type, sys.Error) {
	ctx = ctx.WithMetaData(sys.MetaData{
		"key":    journalVariable,
		"record": fmt.Sprintf("%#v", record),
	})
	var out journal.Impl
	rawItemNode, found := record.Get(journalVariable)
	if !found {
		return nil, ctx.NewError(sys.ErrNeo4JUnmarshallingFailure)
	}

	rin := rawItemNode.(neo4j.Node)
	if ga, err := nodeToGloballyIdentifiable(ctx, rin); err == nil {
		out.GloballyIdentifiable = ga
	} else {
		return nil, err
	}

	if tt, err := nodeToTimeTrackable(ctx, rin); err == nil {
		out.TimeTrackable = tt
	} else {
		return nil, err
	}

	if ld, err := nodeToListDisplayable(ctx, rin); err == nil {
		out.ListDisplayable = ld
	} else {
		return nil, err
	}

	return out, nil
}
