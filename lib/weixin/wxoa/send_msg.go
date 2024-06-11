package wxoa

type SendMessager interface {
	Send(contents map[string]string) error
}
