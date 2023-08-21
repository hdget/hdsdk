package svc

type Service interface {
	Initialize(server any, generators ...Generator) error
}
