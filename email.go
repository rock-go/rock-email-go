package email

import (
	"context"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"net/smtp"
	"time"
)

func (e *Email) SendMail(o Obj) error {
	if e.status == lua.CLOSE {
		return errors.New("mail send service closed")
	}

	e.mailChan <- o
	return nil
}

func (e *Email) Init() {
	e.em = email.NewEmail()
	e.em.From = e.C.from
	e.auth = smtp.PlainAuth("", e.C.from, e.C.password, e.C.server)

	e.mailChan = make(chan Obj, e.C.buffer)
	e.ctx, e.cancel = context.WithCancel(context.Background())
}

func (e *Email) Start() error {
	e.Init()
	addr := fmt.Sprintf("%s:%s", e.C.server, e.C.port)
	e.status = lua.RUNNING
	e.uptime = time.Now().Format("2006-01-02 15:04:05")

	go func() {
		for {
			select {
			case data, ok := <-e.mailChan:
				if !ok {
					logger.Errorf("get email content error")
					continue
				}

				e.em.To = formatAddr(data.to)
				e.em.Text = data.content
				e.em.Subject = data.subject
				if err := e.em.Send(addr, e.auth); err != nil {
					logger.Errorf("send email error: %v", err)
				}
			case <-e.ctx.Done():
				logger.Errorf("%s email sender exit", e.C.name)
				e.status = lua.CLOSE
				return
			}
		}
	}()

	logger.Infof("%s email start successfully", e.C.name)
	return nil
}

func (e *Email) Close() error {
	e.status = lua.CLOSE
	if e.cancel != nil {
		e.cancel()
	}

	close(e.mailChan)
	return nil
}

func (e *Email) Name() string {
	return e.C.name
}

func (e *Email) State() lua.LightUserDataStatus {
	return e.status
}

func (e *Email) Type() string {
	return "email sender"
}

func (e *Email) Status() string {
	return fmt.Sprintf("name: %s, status: %s, uptime: %s",
		e.C.name, e.status, e.uptime)
}
