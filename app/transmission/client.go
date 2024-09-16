package transmission

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/globals"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) ReceiveRequest() (*globals.Request, error) {
	length, err := ReadRequestLength(c.conn)
	if err != nil {
		return nil, err
	}

	request, err := ReadRequest(c.conn, *length)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (c *Client) Respond(response *globals.Response) {
	fmt.Println("Send ", len(response.Body), " bytes")
	utils.LogBytesInHex(response.Body)

	binary.Write(c.conn, binary.BigEndian, uint32(len(response.Body)+4))
	binary.Write(c.conn, binary.BigEndian, response.Request.CorrelationId)
	binary.Write(c.conn, binary.BigEndian, response.Body)
}
