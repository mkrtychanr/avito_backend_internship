package server

import (
	"net/http"
	"strconv"

	"github.com/mkrtychanr/avito_backend_internship/internal/model"
)

func (s *Server) AddMoney(w http.ResponseWriter, r *http.Request) {
	addMoney := model.AddMoney{}
	ok := getDataFromRequest(w, r, &addMoney)
	if !ok {
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

}

func (s *Server) AllowTransaction(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetClientBalance(w http.ResponseWriter, r *http.Request) {
	getClientBalance := model.GetClientBalance{}
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
		makeJsonRespond(w, 403, jsonError("client not found"))
		return
	}
	balance, err := s.getClientBalance(getClientBalance.Id)
	if err != nil {
		internalServerError(w, err)
	}
	makeJsonRespond(w, 200, jsonResult(strconv.FormatFloat(balance, 'f', -1, 64)))
}
