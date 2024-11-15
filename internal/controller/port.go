package controller

type Port interface {
	Open() error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}
