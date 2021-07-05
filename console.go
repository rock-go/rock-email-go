package email

import "github.com/rock-go/rock/lua"

func (e *Email) Header(out lua.Printer) {
	out.Printf("type: %s", e.Type())
	out.Printf("uptime: %s", e.uptime)
	out.Printf("version: v1.0.0")
	out.Println("")
}

func (e *Email) Show(out lua.Printer) {
	e.Header(out)

	out.Printf("name: %s", e.C.name)
	out.Printf("server: %s", e.C.server)
	out.Printf("port: %s", e.C.port)
	out.Printf("from: %s", e.C.from)
	out.Printf("password: ******")
	out.Printf("buffer: %d", e.C.buffer)
	out.Println("")
}

func (e *Email) Help(out lua.Printer) {
	e.Header(out)

	out.Printf(".start() 启动")
	out.Printf(".close() 关闭")
	out.Println("")
}
