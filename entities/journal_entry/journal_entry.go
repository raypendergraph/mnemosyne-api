package journal_entry

import (
	e "mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/notation"
)

const EntityName = "JournalEntry"

type Type interface {
	facets.GloballyIdentifiable
	facets.ListDisplayable
	facets.TimeTrackable
}

type TypeWithAssociations interface {
	GetNotations() e.PagedList[notation.Type]
	facets.Taggable
	Type
}

type Field int

const (
	FieldUUID Field = 1 << iota
	FieldTitle
	FieldCaption
	FieldCreatedAt
	FieldUpdatedAt
	FieldDeletedAt
)

func (r Field) String() string {
	switch r {
	case FieldUUID:
		return "uuid"
	case FieldTitle:
		return "title"
	case FieldCaption:
		return "caption"
	case FieldCreatedAt:
		return "created_at"
	case FieldUpdatedAt:
		return "updated_at"
	case FieldDeletedAt:
		return "deleted_at"
	default:
		return ""
	}
}

type Impl struct {
	facets.GloballyIdentifiable
	facets.TimeTrackable
	facets.ListDisplayable
}

type ImplWithAssociations struct {
	Type
	facets.Taggable
	Notations e.PagedList[notation.Type]
}
