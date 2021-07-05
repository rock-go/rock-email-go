package email

import (
	"context"
	"github.com/jordan-wright/email"
	"github.com/rock-go/rock/lua"
	"net/smtp"
	"strings"
)

type Config struct {
	name   string
	server string
	port   string

	from     string
	password string

	buffer int // 缓存邮件数量
}

type Email struct {
	lua.Super

	C    Config
	em   *email.Email
	auth smtp.Auth

	mailChan chan Obj // 邮件内容队列

	ctx    context.Context
	cancel context.CancelFunc

	status lua.LightUserDataStatus
	uptime string
}

type Obj struct {
	to      string
	subject string
	content []byte
}

func formatAddr(s string) []string {
	s = strings.Trim(s, " ")
	return strings.Split(s, ",")
}
