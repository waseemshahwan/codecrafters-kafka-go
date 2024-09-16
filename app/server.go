package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func LogBytesInHex(data []byte) {
	for i, b := range data {
		hex := "0" + fmt.Sprintf("%x", b)
		hex = hex[len(hex)-2:]

		fmt.Printf("%s ", hex)
		if i%8 == 7 && len(data)%8 != 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func ReadRequestLength(conn net.Conn) (*int, error) {
	buf, err := ReadExactBytes(conn, 4)
	if err != nil {
		return nil, err
	}

	length := int(binary.BigEndian.Uint32(buf))

	return &length, nil
}

func ReadExactBytes(conn net.Conn, length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	fmt.Println("Read ", n, " bytes")
	LogBytesInHex(buf)

	if n != length {
		return nil, fmt.Errorf("expected %d bytes, got %d", length, n)
	}

	return buf, nil
}

func SendBytes(conn net.Conn, bytes []byte) {
	fmt.Println("Send ", len(bytes), " bytes")
	LogBytesInHex(bytes)

	conn.Write(bytes)
}

type Request struct {
	ApiKey        []byte
	ApiVersion    []byte
	CorrelationId []byte
	Body          []byte
}

func ParseApiVersion(request Request) bool {
	VALID_API_VERSIONS := [][]byte{
		{0, 0, 0, 0},
		{0, 0, 0, 1},
		{0, 0, 0, 2},
		{0, 0, 0, 3},
		{0, 0, 0, 4},
	}

	for _, version := range VALID_API_VERSIONS {
		if bytes.Equal(request.ApiVersion, version) {
			return true
		}
	}

	return false
}

func ReadRequest(conn net.Conn, request_length int) (*Request, error) {
	buf, err := ReadExactBytes(conn, request_length)
	if err != nil {
		return nil, err
	}

	request_api_key := buf[:2]
	request_api_version := buf[2:4]
	correlation_id := buf[4:8]
	body := buf[8:]

	request := &Request{
		ApiKey:        request_api_key,
		ApiVersion:    request_api_version,
		CorrelationId: correlation_id,
		Body:          body,
	}

	if !ParseApiVersion(*request) {
		SendBytes(conn, []byte{0, 0, 0, 8})
		SendBytes(conn, request.CorrelationId)
		SendBytes(conn, []byte{0, 35})
		return nil, fmt.Errorf("invalid api version")
	}

	return request, nil
}

func HandleConn(conn net.Conn) {
	for {
		request_length, err := ReadRequestLength(conn)
		if err != nil {
			fmt.Println("Error reading request length: ", err.Error())
			conn.Close()
			return
		}

		request, err := ReadRequest(conn, *request_length)
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			conn.Close()
			return
		}

		// Response length of 4 bytes, for now it is ignored
		SendBytes(conn, []byte{0, 0, 0, 4})

		// Send correlation id back to client
		SendBytes(conn, request.CorrelationId)

		// Send an INVALID_VERSION error back to client
		SendBytes(conn, []byte{35})
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	fmt.Println("Listening on port 9092")

	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go HandleConn(conn)
	}
}
