package converters

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringArrToObjectID(strs []string) ([]primitive.ObjectID, error) {
	res := make([]primitive.ObjectID, 0, len(strs))
	for _, str := range strs {
		id, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil, fmt.Errorf("failed to convert \"%s\" into ObjectID: %w", id, err)
		}
		res = append(res, id)
	}
	return res, nil
}
