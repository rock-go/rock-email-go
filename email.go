package email

import (
	"context"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"net/smtp"
	"time"
)

type Email struct {
	lua.Super

	cfg  *config
	em   *email.Email
	auth smtp.Auth

	mailChan chan Obj // 邮件内容队列

	ctx    context.Context
	cancel context.CancelFunc
}

func newEmail(cfg *config) *Email {
	em := &Email{ cfg:cfg }
	em.S = lua.INIT
	em.T = EMAIL
	return em
}

func (e *Email) SendMail(o Obj) error {
	if e.S == lua.CLOSE {
		return errors.New("mail send service closed")
	}
	e.mailChan <- o
	return nil
}

func (e *Email) Init() {
	e.em = email.NewEmail()
	e.em.From = e.cfg.from
	e.auth = smtp.PlainAuth("", e.cfg.from, e.cfg.password, e.cfg.server)
	e.mailChan = make(chan Obj, e.cfg.buffer)
	e.ctx, e.cancel = context.WithCancel(context.Background())
}

func (e *Email) loop() {
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
			if err := e.em.Send(e.cfg.addr(), e.auth); err != nil {
				logger.Errorf("send email error: %v", err)
			}

		case <-e.ctx.Done():
			logger.Errorf("%s email sender exit", e.Name())
			e.S = lua.CLOSE
			return
		}
	}
}


func (e *Email) Start() error {
	e.Init()
	go e.loop()
	e.S = lua.RUNNING
	e.U = time.Now()
	logger.Infof("%s email start successfully", e.Name())
	return nil
}

func (e *Email) Close() error {
	e.S = lua.CLOSE
	if e.cancel != nil {
		e.cancel()
	}

	close(e.mailChan)
	return nil
}

func (e *Email) Name() string {
	return e.cfg.name
}