package system

import "github.com/google/uuid"

type Void any
type Result[T any] struct {
	Value T
	Error Error
}

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}
func NewUUIDFromString(value string) (UUID, error) {
	u, err := uuid.Parse(value)
	return UUID(u), err
}
func (r *UUID) UnmarshalJSON(bytes []byte) (err error) {
	var u uuid.UUID
	if u, err = uuid.Parse(string(bytes)); err != nil {
		return
	}
	*r = UUID(u)
	return
}

func (r *UUID) MarshalJSON() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return []byte(uuid.UUID(*r).String()), nil
}

func (r *UUID) UnmarshalText(text []byte) (err error) {
	var u uuid.UUID
	if u, err = uuid.Parse(string(text)); err != nil {
		return
	}
	*r = UUID(u)
	return
}

func (r UUID) MarshalText() (text []byte, err error) {
	return []byte(uuid.UUID(r).String()), nil
}

func (r UUID) String() string {
	return uuid.UUID(r).String()
}
