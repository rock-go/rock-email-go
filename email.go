package email

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"mime"
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
	em := &Email{cfg: cfg}
	em.S = lua.INIT
	em.T = EMAIL
	return em
}

// SendMail 参数传递采用interface{}，以供其它模块golang直接调用
func (e *Email) SendMail(v interface{}) error {
	if e.S == lua.CLOSE {
		return errors.New("mail send service closed")
	}

	var o Obj
	switch v.(type) {
	case []byte:
		err := json.Unmarshal(v.([]byte), &o)
		if err != nil {
			return err
		}
	case Obj:
		o = v.(Obj)
	default:
		logger.Errorf("error type mail content")
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

func (e *Email) Start() error {
	e.Init()
	go e.loop()
	e.S = lua.RUNNING
	e.U = time.Now()
	logger.Infof("%s email start successfully", e.Name())
	return nil
}

func (e *Email) loop() {
	for {
		select {
		case data, ok := <-e.mailChan:
			if !ok {
				logger.Errorf("get email content error")
				continue
			}

			if err := e.formatEmail(data); err != nil {
				logger.Errorf("format email object error: %v", err)
				continue
			}

			if err := e.em.SendWithStartTLS(e.cfg.addr(), e.auth,
				&tls.Config{ServerName: e.cfg.server, InsecureSkipVerify: true}); err != nil {
				logger.Errorf("send email error: %v", err)
			}

			e.em.Text = nil
			e.em.HTML = nil
			e.em.Attachments = nil

		case <-e.ctx.Done():
			logger.Errorf("%s email sender exit", e.Name())
			e.S = lua.CLOSE
			return
		}
	}
}

// 格式化邮件
func (e *Email) formatEmail(obj Obj) error {
	e.em.To = formatAddr(obj.To)
	e.em.Text = obj.Content
	if obj.Typ == "html" {
		e.em.Text = nil
		e.em.HTML = obj.Content
	}
	e.em.Subject = obj.Subject

	for i, a := range obj.Attachments {
		_, err := e.em.AttachFile(a)
		if err != nil {
			return err
		}
		// 防止中文乱码
		qEncodedFilename := mime.QEncoding.Encode("UTF-8", e.em.Attachments[i].Filename)
		e.em.Attachments[i].Filename = qEncodedFilename
	}
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
