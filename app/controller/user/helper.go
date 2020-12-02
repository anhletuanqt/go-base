package user

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func buildGetAllQuery(params getAllParams) map[string]interface{} {
	where := make(map[string]interface{})

	if params.FullName != "" {
		where["fullName"] = map[string]interface{}{
			"$regex":   params.FullName,
			"$options": "i",
		}
	}
	if params.Email != "" {
		where["email"] = map[string]interface{}{
			"$regex":   params.Email,
			"$options": "i",
		}
	}
	if params.Gender != "" {
		where["gender"] = params.Gender
	}
	if params.Cursor != "" {
		where["$and"] = bson.A{
			bson.D{{"email", bson.D{{"$gt", params.Cursor}}}},
		}
	}

	return where
}

func buildGetAllOpts(params getAllParams) *options.FindOptions {
	opts := options.Find()

	order := 1
	limit := 10

	if params.Order == "desc" {
		order = -1
	}
	if params.Limit != 0 {
		limit = params.Limit
	}

	opts.SetSort(bson.D{{"email", order}}).SetLimit(int64(limit))

	return opts
}
