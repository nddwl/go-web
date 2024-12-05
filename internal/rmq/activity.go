package rmq

import (
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"go-web/internal/model"
	"time"
)

type Activity struct {
	*Rmq
}

func NewActivity(rmq *Rmq) *Activity {
	activity := &Activity{rmq}
	return activity
}

func (t *Activity) PublishRecord(record model.ActivityRecord) (err error) {
	ch, err := t.get()
	if err != nil {
		return
	}
	defer t.put(ch)
	data, err := json.Marshal(&record)
	if err != nil {
		return
	}
	return ch.Publish("activity", "record", false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     0,
		Body:         data,
	})
}

func (t *Activity) PublishRecords(record []model.ActivityRecord) (err error) {
	ch, err := t.get()
	if err != nil {
		return
	}
	defer t.put(ch)
	data, err := json.Marshal(&record)
	if err != nil {
		return
	}
	return ch.Publish("activity", "record", false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Priority:     0,
		Body:         data,
	})
}

func (t *Activity) ConsumeRecord(fn func(record []model.ActivityRecord) error) (err error) {
	ch, err := t.get()
	if err != nil {
		return
	}
	defer t.put(ch)
	err = ch.Qos(50, 0, false)
	if err != nil {
		return
	}
	msgs, err := ch.Consume("activity_record", "", false, false, false, false, nil)
	if err != nil {
		return
	}
	i := 0
	s := make([]model.ActivityRecord, 50)
	msg := amqp091.Delivery{}
	for {
		select {
		case msg = <-msgs:
			err1 := json.Unmarshal(msg.Body, &s[i])
			if err1 != nil {
				msg.Nack(false, false)
			} else {
				i++
			}
		case <-time.After(time.Second * 3):
			if i > 0 {
				msg.Ack(true)
				err = fn(s[:i])
				if err != nil {
					return t.PublishRecords(s)
				}
				s = make([]model.ActivityRecord, 50)
			}
			return
		}
		if i == 50 {
			msg.Ack(true)
			err = fn(s)
			if err != nil {
				return t.PublishRecords(s)
			}
			s = make([]model.ActivityRecord, 50)
			i = 0
		}
	}
}
