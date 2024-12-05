package rmq

import (
	"errors"
	"github.com/rabbitmq/amqp091-go"
	"go-web/utils/config"
	"sync/atomic"
	"time"
)

type Option struct {
	Url         string
	MaxCap      int
	MaxIdle     int
	IdleTimeout time.Duration
	WaitTimeout time.Duration
}

type Group struct {
	Activity *Activity
	User     *User
}

type idleConn struct {
	channel *amqp091.Channel
	t       time.Time
}

type Rmq struct {
	conn         *amqp091.Connection
	idles        chan *idleConn
	openingIdles int32
	opt          *Option
	Group
}

func New() *Rmq {
	opt := &Option{
		Url:         config.Amqp.Url,
		MaxCap:      50,
		MaxIdle:     10,
		IdleTimeout: 15 * time.Minute,
		WaitTimeout: 5 * time.Minute,
	}
	rmq := &Rmq{
		conn:  nil,
		idles: make(chan *idleConn, opt.MaxIdle),
		opt:   opt,
	}
	conn, err := amqp091.Dial(opt.Url)
	if err != nil {
		panic(err)
	}
	rmq.conn = conn
	for i := 0; i < opt.MaxIdle; i++ {
		channel, err1 := rmq.conn.Channel()
		if err1 != nil {
			panic("CreatePoolFailed")
		}
		rmq.idles <- &idleConn{
			channel: channel,
			t:       time.Now(),
		}
	}
	atomic.AddInt32(&rmq.openingIdles, int32(rmq.opt.MaxIdle))
	rmq.initGroup()
	return rmq
}

func (t *Rmq) initGroup() {
	t.Group = Group{
		Activity: NewActivity(t),
		User:     NewUser(t),
	}
}

func (t *Rmq) get() (channel *amqp091.Channel, err error) {
	for {
		select {
		case idle := <-t.idles:
			if !idle.t.Add(t.opt.IdleTimeout).After(time.Now()) || idle.channel.IsClosed() {
				atomic.AddInt32(&t.openingIdles, -1)
				idle.channel.Close()
				continue
			}
			return idle.channel, nil
		case <-time.After(t.opt.WaitTimeout):
			err = errors.New("WaitTimeout")
			return
		default:
			if atomic.LoadInt32(&t.openingIdles) >= int32(t.opt.MaxCap) {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			channel, err = t.conn.Channel()
			if err == nil {
				atomic.AddInt32(&t.openingIdles, 1)
			}
			return
		}
	}
}

func (t *Rmq) put(channel *amqp091.Channel) {
	select {
	case t.idles <- &idleConn{
		channel: channel,
		t:       time.Now(),
	}:
		return
	default:
		atomic.AddInt32(&t.openingIdles, -1)
		channel.Close()
		return
	}
}
