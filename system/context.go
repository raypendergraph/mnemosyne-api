package system

import (
	"context"
)

type ServiceContext interface {
	NewError(code ErrorCode, causes ...error) Error
	GetPortContext() context.Context
	GetLogger() Logger
	WithMetaData(m MetaData) ServiceContext
}
