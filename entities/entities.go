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

type AssociationKind = string
type EntityKind = string

const (
	EntityKindJournal      EntityKind = "Journal"
	EntityKindJournalEntry EntityKind = "Journal"
)

const (
	AssociationKindAnnotation   AssociationKind = "annotated_by"
	AssociationKindJournalEntry AssociationKind = "journaled_by"
)

var validationMap = map[string]any{
	AssociationKindAnnotation:   nil,
	AssociationKindJournalEntry: nil,
}

func Exists(k AssociationKind) (exists bool) {
	_, exists = validationMap[k]
	return
}

type EntityField int64

const (
	FieldUUID EntityField = 1 << iota

	FieldCaption
	FieldCreatedAt
	FieldDeletedAt
	FieldEntityKind
	FieldTagName
	FieldTitle
	FieldUpdatedAt
)

func (r EntityField) String() string {
	switch r {
	case FieldUUID:
		return "uuid"
	case FieldCaption:
		return "caption"
	case FieldCreatedAt:
		return "created_at"
	case FieldDeletedAt:
		return "deleted_at"
	case FieldEntityKind:
		return "entity_kind"
	case FieldTagName:
		return "tag"
	case FieldTitle:
		return "title"
	case FieldUpdatedAt:
		return "updated_at"
	default:
		return ""
	}
}
