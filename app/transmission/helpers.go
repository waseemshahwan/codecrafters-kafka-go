package transmission

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/globals"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func ReadRequestLength(conn net.Conn) (*int, error) {
	buf, err := ReadExactBytes(conn, 4)
	if err != nil {
		return nil, err
	}

	length := int(binary.BigEndian.Uint32(buf))

	return &length, nil
}

func ReadRequest(conn net.Conn, request_length int) (*globals.Request, error) {
	buf, err := ReadExactBytes(conn, request_length)
	if err != nil {
		return nil, err
	}

	request_api_key := buf[:2]
	request_api_version := buf[2:4]
	correlation_id := buf[4:8]
	body := buf[8:]

	request := &globals.Request{
		ApiKey:        globals.ApiKey(binary.BigEndian.Uint16(request_api_key)),
		ApiVersion:    binary.BigEndian.Uint16(request_api_version),
		CorrelationId: binary.BigEndian.Uint32(correlation_id),
		Body:          body,
	}

	return request, nil
}

func ReadExactBytes(conn net.Conn, length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	fmt.Println("Read ", n, " bytes")
	utils.LogBytesInHex(buf)

	if n != length {
		return nil, fmt.Errorf("expected %d bytes, got %d", length, n)
	}

	return buf, nil
}
