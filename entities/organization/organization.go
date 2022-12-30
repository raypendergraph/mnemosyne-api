package organization

import (
	"mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"time"
)

type Type interface {
	facets.GloballyIdentifiable
	facets.ListDisplayable
	GetDomicile() string
	GetTerminationDate() *time.Time
}

type TypeWithAssociation interface {
	facets.Taggable
	Type

	GetNotations() entities.PagedList[facets.ListDisplayable]
}
