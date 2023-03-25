package server

import (
	"net/http"
	"strconv"

	"github.com/mkrtychanr/avito_backend_internship/internal/model"
	"github.com/shopspring/decimal"
)

func (s *Server) AddMoney(w http.ResponseWriter, r *http.Request) {
	addMoney := model.AddMoney{}
	ok := getDataFromRequest(w, r, &addMoney)
	if !ok {
		return
	}
	if addMoney.Id < 1 || addMoney.Value < 0 {
		invalidData(w)
		return
	}
	ok, err := s.isClientExist(addMoney.Id)
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
	actualValue, err := s.getClientBalance(addMoney.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	actualValue += addMoney.Value
	err = s.setBalance(addMoney.Id, actualValue)
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
	if transaction.ClientId < 1 || transaction.ServiceId < 1 || transaction.OrderId < 1 {
		invalidData(w)
		return
	}
	ok, err := s.isClientExist(transaction.ClientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if !ok {
		clientNotFound(w)
		return
	}
	balance, err := s.getClientBalance(transaction.ClientId)
	if err != nil {
		internalServerError(w, err)
		return
	}
	highBalance := decimal.NewFromFloat(balance)
	highPrice := decimal.NewFromFloat(transaction.Price)
	if highBalance.Cmp(highPrice) == -1 {
		notEnoughMoneyOnTheBalanceSheet(w)
		return
	}
	err = s.createNewReserve(transaction)
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = s.setBalance(transaction.ClientId, balance-transaction.Price)
	if err != nil {
		internalServerError(w, err)
		return
	}
	respondSuccess(w)
}

func (s *Server) AllowTransaction(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetClientBalance(w http.ResponseWriter, r *http.Request) {
	getClientBalance := model.ClientBalance{}
	ok := getDataFromRequest(w, r, &getClientBalance)
	if !ok {
		return
	}
	if getClientBalance.Id < 1 {
		invalidData(w)
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
	balance, err := s.getClientBalance(getClientBalance.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	makeJsonRespond(w, 200, jsonResult(strconv.FormatFloat(balance, 'f', -1, 64)))
}
