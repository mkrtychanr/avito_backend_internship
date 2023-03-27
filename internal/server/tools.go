package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mkrtychanr/avito_backend_internship/internal/model"
	"github.com/shopspring/decimal"
)

func respondSuccess(w http.ResponseWriter) {
	makeJsonRespond(w, 200, jsonResult("success"))
}

func internalServerError(w http.ResponseWriter, err error) {
	makeJsonRespond(w, 500, jsonError("internal server error"))
	log.Println(err)
}

func invalidData(w http.ResponseWriter) {
	makeJsonRespond(w, 400, jsonError("invalid data"))
}

func clientNotFound(w http.ResponseWriter) {
	makeJsonRespond(w, 403, jsonError("client not found"))
}

func transactionNotFound(w http.ResponseWriter) {
	makeJsonRespond(w, 403, jsonError("transaction not found"))
}

func notEnoughMoneyOnTheBalanceSheet(w http.ResponseWriter) {
	makeJsonRespond(w, 406, jsonError("not enough money on the balance sheet"))
}

func reportNotFound(w http.ResponseWriter) {
	makeJsonRespond(w, 404, jsonError("report not found"))
}

func respondLink(w http.ResponseWriter, filename, address string) {
	makeJsonRespond(w, 200, jsonResult(fmt.Sprintf("http://%s/reports/%s", address, filename)))
}

func getDataFromRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w, err)
		return false
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		invalidData(w)
		return false
	}
	return true
}

func jsonError(respondText string) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, respondText))
}

func jsonResult(respondText string) []byte {
	return []byte(fmt.Sprintf(`{"result": "%s"}`, respondText))
}

func makeJsonRespond(w http.ResponseWriter, code int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func isGreaterThanOrEqualThanZero(value decimal.Decimal) bool {
	zero := decimal.New(0, 0)
	return value.GreaterThanOrEqual(zero)
}

func validateDate(date *model.Report) bool {
	dateSlice := strings.Split(date.Date, "-")
	if len(dateSlice) != 2 {
		return false
	}
	_, err := strconv.Atoi(dateSlice[0])
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(dateSlice[1])
	if err != nil {
		return false
	}
	if !(month > 0 && month < 13) {
		return false
	}
	date.Year = dateSlice[0]
	date.Month = dateSlice[1]
	return true
}

func createCSVReport(filename string, reportSlice []model.CSVReport) error {
	file, err := os.Create("reports/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'
	err = writer.WriteAll(reformSliceToCSVForm(reportSlice))
	if err != nil {
		return err
	}
	return nil
}

func reformSliceToCSVForm(data []model.CSVReport) [][]string {
	result := make([][]string, 0)
	for i, line := range data {
		result = append(result, make([]string, 2))
		result[i][0] = line.Service
		result[i][1] = line.Price
	}
	return result
}
