package http_server

type Server interface {
	Init()
	Run()
	Close() error
}
