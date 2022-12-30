package user

type User interface {
	GloballyIdentified
	ListDisplayable
	GetIdentityID() string
	GetEmail() string
}
