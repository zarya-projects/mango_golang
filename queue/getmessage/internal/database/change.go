package database

import (
	sl "main/internal/logs"
	"time"
)

func (s StructDatabase) deleteQueue(crm string) error {

	q := `DELETE FROM mango_queue WHERE crm_id = ?`

	_, err := s.BaseOpen.Exec(q, crm)
	if err != nil {
		s.Logger.Error("err in mango_queue.", sl.Error(err))
	}
	s.Logger.Info("DELETE FROM mango_queue IS GOOD")
	return nil
}

func (s StructDatabase) deleteQueueChanges(crm, phone string, city int) error {

	q := `INSERT INTO changes_queue 
		  SET ` + "`change`" + ` = "DELETE", 
		      crm_id = ?,
		  	  phone = ?,
		  	  time = ?,
		  	  city = ?`

	_, err := s.BaseOpen.Exec(q, crm, phone, time.Now().Unix(), city)
	if err != nil {
		s.Logger.Error("err in changes_queue.", sl.Error(err))
	}
	s.Logger.Info("DELETE FROM mango_queue IS GOOD")
	return nil
}
