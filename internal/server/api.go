package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mkrtychanr/avito_backend_internship/internal/model"
	"github.com/shopspring/decimal"
)

func (s *Server) AddMoney(w http.ResponseWriter, r *http.Request) {
	addMoney := model.AddMoney{}
	ok := getDataFromRequest(w, r, &addMoney)
	if !ok {
		return
	}
	valueFromJson, err := decimal.NewFromString(addMoney.Value)
	if err != nil {
		invalidData(w)
		return
	}
	if !isGreaterThanOrEqualThanZero(valueFromJson) {
		invalidData(w)
		return
	}
	ok, err = s.isClientExist(addMoney.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		err = s.createClient(addMoney.Id)
		if err != nil {
			internalServerError(w, err)
			return
		}
	}
	actualValueString, err := s.getClientBalance(addMoney.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	actualValue, err := decimal.NewFromString(actualValueString)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.setBalance(addMoney.Id, actualValue.Add(valueFromJson).String())
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.addClientSheetChange(addMoney.Id, addMoney.Value, 0)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondSuccess(w)
}

func (s *Server) ReserveMoney(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	ok := getDataFromRequest(w, r, &transaction)
	if !ok {
		return
	}
	valueFromJson, err := decimal.NewFromString(transaction.Price)
	if err != nil {
		invalidData(w)
		return
	}
	if !isGreaterThanOrEqualThanZero(valueFromJson) {
		invalidData(w)
		return
	}
	ok, err = s.isClientExist(transaction.ClientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		clientNotFound(w)
		return
	}
	actualValueString, err := s.getClientBalance(transaction.ClientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	actualValue, err := decimal.NewFromString(actualValueString)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if actualValue.Cmp(valueFromJson) == -1 {
		notEnoughMoneyOnTheBalanceSheet(w)
		return
	}
	err = s.createNewReserve(transaction)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.setBalance(transaction.ClientId, actualValue.Sub(valueFromJson).String())
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.addClientSheetChange(transaction.ClientId, transaction.Price, 1)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondSuccess(w)
}

func (s *Server) AllowTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	ok := getDataFromRequest(w, r, &transaction)
	if !ok {
		return
	}
	id, ok, err := s.isTransactionInReserve(transaction)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		transactionNotFound(w)
		return
	}
	err = s.createNewReport(transaction)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.deleteReserve(id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondSuccess(w)
}

func (s *Server) GetClientBalance(w http.ResponseWriter, r *http.Request) {
	clientId := chi.URLParam(r, "client_id")
	fmt.Println("id: ", clientId)
	ok, err := s.isClientExist(clientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		clientNotFound(w)
		return
	}
	value, err := s.getClientBalance(clientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	makeJsonRespond(w, 200, jsonResult(value))
}

func (s *Server) UnreserveMoney(w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}
	ok := getDataFromRequest(w, r, &transaction)
	if !ok {
		return
	}
	valueFromJson, err := decimal.NewFromString(transaction.Price)
	if err != nil {
		invalidData(w)
		return
	}
	if !isGreaterThanOrEqualThanZero(valueFromJson) {
		invalidData(w)
		return
	}
	id, ok, err := s.isTransactionInReserve(transaction)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		transactionNotFound(w)
		return
	}
	actualValueString, err := s.getClientBalance(transaction.ClientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	actualValue, err := decimal.NewFromString(actualValueString)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.deleteReserve(id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.setBalance(transaction.ClientId, actualValue.Add(valueFromJson).String())
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.addClientSheetChange(transaction.ClientId, transaction.Price, 0)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondSuccess(w)
}

func (s *Server) GenerateReport(w http.ResponseWriter, r *http.Request) {
	date := model.Report{}
	ok := getDataFromRequest(w, r, &date)
	if !ok {
		return
	}
	if !validateDate(&date) {
		invalidData(w)
		return
	}
	reportSlice, err := s.getReportSlice(date)
	if err != nil {
		internalServerError(w, err)
		return
	}
	reportFileName := strings.Replace(time.Now().String(), " ", "_", -1) + ".csv"
	err = createCSVReport(reportFileName, reportSlice)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondLink(w, reportFileName, s.GetAddres())
}

func (s *Server) GetReport(w http.ResponseWriter, r *http.Request) {
	filename := "reports/" + chi.URLParam(r, "file")
	file, err := os.Open(filename)
	if err != nil {
		reportNotFound(w)
		return
	}
	defer file.Close()
	header := make([]byte, 512)
	file.Read(header)
	fileContentType := http.DetectContentType(header)
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)
	file.Seek(0, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(w, file)
		wg.Done()
	}()
	wg.Wait()
}
