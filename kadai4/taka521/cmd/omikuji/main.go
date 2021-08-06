package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/taka521/gopherdojo-studyroom/kadai4/taka521/omikuji"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	http.HandleFunc("/", omikuji.Handler)
	log.Printf("🚀 start server - http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
