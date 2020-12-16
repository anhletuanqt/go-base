package appointment

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func checkFacilityAdmin(user map[string]interface{}, facilityID string) bool {
	var isValid = false
	var isAdmin = false
	var facility string
	var types []interface{}
	var ok bool

	if facility, ok = user["facility"].(string); !ok {
		return false
	}

	if types, ok = user["types"].(primitive.A); !ok {
		fmt.Println("ok: ", ok)

		return false
	}

	for _, v := range types {
		v = v.(string)
		if v == "Facility Admin" {
			isAdmin = true
		}
	}

	if facility == facilityID && isAdmin {
		isValid = true
	}

	return isValid
}
