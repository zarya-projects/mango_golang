package buisness

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	sl "main/internal/logger"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (m Message) Add(w http.ResponseWriter, r *http.Request) {

	m.DB.Cfg = m.Config

	// Создаем подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://" + m.Config.RABBIT_LOGIN + ":" + m.Config.RABBIT_PASS + "@" + m.Config.RABBIT_HOST + ":" + m.Config.RABBIT_PORT + "/")
	if err != nil {
		m.Logger.Error("unable to open connect to RabbitMQ server.", sl.Error(err))
	}
	// Закрываем подключение в случае удачной попытки
	defer func() {
		_ = conn.Close()
	}()

	if err := m.getData(r.FormValue("id")); err != nil {
		m.Logger.Error("failed to get request from server crm-zennit.ru.", sl.Error(err))
	}

	// Открытие нового канала
	ch, err := conn.Channel()
	if err != nil {
		m.Logger.Error("failed to open channel.", sl.Error(err))
	}
	// Закрываем канал в случае удачной попытки открытия
	defer func() {
		_ = ch.Close()
	}()

	m.Logger.Info("first value:", slog.Any("value", m.Val))

	if m.Val.Region == 178 {
		m.Val.City, m.Val.Region, err = m.DB.CheckCity(m.Val.CityName)
		if err != nil {
			m.Logger.Error("fatal in select reg and id from RUSSIA.", sl.Error(err))
		}
	}

	if m.Val.Stage == "NEW" &&
		m.Val.Phone != "" &&
		m.Val.City != 0 &&
		m.Val.Source != "CALL" &&
		m.Val.Source != "21" &&
		m.Val.Napr != 1700 &&
		m.Val.Region == 132 {
		if err := m.addFromRabbitMQ(ch); err != nil {
			m.Logger.Error("failed to add message in RabbitMQ. ", sl.Error(err))
		}
		m.Logger.Info("add value in RabbitMQ:", slog.Any("value", m.Val))

		go m.DB.Launch(m.Val.ID, m.Val.Phone, m.Val.City)
	}

}

// Получение данных с сервера Битрикс
func (m *Message) getData(id string) error {

	res, err := http.Get("https://crm-zennit.ru/local/mango/?get=1&apikey=4723adc21f8f06d7bd5f848438411161&lead=" + id)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	if err := json.Unmarshal(body, &m.Val); err != nil {
		return err
	}

	return nil
}

// Добавление нового сообщения в брокер
func (m Message) addFromRabbitMQ(ch *amqp.Channel) error {

	q, err := ch.QueueDeclare(
		m.Config.RABBIT_NAME, // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(m.Val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}
