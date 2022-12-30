package entities

type PageRequester interface {
	GetSkip() uint
	GetLimit() uint
}

type PagedList[T any] interface {
	GetSkipped() uint
	GetSize() uint
	GetValues() []T
}

type EntityValidator interface {
	Validate() []string
}
