package domain

import (
	e "mnemosyne-api/entities"
	"mnemosyne-api/entities/concept"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	"mnemosyne-api/entities/notation"
	"mnemosyne-api/entities/organization"
	"mnemosyne-api/entities/person"
	sys "mnemosyne-api/system"
)

type ConceptRepositoryAdapter interface {
	FetchConceptWithAssociations(ctx sys.ServiceContext, conceptID sys.UUID) sys.Result[concept.TypeWithAssociations]
	TagConcept(ctx sys.ServiceContext, conceptID sys.UUID, tags []string) sys.Result[sys.Void]
	UntagConcept(ctx sys.ServiceContext, conceptID sys.UUID, tags []string) sys.Result[sys.Void]
	DeleteConcept(ctx sys.ServiceContext, conceptID sys.UUID) sys.Result[sys.Void]
	AddNotation(ctx sys.ServiceContext, conceptID sys.UUID)
}

type JournalRepositoryAdapter interface {
	CreateJournal(ctx sys.ServiceContext, title, caption string) sys.Result[journal.Type]
	FetchJournalWithAssociations(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[journal.TypeWithAssociations]
	ListJournalEntries(ctx sys.ServiceContext, journalID sys.UUID) sys.Result[e.PagedList[facets.ListDisplayable]]
	TagJournal(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void]
	UntagJournal(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void]
	DeleteJournal(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[sys.Void]
	ModifyJournalFields(ctx sys.ServiceContext, fields map[journal.Field]any) sys.Result[sys.Void]
}

type JournalEntryRepositoryAdapter interface {
	CreateJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[journal_entry.Type]
	FetchJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[journal_entry.TypeWithAssociations]
	TagJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void]
	UntagJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID, tags []string) sys.Result[sys.Void]
	DeleteJournalEntry(ctx sys.ServiceContext, journalEntryID sys.UUID) sys.Result[sys.Void]
}

type NotationRepositoryAdapter interface {
	CreateNotation(ctx sys.ServiceContext, annotatedID sys.UUID) sys.Result[notation.Type]
	UpdateNotation(ctx sys.ServiceContext, notation notation.Type) sys.Result[notation.Type]
	DeleteNotation(ctx sys.ServiceContext, notationID sys.UUID) sys.Result[sys.Void]
}

type OrganizationRepositoryAdapter interface {
	FetchOrganization(ctx sys.ServiceContext, organizationID sys.UUID) sys.Result[organization.TypeWithAssociation]
	TagOrganization(ctx sys.ServiceContext, organizationID sys.UUID, tags []string) sys.Result[sys.Void]
	UntagOrganization(ctx sys.ServiceContext, organizationID sys.UUID, tags []string) sys.Result[sys.Void]
	DeleteOrganization(ctx sys.ServiceContext, organizationID sys.UUID) sys.Result[sys.Void]
	AddNotation(ctx sys.ServiceContext, organizationID sys.UUID)
}

type PersonRepositoryAdapter interface {
	FetchPerson(ctx sys.ServiceContext, personID sys.UUID) sys.Result[person.TypeWithAssociations]
	TagPerson(ctx sys.ServiceContext, personID sys.UUID, tags []string) sys.Result[sys.Void]
	UntagPerson(ctx sys.ServiceContext, personID sys.UUID, tags []string) sys.Result[sys.Void]
	DeletePerson(ctx sys.ServiceContext, personID sys.UUID) sys.Result[sys.Void]
	AddNotation(ctx sys.ServiceContext, personID sys.UUID)
}
