package email

import (
	"github.com/rock-go/rock/lua"
	"strings"
	"github.com/rock-go/rock/utils"
	"errors"
	"fmt"
)

type config struct {
	name      string
	server    string
	port      int

	from      string
	password  string

	buffer int // 缓存邮件数量
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}

	tab.ForEach(func(key lua.LValue, val lua.LValue) {
		if key.Type() != lua.LTString {
			L.RaiseError("invalid options %s" , key.Type().String())
			return
		}

		switch key.String() {
		case "name": cfg.name = val.String()
		case "from": cfg.from = val.String()
		case "server": cfg.server = val.String()
		case "password": cfg.password = val.String()
		case "port": cfg.port = utils.LValueToInt(val , 0)
		case "buffer": cfg.buffer = utils.LValueToInt(val , 10)
		default:
			L.RaiseError("invalid options %s key" , key.String())
		}
	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v" , e)
		return nil
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := utils.Name(cfg.name) ; e != nil {
		return e
	}

	if cfg.port <= 0 || cfg.port >= 65535 {
		return errors.New("invalid port")
	}

	return nil
}

func (cfg *config) addr() string {
	return fmt.Sprintf("%s:%d" , cfg.server , cfg.port)
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
