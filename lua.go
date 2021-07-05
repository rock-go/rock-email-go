package email

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
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

	return lua.LNil
}

func (e *Email) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		e.C.name = lua.CheckString(L, val)
	case "server":
		e.C.server = lua.CheckString(L, val)
	case "port":
		e.C.port = lua.CheckString(L, val)
	case "from":
		e.C.from = lua.CheckString(L, val)
	case "password":
		e.C.password = lua.CheckString(L, val)
	case "buffer":
		e.C.buffer = lua.CheckInt(L, val)
	}
}

func (e *Email) start(L *lua.LState) int {
	if e.status == lua.RUNNING {
		L.RaiseError("%s email is already running", e.C.name)
		return 0
	}

	err := e.Start()
	if err != nil {
		L.RaiseError("email sender start error: %v", err)
	}

	return 0
}

func (e *Email) close(L *lua.LState) int {
	if e.status == lua.CLOSE {
		L.RaiseError("%s email is already closed", e.C.name)
		return 0
	}

	err := e.Close()
	if err != nil {
		L.RaiseError("email sender close error: %v", err)
	}

	return 0
}

func (e *Email) LSend(L *lua.LState) int {
	to := L.CheckString(1)
	subject := L.CheckString(2)
	content := []byte(L.CheckString(3))
	obj := Obj{
		to:      to,
		subject: subject,
		content: content,
	}

	if err := e.SendMail(obj); err != nil {
		L.RaiseError("send email error: %v", err)
	}

	return 0
}

func createEmailUserData(L *lua.LState) int {
	opt := L.CheckTable(1)
	cfg := Config{
		name:     opt.CheckString("name", "email"),
		server:   opt.CheckString("server", "mail.eastmoney.com"),
		port:     opt.CheckString("port", "25"),
		from:     opt.CheckString("from", "am35@eastmoney.com"),
		password: opt.CheckString("password", "53xcxWeiXin*0.aq"),
		buffer:   opt.CheckInt("buffer", 10),
	}

	email := &Email{C: cfg}

	var obj *Email
	var ok bool

	proc := L.NewProc(email.C.name)
	if proc.Value == nil {
		proc.Value = email
		goto done
	}

	obj, ok = proc.Value.(*Email)
	if !ok {
		L.RaiseError("invalid email proc")
		return 0
	}
	obj.C = cfg

done:
	L.Push(proc)
	return 1
}

func LuaInjectApi(env xcall.Env) {
	env.Set("email", lua.NewFunction(createEmailUserData))
}
