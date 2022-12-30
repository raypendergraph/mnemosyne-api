package system

type Logger interface {
	LogAlways(msg string)
	LogDebug(msg string)
	LogError(e Error)
	LogWarning(msg string, errors ...error)
	WithAdditionalMetaData(values MetaData) Logger
}
