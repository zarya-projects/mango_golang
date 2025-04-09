package logic

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"main/internal/config"
	sl "main/internal/logs"
	"main/internal/model"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Connection) Launch(cfg *config.Config, log *slog.Logger) error {

	// Забираем конфиг для дальнейшего использования в БД
	c.DB.Cfg = cfg
	c.DB.Logger = log

	conn, err := amqp.Dial("amqp://" + cfg.RABBIT_LOGIN + ":" + cfg.RABBIT_PASS + "@" + cfg.RABBIT_HOST + ":" + cfg.RABBIT_PORT + "/")
	if err != nil {
		return err
	}

	// Закрываем подключение в случае удачной попытки подключения
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	// Закрываем подключение в случае удачной попытки подключения
	defer func() {
		_ = ch.Close()
	}()

	q, err := ch.QueueDeclare(
		cfg.RABBIT_NAME, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			log.Info(string(message.Body))
			r, err := c.checkQeueu(message.Body)
			if err != nil {
				log.Error("error in Queue script.", sl.Error(err))
			}
			log.Info("call with.", slog.Any("VALUE", r))
			if err := call(r); err != nil {
				log.Error("Error from call.", sl.Error(err))
			}
		}
	}()

	// Тут логика в вечно открытом канале для читки, в будущем доработаем, чтобы закрывался
	// Если канал закроется, то горутина закроется и читать некому будет
	<-forever

	return nil
}

func (c Connection) checkQeueu(message []byte) (*model.Response, error) {

	json.Unmarshal(message, &c.Val)
	resp, err := c.DB.Launch(c.Val.City, c.Val.ID, c.Val.Phone)
	if err != nil {
		return nil, err
	}

	return &model.Response{
		Extension: resp.Extension,
		PhoneCity: resp.PhoneCity,
		CrmID:     c.Val.ID,
		PhoneLine: c.Val.Phone,
		UserCrm:   resp.UserCrm,
	}, nil
}

func call(data *model.Response) error {

	client := &http.Client{}
	url := "https://ykcomp.ru/mango/mango/queue/"
	param, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
