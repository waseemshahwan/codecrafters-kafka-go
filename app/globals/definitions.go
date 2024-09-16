package globals

import "encoding/binary"

type Request struct {
	ApiKey        ApiKey
	ApiVersion    uint16
	CorrelationId uint32
	Body          []byte
}

type Response struct {
	Request *Request
	Body    []byte
}

type ErrorCode uint16

const (
	Success            ErrorCode = 0
	UnsupportedVersion ErrorCode = 35
	InvalidRequest     ErrorCode = 42
)

func (code ErrorCode) Bytes() []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(code))
	return bytes
}

type ApiKey uint16

const (
	ApiVersions ApiKey = 18
	Fetch       ApiKey = 1
)
