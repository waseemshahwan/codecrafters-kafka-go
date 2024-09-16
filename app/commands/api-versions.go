package commands

import (
	"encoding/binary"
	"sort"

	"github.com/codecrafters-io/redis-starter-go/app/globals"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

func ApiVersionsV4(request *globals.Request) *globals.Response {
	api_keys := GetApiKeys()

	total_api_keys_size := 1
	api_key_size := 2
	minimum_version_size := 2
	maximum_version_size := 2
	tagged_fields_size := 1
	api_key_item_size := api_key_size + minimum_version_size + maximum_version_size + tagged_fields_size
	throttle_time_size := 4

	response_payload := make(
		[]byte,
		total_api_keys_size+
			len(api_keys)*api_key_item_size+
			throttle_time_size+
			tagged_fields_size,
	)

	response_payload[0] = 1 + uint8(len(api_keys))

	index := 0
	for api_key, versions_map := range api_keys {
		versions := make([]int, len(versions_map))
		i := 0
		for version := range versions_map {
			versions[i] = int(version)
			i++
		}

		// Assumption: all versions will be contiguous
		sort.Ints(versions)
		minimum_version := versions[0]
		maximum_version := versions[len(versions)-1]

		payload := make([]byte, 7)
		binary.BigEndian.PutUint16(payload[0:2], uint16(api_key))
		binary.BigEndian.PutUint16(payload[2:4], uint16(minimum_version))
		binary.BigEndian.PutUint16(payload[4:6], uint16(maximum_version))
		// tagged_fields
		payload[6] = 0

		offset := total_api_keys_size
		copy(response_payload[index*api_key_item_size+offset:], payload)
		index++
	}

	copy(response_payload[len(api_keys)*api_key_item_size+total_api_keys_size:], []byte{0, 0, 0, 0})
	response_payload[len(api_keys)*api_key_item_size+total_api_keys_size+throttle_time_size] = 0

	return utils.SuccessResponse(request, response_payload)
}
