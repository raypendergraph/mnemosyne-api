package notation

import (
	"mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
)

type Format = string

const (
	FormatCreole    Format = "creole"
	FormatPlainText Format = "plaintext"
)

type Kind string

const (
	KindInline       Kind = "inline"
	KindUrl          Kind = "url"
	KindBibliography Kind = "bibliography"
)

type Type interface {
	facets.GloballyIdentifiable
	facets.TimeTrackable
	GetContent() string
	GetFormat() Format
}

type TypeWithAssociations interface {
	Type
}

type Annotated interface {
	GetNotations() entities.PagedList[Type]
}

type AnnotatedImpl struct {
	Notations entities.PagedList[Type]
}

func (r AnnotatedImpl) GetNotations() entities.PagedList[Type] {
	return r.Notations
}