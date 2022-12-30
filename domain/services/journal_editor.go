package services

import (
	"mnemosyne-api/domain"
	"mnemosyne-api/entities/fieldmask"
	"mnemosyne-api/entities/journal"
	"mnemosyne-api/entities/journal_entry"
	sys "mnemosyne-api/system"
	"time"
)

type journalEditor struct {
	journalRepository domain.JournalRepositoryAdapter
}

func (r journalEditor) CreateJournal(ctx sys.ServiceContext, command domain.CreateJournalCommand) (out sys.Result[journal.Type]) {
	t := time.Now().UTC()
	entity := journal.Impl{
		UUID:      sys.NewUUID(),
		Title:     command.GetTitle(),
		Caption:   command.GetCaption(),
		CreatedAt: t,
		UpdatedAt: t,
	}

	if messages := entity.Validate(); len(messages) > 0 {
		out.Error = ctx.WithMetaData(sys.MetaData{"messages": messages}).
			NewError(sys.ErrDomainInvariantViolation)
		return
	}
	return r.journalRepository.CreateJournal(ctx, entity)
}

func (r journalEditor) UpdateJournalAttributes(ctx sys.ServiceContext, command domain.UpdateJournalAttributesCommand) sys.Result[sys.Void] {
	var fields []journal.Field
	fieldmask.EnumerateFields(command.GetFieldMask(), &fields)
	updateMap := make(map[journal.Field]any)
	for _, field := range fields {
		switch field {
		case journal.FieldTitle:
			updateMap[journal.FieldTitle] = command.GetTitle()
		case journal.FieldCaption:
			updateMap[journal.FieldCaption] = command.GetCaption()
		default:
			//return error
		}
	}
	return sys.Result[sys.Void]{
		Error: r.journalRepository.ModifyJournalFields(ctx, updateMap).Error,
	}
}

func (r journalEditor) RemoveJournal(ctx sys.ServiceContext, command domain.DeleteJournalCommand) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func (r journalEditor) CreateJournalEntry(ctx sys.ServiceContext, command domain.AddJournalEntryCommand) sys.Result[journal_entry.Type] {
	//TODO implement me
	panic("implement me")
}

func (r journalEditor) RemoveJournalEntry(ctx sys.ServiceContext) sys.Result[sys.Void] {
	//TODO implement me
	panic("implement me")
}

func NewJournalEditor(repository domain.JournalRepositoryAdapter) domain.JournalEditor {
	return journalEditor{journalRepository: repository}
}
