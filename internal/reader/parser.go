package reader

type Parser interface {
	Chunk(inputBuf []byte) (bool, []byte, error)
	Frame(chunk []byte) error
}
