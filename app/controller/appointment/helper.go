package appointment

func checkFacilityAdmin(user map[string]interface{}, facilityID string) bool {
	var isValid = false
	var isAdmin = false
	var facility string
	var types []string
	var ok bool

	if facility, ok = user["facility"].(string); !ok {
		return false
	}

	if types, ok = user["types"].([]string); !ok {
		return false
	}

	for _, v := range types {
		if v == "Facility Admin" {
			isAdmin = true
		}
	}

	if facility == facilityID && isAdmin {
		isValid = true
	}

	return isValid
}
