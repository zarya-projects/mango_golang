package database

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"
)

func (s *StructDatabase) Launch(city int, crm, phone string) (Value, error) {

	var (
		v   Value
		err error
	)

	// Открытие соединения с базой
	s.BaseOpen, err = s.DBJoin()
	if err != nil {
		return v, err
	}

	// Открываем новый канал для асинхронного запроса к БД
	c := make(chan Value)

	// Смотрим какой номер телефона у города
	groupID, err := s.checkGroupId(city)
	if err != nil {
		return v, err
	}

	s.Logger.Info("START", slog.Any("GROUP", groupID))

	// Пускаем 2 функции в горутину, чтобы получить номер телефона города
	// и свободного сотрудника
	go s.getOpenUsers(groupID, c)
	go s.getPhoneForCity(city, c)

	// Читаем канал
	for {
		val, ok := <-c
		if !ok {
			break
		}
		if val.Extension != 0 {
			v.Extension = val.Extension
		}
		if val.PhoneCity != "" {
			v.PhoneCity = val.PhoneCity
		}
		if val.UserCrm != "" {
			v.UserCrm = val.UserCrm
		}
	}

	// Удаление записей по изменениям брокера
	go s.deleteQueue(crm)
	go s.deleteQueueChanges(crm, phone, city)

	return v, nil
}

// Смотрим приоритет у города
func (s StructDatabase) checkGroupId(city int) (int, error) {

	var group int

	q := `SELECT group_id FROM mango_priority WHERE city_id = ?`

	res, err := s.BaseOpen.Query(q, city)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	for res.Next() {
		if err := res.Scan(&group); err != nil {
			return 0, err
		}
	}

	// Если нет привязки у группы к городу, то ставим 1 группу приоритетом
	if group == 0 {
		group = 1
	}

	return group, nil
}

func (s StructDatabase) getOpenUsers(groupID int, c chan Value) error {

	var (
		val Value
		err error
	)

	// Структура работы очереди по группам
	groupStruct := map[int]int{
		1: 2,
		2: 3,
		3: 1,
	}

	// Бегаем по бд в поисках свободно юзера начиная с приоритета
	for val.Extension == 0 {
		log.Printf("Должен взять в работу %s \n", strconv.Itoa(groupID))
		val.Extension, val.UserCrm, err = s.repeatedCheckQueue(groupID)
		if err != nil {
			close(c)
			return err
		}
		groupID = groupStruct[groupID]
		if val.Extension == 0 {
			time.Sleep(3 * time.Second)
		}
	}

	c <- val

	// запускаем в горутину обновление постобработки юзера
	// чтобы не стопать основной поток
	go s.updatePostUser(val.Extension)
	close(c)
	return nil
}

func (s StructDatabase) repeatedCheckQueue(groupID int) (int, string, error) {

	var (
		extension int
		crm       string
	)
	/**
	Логика запроса
	Запрос идет в 3 таблицы

	Забираем из базы всех свободны юзеров, которые сейчас не разговаривают,
	у которых с последнего звонка прошло больше 5 минут или они принудильно завершили постобработку
	и идет селект по группе
	Сортировка идет по тому, кто позже всего перестал разговаривать
	*/
	q := `SELECT mg.user_extension, mu.user_crm_id 
		  FROM mango_group AS mg, mango_priority AS mp, mango_users AS mu 
		  WHERE mp.group_id = mg.group_id AND 
		  mg.group_id = ? AND 
		  mg.user_extension = mu.extension AND 
		  (mu.call_status = "Disconnected" OR mu.call_status IS NULL) AND 
		  mu.mango_status = 1 AND 
		  (? - mu.timestamp > 300 || mu.post = "y")
		  ORDER BY mu.timestamp DESC`

	res, err := s.BaseOpen.Query(q, groupID, time.Now().Unix())
	if err != nil {
		return 0, "", err
	}
	defer res.Close()

	if res.Next() {
		if err := res.Scan(&extension, &crm); err != nil {
			return 0, "", err
		}
	}

	return extension, crm, nil
}

// Данный метод выполняется последний,
// поэтому в нем и будем закрывать соединение с БД
func (s StructDatabase) updatePostUser(extension int) {

	q := `UPDATE mango_users 
		  SET post = "n"
		  WHERE extension = ?`

	_, err := s.BaseOpen.Exec(q, extension)
	if err != nil {
		// Закрытие соединени с БД
		s.DBClose()
		fmt.Println(err)
	}

	// Закрытие соединени с БД
	s.DBClose()
}

// Канал один, поэтому этот метод канал завершить не может при успешном выполении
func (s StructDatabase) getPhoneForCity(city int, c chan Value) error {

	var val Value

	q := `SELECT phone FROM city WHERE id = ?`

	res, err := s.BaseOpen.Query(q, city)
	if err != nil {
		close(c)
		return err
	}
	defer res.Close()

	for res.Next() {
		if err := res.Scan(&val.PhoneCity); err != nil {
			close(c)
			return err
		}
	}

	c <- val
	return nil
}
