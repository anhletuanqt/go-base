package test

import (
	"base/test/helper"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"regexp"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("User", func() {
	t := GinkgoT()
	// facilities := make([]helper.Facility, 0, 100)

	Context("Query -", func() {
		It("Get User By Id", func() {
			mockFacility := helper.SeedUser()
			var data map[string]interface{}
			// http.Request
			req := httptest.NewRequest("GET", BASE_URL+"/users/"+mockFacility.ID.Hex(), nil)
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)
			assert.Equal(t, nil, err)

			item := data["item"].(map[string]interface{})
			assert.Equal(t, mockFacility.ID.Hex(), item["_id"])
		})

		Context("Get All -", func() {
			BeforeEach(func() {
				helper.SeedMultipleUser(100)
			})

			It("Default", func() {
				var data map[string]interface{}
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/users", nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})
				// helper.FormatJson(items, true)
				assert.Equal(t, 10, len(items))

				for i, v := range items {
					if i != len(items)-1 {
						currentEmail := v.(map[string]interface{})["email"].(string)
						nextEmail := items[i+1].(map[string]interface{})["email"].(string)
						assert.Equal(t, true, nextEmail > currentEmail)
					}
				}
			})

			It("Search By Name", func() {
				var data map[string]interface{}
				// Mock facility
				mockUser := helper.SeedUser()
				filterName := mockUser.FullName[:2]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/users?fullName="+filterName, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcher, _ := regexp.MatchString("(?i)"+filterName, current["fullName"].(string))
					assert.Equal(t, true, matcher)

					if i != len(items)-1 {
						currentEmail := current["email"].(string)
						nextEmail := items[i+1].(map[string]interface{})["email"].(string)
						assert.Equal(t, true, nextEmail > currentEmail)
					}
				}
			})

			It("Search By Email", func() {
				var data map[string]interface{}
				// Mock facility
				mockUser := helper.SeedUser()
				filterEmail := mockUser.Email[:2]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/users?email="+filterEmail, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcher, _ := regexp.MatchString("(?i)"+filterEmail, current["email"].(string))
					assert.Equal(t, true, matcher)

					if i != len(items)-1 {
						currentEmail := current["email"].(string)
						nextEmail := items[i+1].(map[string]interface{})["email"].(string)
						assert.Equal(t, true, nextEmail > currentEmail)
					}
				}
			})

			It("Combine Search By Email and By Name", func() {
				var data map[string]interface{}
				// Mock facility
				mockUser := helper.SeedUser()
				filterName := mockUser.FullName[1:3]
				filterEmail := mockUser.Email[:2]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/users?fullName="+filterName+"&email="+filterEmail, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcherEmail, _ := regexp.MatchString("(?i)"+filterEmail, current["email"].(string))
					matcherName, _ := regexp.MatchString("(?i)"+filterName, current["fullName"].(string))
					assert.Equal(t, true, matcherEmail)
					assert.Equal(t, true, matcherName)

					if i != len(items)-1 {
						currentEmail := current["email"].(string)
						nextEmail := items[i+1].(map[string]interface{})["email"].(string)
						assert.Equal(t, true, nextEmail > currentEmail)
					}
				}
			})

			It("Paging", func() {
				var data map[string]interface{}
				// http.Request
				pReq := httptest.NewRequest("GET", BASE_URL+"/users", nil)
				pReq.Header.Set("content-type", "application/json")

				// http.Response
				pResp, _ := App.Test(pReq)
				json.NewDecoder(pResp.Body).Decode(&data)
				pItems := data["items"].([]interface{})

				// Next page
				// http.Request
				cursor := pItems[9].(map[string]interface{})
				cursorEmail := cursor["email"].(string)

				req := httptest.NewRequest("GET", BASE_URL+"/users?cursor="+cursorEmail, nil)
				req.Header.Set("content-type", "application/json")
				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)

				items := data["items"].([]interface{})
				assert.Equal(t, 10, len(items))

				for i, v := range items {
					current := v.(map[string]interface{})
					currentEmail := current["email"].(string)
					if i == 0 {
						assert.Equal(t, true, currentEmail > cursorEmail)
					} else {
						prevEmail := items[i-1].(map[string]interface{})["email"].(string)
						assert.Equal(t, true, currentEmail > prevEmail)
					}
				}
			})
		})
	})

	Context("Mutation", func() {
		It("Create a user", func() {
			var data map[string]interface{}
			input := helper.SeedUserInput()
			bsInput, _ := helper.FormatJson(input, false)
			// http.Request
			req := httptest.NewRequest("POST", BASE_URL+"/users", bytes.NewReader(bsInput))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)
			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			expectTypes := input["types"].([]string)
			actualTypes := item["types"].([]interface{})

			assert.Equal(t, input["fullName"].(string), item["fullName"].(string))
			assert.Equal(t, input["email"].(string), item["email"].(string))
			assert.Equal(t, input["gender"].(string), item["gender"].(string))
			assert.Equal(t, input["dob"].(string), item["dob"].(string))
			assert.Equal(t, len(expectTypes), len(actualTypes))
			assert.Equal(t, expectTypes[0], actualTypes[0])
		})

		It("Update a user", func() {
			var data map[string]interface{}
			mockData := helper.SeedUser()
			input := helper.SeedUserInput()
			bsInput, _ := helper.FormatJson(input, false)
			// http.Request
			req := httptest.NewRequest("PUT", BASE_URL+"/users/"+mockData.ID.Hex(), bytes.NewReader(bsInput))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)

			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			expectTypes := input["types"].([]string)
			actualTypes := item["types"].([]interface{})

			assert.Equal(t, input["fullName"].(string), item["fullName"].(string))
			assert.Equal(t, input["email"].(string), item["email"].(string))
			assert.Equal(t, input["gender"].(string), item["gender"].(string))
			assert.Equal(t, input["dob"].(string), item["dob"].(string))
			assert.Equal(t, len(expectTypes), len(actualTypes))
			assert.Equal(t, expectTypes[0], actualTypes[0])
		})

		It("Delete a user", func() {
			var data map[string]interface{}
			mockData := helper.SeedUser()
			// http.Request
			req := httptest.NewRequest("DELETE", BASE_URL+"/users/"+mockData.ID.Hex(), nil)
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)

			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			assert.Equal(t, item["_id"].(string), mockData.ID.Hex())
		})
	})
})
