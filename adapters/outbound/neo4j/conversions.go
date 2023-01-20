package neo4j

import (
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type nodeMapping[T any] func(ctx sys.ServiceContext, n neo4j.Node) (T, sys.Error)
type resultMapping[T any] func(result neo4j.ResultWithContext) (T, error)

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

func nodeToJournalEntry(ctx sys.ServiceContext, node neo4j.Node) (journal_entry.Impl, sys.Error) {
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

func nodeToJournal(ctx sys.ServiceContext, node neo4j.Node) (journal.Type, sys.Error) {
	var out journal.Impl
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
