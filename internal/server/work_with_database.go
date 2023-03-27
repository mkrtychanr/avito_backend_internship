package server

import (
	"fmt"

	"github.com/mkrtychanr/avito_backend_internship/internal/model"
)

func (s *Server) createClient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("insert into client (id, balance) values ($1, 0)", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) isClientExist(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query("select id from client where id = $1", id)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s *Server) getClientBalance(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query("select balance from client where id = $1", id)
	if err != nil {
		return "", err
	}
	rows.Next()
	var result string
	rows.Scan(&result)
	return result, nil
}

func (s *Server) setBalance(id string, balance string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("update client set balance = $1 where id = $2", balance, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) createNewReserve(transaction model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("insert into reserve (client_id, service_id, order_id, price, reserve_time) values ($1, $2, $3, $4, current_timestamp)", transaction.ClientId, transaction.ServiceId, transaction.OrderId, transaction.Price)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) deleteReserve(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("delete from reserve where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) createNewReport(transaction model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("insert into report (client_id, service_id, order_id, price, report_time) values ($1, $2, $3, $4, current_timestamp)", transaction.ClientId, transaction.ServiceId, transaction.OrderId, transaction.Price)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) isTransactionInReserve(transaction model.Transaction) (int64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query("select id from reserve where client_id=$1 and service_id=$2 and order_id=$3 and price=$4", transaction.ClientId, transaction.ServiceId, transaction.OrderId, transaction.Price)
	if err != nil {
		return 0, false, err
	}
	if rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return 0, false, err
		}
		return id, true, nil
	}
	return 0, false, nil
}

func (s *Server) addClientSheetChange(id, value string, status byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec("insert into client_sheet_change (client_id, difference, status, change_time) values ($1, $2, $3, current_timestamp)", id, value, status)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) getReportSlice(date model.Report) ([]model.CSVReport, error) {
	s.mu.Lock()
	rows, err := s.db.Query(`with time_filter as (
		select * from report
		where 
		date_part('year', report_time) = $1 
		and 
		date_part('month', report_time) = $2
	  )
	  
	  select service_id, sum(price::decimal) from time_filter group by service_id`, date.Year, date.Month)
	s.mu.Unlock()
	if err != nil {
		return nil, err
	}
	data := make([]model.CSVReport, 0)
	for rows.Next() {
		newLine := model.CSVReport{}
		err := rows.Scan(&newLine.Service, &newLine.Price)
		if err != nil {
			return nil, err
		}
		fmt.Println(newLine.Service)
		data = append(data, newLine)
	}
	return data, nil
}
