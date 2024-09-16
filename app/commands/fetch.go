package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/globals"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func FetchV16(request *globals.Request) *globals.Response {
	return utils.SuccessResponse(request, []byte{})
}
