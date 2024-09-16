package commands

import "github.com/codecrafters-io/redis-starter-go/app/globals"

type ApiKeys map[globals.ApiKey]map[uint16]func(*globals.Request) *globals.Response

func GetApiKeys() ApiKeys {
	return ApiKeys{
		globals.ApiVersions: {
			4: ApiVersionsV4,
		},
		globals.Fetch: {
			16: FetchV16,
		},
	}
}
