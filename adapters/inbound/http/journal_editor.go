package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mnemosyne-api/adapters/outbound/neo4j"
	"mnemosyne-api/domain"
	"mnemosyne-api/domain/services"
	sys "mnemosyne-api/system"
	"net/http"
	"sync"
)

var journalEditorOnce = sync.Once{}
var journalEditor domain.JournalEditor

func getJournalEditor(base neo4j.Neo4JAdapter) domain.JournalEditor {
	journalEditorOnce.Do(func() {
		repository := neo4j.NewJournalRepository(base)
		journalEditor = services.NewJournalEditor(repository)
	})
	return journalEditor
}

func (r httpAdapter) createJournal(c *gin.Context) {
	serviceContext := r.newServiceContext(c)
	var command createJournalBinding
	if ok := bind(c, &command); !ok {
		return
	}
	serviceContext = serviceContext.WithMetaData(sys.MetaData{"command": command})
	editor := getJournalEditor(r.neo4JAdapter)
	if result := editor.CreateJournal(serviceContext, command); result.Error == nil {
		c.JSON(http.StatusOK, result.Value)
	} else {
		serviceContext.GetLogger().LogError(result.Error)
		cause := result.Error.GetCause().GetDefinition()
		c.JSON(cause.HTTPCode, cause.Message)
	}
}

type createJournalBinding struct {
	headerData
	Caption   string   `json:"caption"`
	Title     string   `json:"title"`
	JournalID sys.UUID `json:"journal_id"`
}

func (r createJournalBinding) GetJournalID() uuid.UUID {
	return uuid.UUID(r.JournalID)
}

func (r createJournalBinding) GetTitle() string {
	return r.Title
}

func (r createJournalBinding) GetCaption() string {
	return r.Caption
}

type deleteJournalBinding struct {
	headerData
	JournalID uuid.UUID
}

func (r deleteJournalBinding) GetJournalID() uuid.UUID {
	return r.JournalID
}

type updateJournalAttributesBinding struct {
	headerData
	Caption string
	Title   string
}

func (r updateJournalAttributesBinding) GetTitle() string {
	return r.Title
}

func (r updateJournalAttributesBinding) GetCaption() string {
	return r.Caption
}

type addJournalEntryBinding struct {
	headerData
	JournalID uuid.UUID
	Title     string
}

func (r addJournalEntryBinding) GetJournalID() uuid.UUID {
	return r.JournalID
}

func (r addJournalEntryBinding) GetTitle() string {
	return r.Title
}
