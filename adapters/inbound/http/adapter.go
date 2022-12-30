package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mnemosyne-api/adapters/outbound/neo4j"
	"net/http"
	"time"

	sys "mnemosyne-api/system"
)

func NewAdapter(configuration sys.HTTPConfiguration, errs sys.ErrorCatalog, baseDBAdapter neo4j.Neo4JAdapter, logger sys.Logger) *gin.Engine {
	engine := gin.Default()
	adapter := httpAdapter{
		engine:       engine,
		errors:       errs,
		logger:       logger,
		neo4JAdapter: baseDBAdapter,
	}
	adapter.bindRoutes()
	return engine
}

type httpAdapter struct {
	engine       *gin.Engine
	errors       sys.ErrorCatalog
	logger       sys.Logger
	neo4JAdapter neo4j.Neo4JAdapter
}

func (r httpAdapter) bindRoutes() {
	router := r.engine
	v1 := router.Group("/v1")

	journals := v1.Group("/journals")
	journals.POST("", r.createJournal)
	journals.GET("/:id", r.getJournal)
	journals.PATCH("/:id", r.updateJournalAttributes)
	journals.DELETE("/:id", r.deleteJournal)

	journalEntries := journals.Group("/:id/entries")
	journalEntries.GET("", nil)
	journalEntries.POST("", r.addJournalEntry)
}

func (r httpAdapter) getJournal(c *gin.Context) {

}

func (r httpAdapter) updateJournalAttributes(c *gin.Context) {

}

func (r httpAdapter) deleteJournal(c *gin.Context) {

}

func (r httpAdapter) addJournalEntry(c *gin.Context) {

}

type headerData struct {
	TransactionID uuid.UUID `header:"transaction-id" binding:`
	Timestamp     time.Time `header:"-"`
	ActorID       uuid.UUID `header:"actor-id"`
}

func (r headerData) GetTransactionID() uuid.UUID {
	return r.TransactionID
}

func (r headerData) GetTimestamp() time.Time {
	return r.Timestamp
}

func (r headerData) GetActorID() uuid.UUID {
	return r.ActorID
}

func bind[T any](c *gin.Context, binding T) bool {

	if err := c.ShouldBindHeader(binding); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return false
	}
	if err := c.ShouldBindJSON(binding); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return false
	}
	return true
}
