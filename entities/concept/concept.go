package concept

import (
	"mnemosyne-api/entities/facets"
)

type Type interface {
	facets.GloballyIdentifiable
	facets.ListDisplayable
	facets.TimeTrackable
}

type TypeWithAssociations interface {
	facets.Taggable
	Type

	GetNotations()
}
