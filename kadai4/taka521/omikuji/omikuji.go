package omikuji

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

var _ http.HandlerFunc = Handler

func Handler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	r := Execute(time.Now())

	b, err := json.Marshal(r)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Write(b)
}

type response struct {
	Result string
}

var Execute = func(t time.Time) *response {
	if t.Month() == time.January && t.Day() <= 3 {
		return &response{Result: "大吉"}
	}
	return &response{Result: kuji[rand.Intn(len(kuji))]}
}

var kuji = []string{"大吉", "吉", "中吉", "小吉", "凶"}
