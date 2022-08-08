package util

import (
	"encoding/json"
	"fmt"
	"math"

	proto "google.golang.org/protobuf/proto"
)

func FormatRedisKeyForItem(itemId int64, shopId int64) string {
	return fmt.Sprintf("%d:%d", itemId, shopId)
}

func MarshalJson(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func MarshalProto(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

func UnmarshalProto(bytes []byte, message proto.Message) error {
	return proto.Unmarshal(bytes, message)
}

func CalculateNumberOfPages(totalCount int, itemsPerPage int) float64 {
	return math.Ceil(float64(totalCount) / float64(itemsPerPage))
}
