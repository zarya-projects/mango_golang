package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

type product struct {
	crm    string
	phone  string
	city   string
	id     string
	change string
	time   string
}

// сокет
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	// подключение к БД
	db, err := sql.Open("mysql", "admin:2024T,bnt,f,yfcdt;tvctyt@/mango")

	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// обновление соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	// для проверок изменений в БД
	count := 0
	newBool := false

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// цикл обработки сообщений
	for {
		messageType, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var dataMap []map[string]string

		// эхо ансвер
		for i := 0; i >= 0; i++ {

			dataMap, count, newBool = dbGet(db, count, newBool)
			data, err := json.Marshal(dataMap)
			if err != nil {
				log.Println(err)
				break
			}

			// обработка результатов запроса в бд и ответ клиенту
			time.Sleep(5 * time.Second)

			if len(dataMap) > 0 {
				if err := ws.WriteMessage(messageType, data); err != nil {
					log.Println(err)
					break
				}
			} else {
				if err := ws.WriteMessage(messageType, nil); err != nil {
					log.Println(err)
					break
				}
			}
		}
	}
}

// конец сокета

// DB поиск
func dbGet(db *sql.DB, count int, newBool bool) ([]map[string]string, int, bool) {

	var newCount int
	selectCount, err := db.Query("SELECT COUNT(id) FROM mango.changes_queue")

	if err != nil {
		log.Println(err)
	}
	defer selectCount.Close()

	for selectCount.Next() {
		if err := selectCount.Scan(&newCount); err != nil {
			log.Println(err)
		}
	}

	// если запрос первый, то вызывается функция selectedAll
	if newBool == false {
		newBool = true
		return selectedAll(db, count), newCount, newBool
	} else if newCount != count {
		data := selectedChange(db, newCount-count)
		return data, newCount, newBool
	}

	return []map[string]string{}, newCount, newBool
}

// поиск в таблице изменений статутов
func selectedChange(db *sql.DB, count int) []map[string]string {

	rows, err := db.Query("SELECT * FROM mango.changes_queue ORDER BY id DESC LIMIT ?", count)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	products := []product{}

	for rows.Next() {
		p := product{}
		err := rows.Scan(&p.crm, &p.city, &p.change, &p.phone, &p.id, &p.time)

		if err != nil {
			log.Println(err)
			continue
		}

		products = append(products, p)

	}

	data := []map[string]string{}
	log.Println(count)
	for _, p := range products {
		log.Println(p)
		longData := map[string]string{"change": p.crm, "city": p.phone, "crm": p.id, "id": p.city, "phone": p.change, "time": p.time}
		data = append(data, longData)
	}

	return data

}

// данная функция вызывается при первом соединении в сокетом и запрашивает данные из таблицы со всему звонками в статусе ожидания
func selectedAll(db *sql.DB, count int) []map[string]string {

	rows, err := db.Query("SELECT * FROM mango.mango_queue WHERE id >= ?", count)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	products := []product{}

	for rows.Next() {
		p := product{}
		err := rows.Scan(&p.crm, &p.phone, &p.city, &p.id, &p.time)

		if err != nil {
			log.Println(err)
			continue
		}

		products = append(products, p)

	}

	data := []map[string]string{}

	for _, p := range products {
		log.Println(p)
		longData := map[string]string{"id": p.crm, "phone": p.phone, "city": p.city, "time": p.time}
		data = append(data, longData)
	}

	return data

}
