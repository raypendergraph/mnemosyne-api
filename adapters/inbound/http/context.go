package http

import (
	"context"
	"github.com/gin-gonic/gin"
	sys "mnemosyne-api/system"
)

func (r httpAdapter) newServiceContext(c *gin.Context) sys.ServiceContext {
	// TODO logger stuff from gin
	return httpContext{
		portContext: c,
		errors:      r.errors,
		logger:      r.logger,
		metaData:    sys.MetaData{},
	}
}

type httpContext struct {
	portContext *gin.Context
	errors      sys.ErrorCatalog
	logger      sys.Logger
	metaData    sys.MetaData
}

func (r httpContext) NewError(code sys.ErrorCode, causes ...error) sys.Error {
	return r.errors.NewError(code, r.metaData, causes...)
}

func (r httpContext) GetPortContext() context.Context {
	return r.portContext
}

func (r httpContext) GetLogger() sys.Logger {
	return r.logger.WithAdditionalMetaData(r.metaData)
}

func (r httpContext) WithMetaData(additionalMeta sys.MetaData) sys.ServiceContext {
	return httpContext{
		portContext: r.portContext,
		errors:      r.errors,
		metaData:    r.metaData.CombinedWith(additionalMeta),
		logger:      r.logger,
	}
}
