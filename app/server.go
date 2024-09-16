package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func LogBytesInHex(message string, data []byte) {
	fmt.Print(message + " [")
	for i, b := range data {
		fmt.Printf("%x", b)
		if i < len(data)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")
}

func ReadRequestLength(conn net.Conn) (*int, error) {
	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	data := buf[:n]
	LogBytesInHex("Received bytes:", data)
	if n != 4 {
		return nil, fmt.Errorf("expected 4 bytes, got %d", n)
	}

	length := int(binary.BigEndian.Uint32(data))

	fmt.Println("Received length: ", length)

	return &length, nil
}

func ReadExactBytes(conn net.Conn, length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if n != length {
		return nil, fmt.Errorf("expected %d bytes, got %d", length, n)
	}

	return buf, nil
}

type Request struct {
	ApiKey        []byte
	ApiVersion    []byte
	CorrelationId []byte
	Body          []byte
}

func ReadRequest(conn net.Conn, request_length int) (*Request, error) {
	buf, err := ReadExactBytes(conn, request_length)
	if err != nil {
		return nil, err
	}

	LogBytesInHex("Received raw request:", buf)

	request_api_key := buf[:2]
	request_api_version := buf[2:4]
	correlation_id := buf[4:12]
	body := buf[12:]

	request := &Request{
		ApiKey:        request_api_key,
		ApiVersion:    request_api_version,
		CorrelationId: correlation_id,
		Body:          body,
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

		// sample response hardcoded
		response_length := make([]byte, 4)
		binary.BigEndian.PutUint32(response_length, uint32(4))
		LogBytesInHex("Sending response length:", response_length)
		conn.Write(response_length)

		LogBytesInHex("Sending response:", request.CorrelationId)
		conn.Write(request.CorrelationId)
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
