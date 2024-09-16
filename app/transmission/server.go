package transmission

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/globals"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

var VALID_API_VERSIONS = []uint16{0, 1, 2, 3, 4}

func HandleConn(conn net.Conn) {
	client := NewClient(conn)

	for {
		request, err := client.ReceiveRequest()
		if err != nil {
			fmt.Println("Error receiving request: ", err.Error())
			return
		}

		for i, version := range VALID_API_VERSIONS {
			if request.ApiVersion == version {
				break
			}

			if i == len(VALID_API_VERSIONS)-1 {
				client.Respond(
					utils.ErrorResponse(request, globals.UnsupportedVersion),
				)

				return
			}
		}

		// Check if request.ApiKey is supported at all
		if _, ok := commands.GetApiKeys()[request.ApiKey]; !ok {
			client.Respond(
				utils.ErrorResponse(request, globals.InvalidRequest),
			)

			return
		}

		// Check if request.ApiVersion is supported for the given API key
		if _, ok := commands.GetApiKeys()[request.ApiKey][request.ApiVersion]; !ok {
			client.Respond(
				utils.ErrorResponse(request, globals.InvalidRequest),
			)

			return
		}

		command := commands.GetApiKeys()[request.ApiKey][request.ApiVersion]
		response := command(request)
		client.Respond(response)
	}
}

func MakeServer(port int) {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	l, err := net.Listen("tcp", address)
	fmt.Printf("Listening on port %d\n", port)

	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
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
