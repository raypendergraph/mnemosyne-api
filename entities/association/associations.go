package association

import (
	"mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/notation"
)

type Type interface {
	facets.ListDisplayable
	GetAssociationAssociationKind() entities.AssociationKind
}
type Impl struct {
	facets.ListDisplayableImpl
	AssociationKind entities.AssociationKind
}

func (r Impl) GetEntityKind() string {
	return r.EntityKind
}

func (r Impl) GetAssociationKind() string {
	return r.AssociationKind
}

func (r Impl) GetAssociationAssociationKind() entities.AssociationKind {
	return r.AssociationKind
}

type Associating interface {
	GetAssociations() []Type
	GetNotations() []notation.Type
	GetTags() []string
}

type AssociatingImpl struct {
	Associations []Type
	Notations    []notation.Type
	Tags         []string
}

func (r AssociatingImpl) GetAssociations() []Type {
	return r.Associations
}

func (r AssociatingImpl) GetNotations() []notation.Type {
	return r.Notations
}

func (r AssociatingImpl) GetTags() []string {
	return r.Tags
}
