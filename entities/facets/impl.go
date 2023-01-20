package facets

import (
	sys "mnemosyne-api/system"
	"time"
)

type GloballyIdentifiableImpl struct {
	UUID sys.UUID `json:"uuid"`
}

func (r GloballyIdentifiableImpl) GetUUID() sys.UUID {
	return r.UUID
}

type ListDisplayableImpl struct {
	GloballyIdentifiableImpl
	Title      string `json:"title"`
	Caption    string `json:"caption"`
	EntityKind string `json:"entity_kind"`
}

func (r ListDisplayableImpl) GetTitle() string {
	return r.Title
}

func (r ListDisplayableImpl) GetCaption() string {
	return r.Caption
}

func (r ListDisplayableImpl) GetEntityKind() string {
	return r.EntityKind
}

type TimeTrackableImpl struct {
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r TimeTrackableImpl) GetCreatedAt() time.Time {
	return r.CreatedAt
}

func (r TimeTrackableImpl) GetDeletedAt() time.Time {
	return r.DeletedAt
}

func (r TimeTrackableImpl) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}
