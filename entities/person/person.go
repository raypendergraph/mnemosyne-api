package person

import (
	"mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/notation"
	"time"
)

type Type interface {
	facets.GloballyIdentifiable
	facets.ListDisplayable
	GetBirthTime() time.Time
	GetDeathTime() *time.Time
}

type TypeWithAssociations interface {
	Type
	GetNotations() entities.PagedList[notation.Type]
}
