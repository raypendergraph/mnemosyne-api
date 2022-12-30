package journal

import (
	e "mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal_entry"
)

type Type interface {
	facets.GloballyIdentifiable
	facets.ListDisplayable
	facets.TimeTrackable
}

type TypeWithAssociations interface {
	facets.Taggable
	Type
	GetEntries() e.PagedList[journal_entry.Type]
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
	Impl
	Tags    []string
	Entries e.PagedList[journal_entry.Type]
}

func (r ImplWithAssociations) GetTags() []string {
	return r.Tags
}

func (r ImplWithAssociations) GetEntries() e.PagedList[journal_entry.Type] {
	return r.Entries
}
