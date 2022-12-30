package domain

import (
	"github.com/google/uuid"
	e "mnemosyne-api/entities"
	"mnemosyne-api/entities/facets"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"
	"time"
)

type JournalEditor interface {
	CreateJournal(ctx sys.ServiceContext, command CreateJournalCommand) sys.Result[journal.Type]
	UpdateJournalAttributes(ctx sys.ServiceContext, command UpdateJournalAttributesCommand) sys.Result[sys.Void]
	RemoveJournal(ctx sys.ServiceContext, command DeleteJournalCommand) sys.Result[sys.Void]
	CreateJournalEntry(ctx sys.ServiceContext, command AddJournalEntryCommand) sys.Result[journal_entry.Type]
	RemoveJournalEntry(ctx sys.ServiceContext) sys.Result[sys.Void]
}
type CreateJournalCommand interface {
	Command
	GetJournalID() uuid.UUID
	GetTitle() string
	GetCaption() string
}

type UpdateJournalAttributesCommand interface {
	Command
	GetFieldMask() journal.Field
	GetTitle() string
	GetCaption() string
}

type DeleteJournalCommand interface {
	Command
	GetJournalID() uuid.UUID
}

type AddJournalEntryCommand interface {
	Command
	GetJournalID() uuid.UUID
	GetTitle() string
}

type JournalReader interface {
	ListJournals(pager e.PageRequester) (e.PagedList[facets.ListDisplayable], error)
	ReadJournal(journalID uuid.UUID, pager e.PageRequester) (journal.TypeWithAssociations, error)
}

type Command interface {
	// GetTransactionID gets the transaction ID for this command. If there is not one, it will be generated.
	GetTransactionID() uuid.UUID

	// GetTimestamp when this command was initiated
	GetTimestamp() time.Time

	// GetActorID the actor for which this command was initiated
	GetActorID() uuid.UUID
}
