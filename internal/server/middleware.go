package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Логирующее middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		logEntry := map[string]interface{}{
			"time":         start.Format(time.RFC3339),
			"method":       r.Method,
			"url":          r.URL.Path,
			"query_params": r.URL.Query(),
			"request_body": string(body),
		}

		// Обертка для записи статуса и ответа
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, responseBody: new(bytes.Buffer)}

		// Вызов следующего обработчика
		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logEntry["duration"] = duration.String()
		logEntry["status_code"] = rec.statusCode
		logEntry["response_body"] = rec.responseBody.String()

		logJSON, _ := json.Marshal(logEntry)
		log.Println(string(logJSON))
	})
}

// Пример основного обработчика
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Hello, world!"}
	json.NewEncoder(w).Encode(response)
}

// Обертка для записи статуса и ответа
type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	responseBody *bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.responseBody.Write(b)
	return rec.ResponseWriter.Write(b)
}

// recoverMiddleware для перехвата паники и логирования ошибки
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
