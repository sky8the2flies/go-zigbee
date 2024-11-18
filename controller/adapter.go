package controller

type Adapter interface {
	Start() error
	Read() error
	Close() error
}
