package http

import (
	"encoding/json"
	"testing"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Age     int      `json:"age"`
	Country string   `json:"country"`
	Roles   []string `json:"roles"`
}

var sampleUser = User{
	ID:      123,
	Name:    "John Doe",
	Email:   "john@test.com",
	Age:     30,
	Country: "India",
	Roles:   []string{"admin", "editor"},
}

var sampleJSON = []byte(`{"id":123,"name":"John Doe","email":"john@test.com","age":30,"country":"India","roles":["admin","editor"]}`)

func BenchmarkJSON_Marshal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := json.Marshal(sampleUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSON_Unmarshal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var u User
		err := json.Unmarshal(sampleJSON, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}
