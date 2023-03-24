package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func internalServerError(w http.ResponseWriter, err error) {
	makeJsonRespond(w, 500, jsonError("internal server error"))
	log.Println(err)
}

func respondSuccess(w http.ResponseWriter) {
	makeJsonRespond(w, 200, jsonResult("success"))
}

func getDataFromRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w, err)
		return false
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		makeJsonRespond(w, 400, jsonError("invalid data"))
		return false
	}
	return true
}

func jsonError(respondText string) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, respondText))
}

func jsonResult(respondText string) []byte {
	return []byte(fmt.Sprintf(`{"result": %s}`, respondText))
}

func makeJsonRespond(w http.ResponseWriter, code int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
