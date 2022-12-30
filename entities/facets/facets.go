package facets

import (
	sys "mnemosyne-api/system"
	"time"
)

type GloballyIdentifiable interface {
	GetUUID() sys.UUID
}

type ListDisplayable interface {
	GloballyIdentifiable
	GetTitle() string
	GetCaption() string
}

type Taggable interface {
	GetTags() []string
}

type TimeTrackable interface {
	GetCreatedAt() time.Time
	GetDeletedAt() time.Time
	GetUpdatedAt() time.Time
}
