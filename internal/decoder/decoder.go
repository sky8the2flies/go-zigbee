package decoder

type Decoder struct {
	msgBuf    []byte
	msgBufIdx int
}

func NewDecoder() *Decoder {
	bufferSize := 1024
	return &Decoder{msgBuf: make([]byte, bufferSize), msgBufIdx: 0}
}

const (
	END = 192
)

func (d *Decoder) Decode(data []byte) (bool, []byte) {
	for _, b := range data {
		if b == END && d.msgBufIdx > 0 {
			msg := d.msgBuf[:d.msgBufIdx]
			d.msgBufIdx = 0
			return true, msg
		} else if b == END {
			continue
		}
		d.msgBuf[d.msgBufIdx] = b
		d.msgBufIdx++
	}
	msg := d.msgBuf[:d.msgBufIdx]
	return false, msg
}
