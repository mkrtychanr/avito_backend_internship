package server

func (s *Server) createClient(id int64) error {
	_, err := s.db.Exec("insert into client (id, balance) values ($1, 0)", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) isClientExist(id int64) (bool, error) {
	rows, err := s.db.Query("select id from client where id = $1", id)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s *Server) getClientBalance(id int64) (float64, error) {
	rows, err := s.db.Query("select balance from client where id = $1", id)
	if err != nil {
		return 0, err
	}
	rows.Next()
	var result float64
	rows.Scan(&result)
	return result, nil
}

func (s *Server) setBalance(id int64, balance float64) error {
	_, err := s.db.Exec("update client set balance = $1 where id = $2", balance, id)
	if err != nil {
		return err
	}
	return nil
}
