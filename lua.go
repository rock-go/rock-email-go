package email

import (
	"bytes"
	"fmt"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"reflect"
	"strings"
)

var (
	EMAIL = reflect.TypeOf((*Email)(nil)).String()
)

func (e *Email) Index(L *lua.LState, key string) lua.LValue {
	if key == "start" {
		return lua.NewFunction(e.start)
	}
	if key == "close" {
		return lua.NewFunction(e.close)
	}
	if key == "send" {
		return lua.NewFunction(e.LSend)
	}
	if key == "send_obj" {
		return lua.NewFunction(e.LSendObj)
	}

	return lua.LNil
}

func (e *Email) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		e.cfg.name = lua.CheckString(L, val)
	case "server":
		e.cfg.server = lua.CheckString(L, val)
	case "port":
		e.cfg.port = lua.CheckInt(L, val)
	case "from":
		e.cfg.from = lua.CheckString(L, val)
	case "password":
		e.cfg.password = lua.CheckString(L, val)
	case "buffer":
		e.cfg.buffer = lua.CheckInt(L, val)
	}
}

func (e *Email) start(L *lua.LState) int {
	if e.State() == lua.RUNNING {
		L.RaiseError("%s email is already running", e.cfg.name)
		return 0
	}

	err := e.Start()
	if err != nil {
		L.RaiseError("email sender start error: %v", err)
	}

	return 0
}

func (e *Email) close(L *lua.LState) int {
	if e.S == lua.CLOSE {
		L.RaiseError("%s email is already closed", e.cfg.name)
		return 0
	}

	err := e.Close()
	if err != nil {
		L.RaiseError("email sender close error: %v", err)
	}
	return 0
}

func (e *Email) LSend(L *lua.LState) int {
	n := L.GetTop()
	if n < 3 {
		L.TypeError(3, lua.LTString)
		return 0
	}

	to := L.CheckString(1)
	subject := L.CheckString(2)

	var buff bytes.Buffer
	for i := 3; i <= n; i++ {
		item := lua.S2B(fmt.Sprintf("%v", L.Get(i)))
		buff.Write(item)
	}

	obj := Obj{
		To:      to,
		Subject: subject,
		Content: buff.Bytes(),
	}

	if err := e.SendMail(obj); err != nil {
		L.RaiseError("send email error: %v", err)
	}

	return 0
}

// LSendObj 发送更为详细的邮件.参数：收件人，主题，邮件格式，邮件正文，附件
func (e *Email) LSendObj(L *lua.LState) int {
	obj := checkObj(L)
	if err := e.SendMail(*obj); err != nil {
		L.RaiseError("send email error: %v", err)
	}

	return 0
}

func checkObj(L *lua.LState) *Obj {
	tb := L.CheckTable(1)
	to := tb.RawGetString("to").String()
	subject := tb.RawGetString("subject").String()
	typ := tb.RawGetString("typ").String()
	content := lua.S2B(tb.RawGetString("content").String())
	attach := tb.RawGetString("attach").String()
	attachments := make([]string, 0)
	if attach != "" {
		attachments = strings.Split(attach, ",")
	}

	return &Obj{
		To:          to,
		Subject:     subject,
		Typ:         typ,
		Content:     content,
		Attachments: attachments,
	}
}

func newLuaEmail(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, EMAIL)
	if proc.IsNil() {
		proc.Set(newEmail(cfg))
		goto done
	}
	proc.Value.(*Email).cfg = cfg

done:
	L.Push(proc)
	return 1
}

func LuaInjectApi(env xcall.Env) {
	env.Set("email", lua.NewFunction(newLuaEmail))
}
