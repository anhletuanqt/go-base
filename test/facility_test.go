package test

import (
	"base/test/helper"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Facility", func() {
	t := GinkgoT()
	// facilities := make([]helper.Facility, 0, 100)

	Context("Query -", func() {
		It("Get Facility By Id", func() {
			mockFacility := helper.SeedFacility()
			var data map[string]interface{}
			// http.Request
			req := httptest.NewRequest("GET", BASE_URL+"/facilities/"+mockFacility.ID.Hex(), nil)
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
				helper.SeedMultipleFacility(100)
			})

			It("Default", func() {
				var data map[string]interface{}
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/facilities", nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})
				assert.Equal(t, 10, len(items))
				for i, v := range items {
					if i != len(items)-1 {
						currentName := v.(map[string]interface{})["name"].(string)
						nextName := items[i+1].(map[string]interface{})["name"].(string)
						assert.Equal(t, true, nextName > currentName)
					}
				}
			})

			It("Search By Name", func() {
				var data map[string]interface{}
				// Mock facility
				mockFacility := helper.SeedFacility()
				filterName := mockFacility.Name[:2]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/facilities?name="+filterName, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcher, _ := regexp.MatchString("(?i)"+filterName, current["name"].(string))
					assert.Equal(t, true, matcher)

					if i != len(items)-1 {
						currentName := current["name"].(string)
						nextName := items[i+1].(map[string]interface{})["name"].(string)
						assert.Equal(t, true, nextName > currentName)
					}
				}
			})

			It("Search By Phone", func() {
				var data map[string]interface{}
				// Mock facility
				mockFacility := helper.SeedFacility()
				filterPhone := mockFacility.Phone[1:3]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/facilities?phone="+filterPhone, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcher, _ := regexp.MatchString("(?i)"+filterPhone, current["phone"].(string))
					assert.Equal(t, true, matcher)

					if i != len(items)-1 {
						currentName := current["name"].(string)
						nextName := items[i+1].(map[string]interface{})["name"].(string)
						assert.Equal(t, true, nextName > currentName)
					}
				}
			})

			It("Combine Search By Phone and By Name", func() {
				var data map[string]interface{}
				// Mock facility
				mockFacility := helper.SeedFacility()
				filterPhone := mockFacility.Phone[1:3]
				filterName := mockFacility.Name[:2]
				// http.Request
				req := httptest.NewRequest("GET", BASE_URL+"/facilities?phone="+filterPhone+"&name="+filterName, nil)
				req.Header.Set("content-type", "application/json")

				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)
				items := data["items"].([]interface{})

				for i, v := range items {
					current := v.(map[string]interface{})
					matcherPhone, _ := regexp.MatchString("(?i)"+filterPhone, current["phone"].(string))
					matcherName, _ := regexp.MatchString("(?i)"+filterName, current["name"].(string))
					assert.Equal(t, true, matcherPhone)
					assert.Equal(t, true, matcherName)

					if i != len(items)-1 {
						currentName := current["name"].(string)
						nextName := items[i+1].(map[string]interface{})["name"].(string)
						assert.Equal(t, true, nextName > currentName)
					}
				}
			})

			It("Paging", func() {
				var data map[string]interface{}
				// http.Request
				pReq := httptest.NewRequest("GET", BASE_URL+"/facilities", nil)
				pReq.Header.Set("content-type", "application/json")

				// http.Response
				pResp, _ := App.Test(pReq)
				json.NewDecoder(pResp.Body).Decode(&data)
				pItems := data["items"].([]interface{})

				// Next page
				// http.Request
				cursor := pItems[9].(map[string]interface{})
				cursorName := cursor["name"].(string)

				req := httptest.NewRequest("GET", BASE_URL+"/facilities?cursor="+strings.Split(cursorName, " ")[1], nil)
				req.Header.Set("content-type", "application/json")
				// http.Response
				resp, _ := App.Test(req)
				err := json.NewDecoder(resp.Body).Decode(&data)
				assert.Equal(t, nil, err)

				items := data["items"].([]interface{})
				assert.Equal(t, 10, len(items))

				for i, v := range items {
					current := v.(map[string]interface{})
					currentName := current["name"].(string)
					if i == 0 {
						cursorName := cursor["name"].(string)
						assert.Equal(t, true, currentName > cursorName)
					} else {
						prevName := items[i-1].(map[string]interface{})["name"].(string)
						assert.Equal(t, true, currentName > prevName)
					}
				}
			})
		})
	})

	Context("Mutation", func() {
		It("Create a facility", func() {
			var data map[string]interface{}
			input := helper.SeedFacilityInput()
			bsInput, _ := helper.FormatJson(input, false)
			// http.Request
			req := httptest.NewRequest("POST", BASE_URL+"/facilities", bytes.NewReader(bsInput))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)
			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			assert.Equal(t, input["name"].(string), item["name"].(string))
			assert.Equal(t, input["phone"].(string), item["phone"].(string))
			assert.Equal(t, input["status"].(string), item["status"].(string))
		})

		It("Update a facility", func() {
			var data map[string]interface{}
			mockData := helper.SeedFacility()
			helper.FormatJson(mockData, false)
			input := helper.SeedFacilityInput()
			bsInput, _ := helper.FormatJson(input, false)
			// http.Request
			req := httptest.NewRequest("PUT", BASE_URL+"/facilities/"+mockData.ID.Hex(), bytes.NewReader(bsInput))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)

			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			assert.Equal(t, input["name"].(string), item["name"].(string))
			assert.Equal(t, input["phone"].(string), item["phone"].(string))
			assert.Equal(t, input["status"].(string), item["status"].(string))
		})

		It("Delete a facility", func() {
			var data map[string]interface{}
			mockData := helper.SeedFacility()
			// http.Request
			req := httptest.NewRequest("DELETE", BASE_URL+"/facilities/"+mockData.ID.Hex(), nil)
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
