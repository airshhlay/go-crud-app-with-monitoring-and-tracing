package util

import (
	"fmt"
	"math"

	proto "google.golang.org/protobuf/proto"
)

// FormatRedisKeyForItem returns a string in the form itemid:shopId
func FormatRedisKeyForItem(itemID int64, shopID int64) string {
	return fmt.Sprintf("%d:%d", itemID, shopID)
}

// MarshalProto marshals a protobuf message into bytes
func MarshalProto(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

// UnmarshalProto unmarshals bytes into the given protobuf message
func UnmarshalProto(bytes []byte, message proto.Message) error {
	return proto.Unmarshal(bytes, message)
}

// CalculateNumberOfPages calculates the number of pages needed to display
// <totalCount> items with <itemsPerPage> per page
func CalculateNumberOfPages(totalCount int, itemsPerPage int) float64 {
	return math.Ceil(float64(totalCount) / float64(itemsPerPage))
}
