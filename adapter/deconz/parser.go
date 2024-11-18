package adapter_deconz

const (
	END = 192
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Encode(inputBuf []byte) ([]byte, error) {
	return nil, nil
}

func (p *Parser) Decode(data []byte, msgBufCap int) (bool, []byte, error) {
	complete, msg := extractMessage(data, msgBufCap)
	return complete, msg, nil
}

func extractMessage(data []byte, msgBufCap int) (bool, []byte) {
	msgBuf := make([]byte, msgBufCap)
	msgBufIdx := 0

	for _, b := range data {
		if b == END && msgBufIdx > 0 {
			msg := msgBuf[:msgBufIdx]
			msgBufIdx = 0
			return true, msg
		} else if b == END {
			continue
		}
		msgBuf[msgBufIdx] = b
		msgBufIdx++
	}
	msg := msgBuf[:msgBufIdx]
	return false, msg
}
