package system

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type ErrorCode uint

// String - implementation of fmt.Stringer
func (r ErrorCode) String() string {
	return fmt.Sprintf("%d", r)
}

type ErrorDefinition struct {
	Code     ErrorCode
	Message  string
	HTTPCode int
}

type Error interface {
	error
	GetCause() Cause
	NextFrame() Error
}

type Cause interface {
	error
	GetMetaData() MetaData
	GetDefinition() ErrorDefinition
}

type HttpStatusCode int

type ErrorCatalog interface {
	NewError(code ErrorCode, meta MetaData, causes ...error) Error
}
type CatalogConfig = map[ErrorCode]ErrorDefinition
type causeImpl struct {
	causes     []error
	definition ErrorDefinition
	metaData   MetaData
}

// Error implementation of errors.Error
func (r causeImpl) Error() string {
	causeStrings := make([]string, len(r.causes))
	for i, cause := range r.causes {
		causeStrings[i] = cause.Error()
	}
	return fmt.Sprintf("{code:%s, message: %s, causes: [%s]",
		r.definition.Code, r.definition.Message, strings.Join(causeStrings, ", "))
}

func (r causeImpl) GetMetaData() MetaData {
	return r.metaData
}

func (r causeImpl) GetDefinition() ErrorDefinition {
	return r.definition
}

type frame struct {
	next        *frame
	cause       Cause
	file        string
	function    string
	line        int
	frameOffset int
}

type frameVisitor func(f frame, i int)

func (r frame) iterateList(it frameVisitor) {
	for f, i := &r, 0; f != nil; f, i = f.next, i+1 {
		it(*f, i)
	}
}

// Error implementation of errors.Error returns the error message as a string in the form of:
//
//	"[main.FunctionA]->[module/package1.FunctionB]->[module/package2.FunctionC]->[original error message]"
func (r frame) Error() string {
	const separator = "->"
	var frames []string
	var last frame
	r.iterateList(func(f frame, i int) {
		last = f
		pathParts := strings.Split(f.file, string(os.PathSeparator))
		functionParts := strings.Split(f.function, string(os.PathSeparator))
		frames = append(frames, fmt.Sprintf("[%s/%s:%d]", pathParts[len(pathParts)-1], functionParts[len(functionParts)-1], f.line))
	})
	// join all frames with a "->"
	return fmt.Sprintf("%s - [%s]", strings.Join(frames, separator), last.cause.Error())
}

// GetCause gets the initiating Cause of this Error by traversing its internal list.
func (r frame) GetCause() Cause {
	var last frame
	r.iterateList(func(f frame, i int) {
		last = f
	})
	return last.cause
}

// NextFrame adds a frame to this error as it is passed up the call chain. This is used to build the call stack when
// printed. If this is not called before a returned error is returned by the current function, it will be missing in the
// stack trace.
func (r frame) NextFrame() Error {
	pc, file, line, ok := runtime.Caller(1 + r.frameOffset)
	if !ok {
		return nil
	}
	functionStr := runtime.FuncForPC(pc).Name()
	return frame{
		next:        &r,
		cause:       nil,
		file:        file,
		function:    functionStr,
		line:        line,
		frameOffset: r.frameOffset,
	}
}

func newFrame(skip int) frame {
	pc, file, line, _ := runtime.Caller(skip)
	functionStr := runtime.FuncForPC(pc).Name()

	return frame{
		file:     file,
		function: functionStr,
		line:     line,
	}
}

type MetaData map[string]any

func (r MetaData) CombinedWith(additional MetaData) MetaData {
	// Shallow copy will work fine here, these are read only values anyway.
	dup := make(MetaData, len(r)+len(additional))
	for k, v := range r {
		dup[k] = v
	}
	for k, v := range additional {
		dup[k] = v
	}
	return dup
}

func NewCatalog(frameOffset int) ErrorCatalog {
	return catalog{
		frameSkip: frameOffset,
		config:    catalogConfig,
	}
}

type catalog struct {
	frameSkip int
	config    CatalogConfig
}

func (c catalog) NewError(code ErrorCode, meta MetaData, causes ...error) Error {
	const codeKey = "code"
	definition := c.config[code]
	if meta == nil {
		meta = MetaData{
			codeKey: definition.Code,
		}
	} else {
		meta[codeKey] = definition.Code
	}
	cause := causeImpl{
		causes:     causes,
		definition: definition,
		metaData:   meta,
	}
	f := newFrame(1 + c.frameSkip)
	f.cause = &cause
	f.frameOffset = c.frameSkip
	return f
}

// MaybeError - implementation of ErrorCatalog
func (c catalog) MaybeError(code ErrorCode, meta MetaData, err error) Error {
	if err == nil {
		return nil
	}
	return c.NewError(code, meta, err)
}

const (
	ErrHTTPServerConfigurationFailure ErrorCode = 10
	ErrNeo4JConfigurationFailure      ErrorCode = 100
	ErrNeo4JTransactedWriteFailure    ErrorCode = 101
	ErrNeo4JUnmarshallingFailure      ErrorCode = 102
	ErrDomainInvariantViolation       ErrorCode = 1000
)

var catalogConfig = CatalogConfig{
	ErrHTTPServerConfigurationFailure: ErrorDefinition{
		Code:     ErrHTTPServerConfigurationFailure,
		Message:  "failed to stand up the HTTP server",
		HTTPCode: http.StatusInternalServerError,
	},
	ErrNeo4JConfigurationFailure: ErrorDefinition{
		Code:     ErrNeo4JConfigurationFailure,
		Message:  "failed to initially configure neo4j",
		HTTPCode: http.StatusFailedDependency,
	},
	ErrNeo4JTransactedWriteFailure: ErrorDefinition{
		Code:     ErrNeo4JTransactedWriteFailure,
		Message:  "could not perform a transacted write to neo4j",
		HTTPCode: http.StatusFailedDependency,
	},
	ErrNeo4JUnmarshallingFailure: ErrorDefinition{
		Code:     ErrNeo4JUnmarshallingFailure,
		Message:  "could not properly unmarshal neo4j results",
		HTTPCode: http.StatusFailedDependency,
	},
	ErrDomainInvariantViolation: ErrorDefinition{
		Code:     ErrDomainInvariantViolation,
		Message:  "an error occurred trying to apply invariants",
		HTTPCode: http.StatusConflict,
	},
}
