package attendanceops

import (
	"context"
	"math/rand"
	"time"
)

type Msg struct {
	message string
}

func(m *Msg) AckMsg() error {
	return nil
}

func(m *Msg) NackMsg() error {
	return nil
}

func(m *Msg) KeepAlive() error {
	return nil
}


func fetchMessages(numOfMessages int, timeout time.Duration) ([]Msg, error) {
	return nil, nil
}

type MsgProcessor struct {
	Shutdown bool
}

func (w *MsgProcessor) Stop() {
	w.Shutdown = true
}

func (w *MsgProcessor) Run() error {
	for !w.Shutdown {
		msg, err	:= 
			fetchMessages(10, 10*time.Second)
		if err != nil {
			// log error
			continue
		}

		if err = w.ProcessMsgs(msg); err != nil {
			// log error
			continue
		}
	}

	return nil
}

func (w *MsgProcessor) ProcessMsgs(msg []Msg) error {
	for _, m := range msg {
		go w.ProcessMsg(m)
	}

	return nil
}

func (w *MsgProcessor) ProcessMsg(m Msg) error {
	// create cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// keep alive
	go w.KeepAlive(ctx, m)

	// process the message
	time.Sleep(1 * time.Second)

	// create random error
	if rand.Intn(2) == 0 {
		// if error occurs Nack the message
		return m.NackMsg()
	} else {
		// if no error Ack the message
		return m.AckMsg()
	}
}

func (w *MsgProcessor) KeepAlive(ctx context.Context, m Msg) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		default:
			if err := m.KeepAlive(); err != nil {
				// log error
				return err
			}
			time.Sleep(1 * time.Second)
		}
	}
}