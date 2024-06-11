package wxoa

type SendMessager interface {
	Send(contents ...string) error
}
