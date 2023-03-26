package server

import (
	"net/http"

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
	getClientBalance := model.ClientBalance{}
	ok := getDataFromRequest(w, r, &getClientBalance)
	if !ok {
		return
	}
	ok, err := s.isClientExist(getClientBalance.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		clientNotFound(w)
		return
	}
	value, err := s.getClientBalance(getClientBalance.Id)
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
	respondSuccess(w)
}
