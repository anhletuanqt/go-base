package facilities

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func buildGetAllQuery(params getAllParams) map[string]interface{} {
	where := make(map[string]interface{})

	if params.Name != "" {
		where["name"] = map[string]interface{}{
			"$regex":   params.Name,
			"$options": "i",
		}
	}
	if params.Phone != "" {
		where["phone"] = map[string]interface{}{
			"$regex":   params.Phone,
			"$options": "i",
		}
	}
	if params.Status != "" {
		where["status"] = params.Status
	}
	if params.Cursor != "" {
		where["$and"] = bson.A{
			bson.D{{"name", bson.D{{"$gt", params.Cursor}}}},
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

	opts.SetSort(bson.D{{"name", order}}).SetLimit(int64(limit))

	return opts
}
