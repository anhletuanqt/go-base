package test

import (
	"base/test/helper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("Question Set", func() {
	t := GinkgoT()
	BeforeEach(func() {

	})

	Context("Query", func() {
		It("Get Question By Id", func() {
			item := helper.SeedQuestionSet()
			var bodyResp map[string]interface{}
			// http.Request
			req := httptest.NewRequest("GET", BASE_URL+"/qs/"+item.ID.Hex(), nil)
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&bodyResp)
			assert.Equal(t, nil, err)
			doc := bodyResp["item"].(map[string]interface{})
			assert.Equal(t, item.ID.Hex(), doc["_id"])
		})

		It("Get All", func() {
			helper.SeedMultipleQuestionSet(20)
			var bodyResp map[string]interface{}
			// http.Request
			req := httptest.NewRequest("GET", BASE_URL+"/qs/", nil)
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&bodyResp)
			assert.Equal(t, nil, err)
			docs := bodyResp["items"].([]interface{})
			count, _ := DB.Collection("questionsets").CountDocuments(context.Background(), bson.D{{}})
			assert.Equal(t, int(count), len(docs))
		})

		It("Get All - Search by name", func() {
			mockData := helper.SeedMultipleQuestionSet(20)
			filterName := mockData[0].Name[:2]
			var bodyResp map[string]interface{}
			// http.Request
			req := httptest.NewRequest("GET", BASE_URL+"/qs?name="+filterName, nil)
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&bodyResp)
			assert.Equal(t, nil, err)
			docs := bodyResp["items"].([]interface{})
			fmt.Println(filterName)

			for _, v := range docs {
				d := v.(map[string]interface{})
				matcher, _ := regexp.MatchString("(?i)"+filterName, d["name"].(string))
				assert.Equal(t, true, matcher)
			}
		})
	})

	Context("Mutation", func() {
		It("Create a a question set", func() {
			var data map[string]interface{}
			input := helper.SeedQuestionSetInput(true).(map[string]string)
			bsInput, _ := helper.FormatJson(input, true)
			// http.Request
			req := httptest.NewRequest("POST", BASE_URL+"/qs", bytes.NewReader(bsInput))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, _ := App.Test(req)
			err := json.NewDecoder(resp.Body).Decode(&data)
			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			assert.Equal(t, input["name"], item["name"])
		})

		It("Update a question set", func() {
			var data map[string]interface{}
			input := helper.SeedQuestionSet()
			updateData, _ := json.Marshal(struct{ Name string }{Name: "Updated data"})
			// http.Request
			req := httptest.NewRequest("PUT", BASE_URL+"/qs/"+input.ID.Hex(), bytes.NewReader(updateData))
			req.Header.Set("content-type", "application/json")

			// http.Response
			resp, err1 := App.Test(req)
			fmt.Println(err1)
			err := json.NewDecoder(resp.Body).Decode(&data)
			assert.Equal(t, nil, err)
			item := data["item"].(map[string]interface{})
			assert.Equal(t, "Updated data", item["name"])
		})
	})
})
