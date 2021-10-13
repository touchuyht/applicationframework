package appentrance

type Government interface {
	Register(Department) error
	Start() []error
	GracefulShutdown() error
}

type Department interface {
	Name() string
	Start(ch chan struct{}) error
}