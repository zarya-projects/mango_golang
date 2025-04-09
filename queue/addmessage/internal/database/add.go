package database

import (
	"fmt"
	"time"
)

func (s StructDatabase) Launch(id, phone string, city int) error {

	var err error

	s.BaseOpen, err = s.DBJoin()
	if err != nil {
		return err
	}

	if err = s.addData(id, phone, city); err != nil {
		return err
	}
	if err = s.addChange(id, phone, city); err != nil {
		return err
	}
	// Закрытие соединени с БД
	s.DBClose()
	return nil
}

func (s StructDatabase) addData(id, phone string, city int) error {

	q := `INSERT INTO mango_queue 
		  SET crm_id = ?,
		  	  phone = ?,
		  	  time = ?,
		  	  city = ?`

	res, err := s.BaseOpen.Exec(q, id, phone, time.Now().Unix(), city)
	if err != nil {
		fmt.Println(err)
	}

	if _, err = res.LastInsertId(); err != nil {
		return err
	}

	return nil
}

func (s StructDatabase) addChange(id, phone string, city int) error {

	q := `INSERT INTO changes_queue 
		  SET ` + "`change`" + ` = "ADD", 
		      crm_id = ?,
		  	  phone = ?,
		  	  time = ?,
		  	  city = ?`

	res, err := s.BaseOpen.Exec(q, id, phone, time.Now().Unix(), city)
	if err != nil {
		fmt.Println(err)
	}

	if _, err = res.LastInsertId(); err != nil {
		return err
	}

	return nil
}
