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
	GetEntityKind() string
}

type Taggable interface {
	GetTags() []string
}

type TaggableImpl struct {
	Tags []string
}

func (r TaggableImpl) GetTags() []string {
	return r.Tags
}

type TimeTrackable interface {
	GetCreatedAt() time.Time
	GetDeletedAt() time.Time
	GetUpdatedAt() time.Time
}
