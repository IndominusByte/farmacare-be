package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	handler_http "github.com/IndominusByte/farmacare-be/api/cmd/http/handler"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestValidationPokemons(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "empty",
			url:  "/pokemons",
		},
		{
			name: "type data",
			url:  "/pokemons" + "?page=a&per_page=a&order_by=a",
		},
		{
			name: "minimum",
			url:  "/pokemons" + "?page=-1&per_page=-1",
		},
		{
			name: "oneof",
			url:  "/pokemons" + "?page=1&per_page=1&order_by=a",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["per_page"].(string))
			case "type data":
				assert.Equal(t, "Invalid input type.", data["detail_message"].(map[string]interface{})["_body"].(string))
			case "minimum":
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["per_page"].(string))
			case "oneof":
				assert.Equal(t, "Must be one of: asc, desc.", data["detail_message"].(map[string]interface{})["order_by"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestPokemons(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	req, _ := http.NewRequest(http.MethodGet, "/pokemons?page=1&per_page=1&order_by=asc", nil)

	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	assert.NotNil(t, data["results"].(map[string]interface{})["data"])
	assert.Equal(t, 200, response.Result().StatusCode)
}

func TestBattle(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	req, _ := http.NewRequest(http.MethodGet, "/battle?page=1&per_page=1", nil)

	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	assert.NotNil(t, data["results"])
	assert.Equal(t, 200, response.Result().StatusCode)
}

func TestValidationBattleHistory(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "empty",
			url:  "/battle-history",
		},
		{
			name: "type data",
			url:  "/battle-history" + "?page=a&per_page=a&start_datetime=a&end_datetime=a",
		},
		{
			name: "minimum",
			url:  "/battle-history" + "?page=-1&per_page=-1&start_datetime=a&end_datetime=a",
		},
		{
			name: "datetime",
			url:  "/battle-history" + "?page=1&per_page=1&start_datetime=a&end_datetime=a",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["per_page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["start_datetime"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["end_datetime"].(string))
			case "type data":
				assert.Equal(t, "Invalid input type.", data["detail_message"].(map[string]interface{})["_body"].(string))
			case "minimum":
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["per_page"].(string))
			case "datetime":
				assert.Equal(t, "Invalid date time.", data["detail_message"].(map[string]interface{})["start_datetime"].(string))
				assert.Equal(t, "Must be greater than or equal to start_datetime.", data["detail_message"].(map[string]interface{})["end_datetime"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}

}

func TestBattleHistory(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	req, _ := http.NewRequest(http.MethodGet, "/battle-history?page=1&per_page=1&start_datetime=2022-11-05 00:00:00&end_datetime=2023-11-05 00:00:00", nil)

	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	assert.NotNil(t, data["results"].(map[string]interface{})["data"])
	assert.Equal(t, 200, response.Result().StatusCode)
}

func TestValidationPokemonCheating(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "required",
			payload: map[string]string{"name": ""},
		},
		{
			name:    "minimum",
			payload: map[string]string{"name": "a"},
		},
		{
			name:    "maximum",
			payload: map[string]string{"name": createMaximum(200)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/pokemon-cheating/asd", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["name"].(string))
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 3.", data["detail_message"].(map[string]interface{})["name"].(string))
			case "maximum":
				assert.Equal(t, "Longer than maximum length 100.", data["detail_message"].(map[string]interface{})["name"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestPokemonCheating(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}

	// get one battles history
	var poke handler_http.PokemonOut
	s.Collection = s.Db.Collection("battles")
	s.Collection.FindOne(s.Ctx, bson.M{}).Decode(&poke)

	tests := [...]struct {
		name       string
		url        string
		payload    map[string]string
		expected   string
		statusCode int
	}{
		{
			name:       "battle not found",
			url:        "/pokemon-cheating/asd",
			payload:    map[string]string{"name": "asd"},
			expected:   "Battle not found.",
			statusCode: 404,
		},
		{
			name:       "battle not found",
			url:        "/pokemon-cheating/" + poke.Id.Hex(),
			payload:    map[string]string{"name": "asdd"},
			expected:   "Battle not found.",
			statusCode: 404,
		},
		{
			name:       "success",
			url:        "/pokemon-cheating/" + poke.Id.Hex(),
			payload:    map[string]string{"name": poke.PokemonScore[0].Name},
			expected:   "",
			statusCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, test.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "success":
				assert.NotNil(t, data["results"])
			default:
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_body"].(string))
			}
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}

}
