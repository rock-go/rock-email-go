# rock-email-go

rock-go 系统的邮件发送组件。

# 使用说明

## 导入

```go
import email "github.com/rock-go/rock-email-go"
```

## 注册

```go
rock.Inject(xcall.Rock, email.LuaInjectApi)
```

## lua脚本调用

```lua
-- 邮件发送模块
local email = rock.email {
    name = "email",
    server = "mail.xxx.com", -- 邮件发送服务器
    port = "587", --端口
    from = "sender@xxx.com", -- 邮件发送者
    password = "password", -- 邮件发送者密码
    buffer = 100 -- 缓存的邮件数量，针对大量邮件发送
}

proc.start(email) -- 启动

--发送的时候，调用send(param1,param2,param3 ,param4 , param5, ....)
-- param1: 邮件接收者列表，多个邮箱以逗号分隔
-- param2: 邮件主题
-- param3: 邮件内容
-- param4: 邮件内容
-- param4: 邮件内容
email.send("961756805@qq.com", "测试邮件",
        "这是一封测试邮件，请忽略", " 这是第二段内容")

-- 发送html或附件
local data = "<body background='https://security.eastmoney.com/assets/img/fake_email.png?m=sunke@eastmoney.com'>尊敬的同事：<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您好！VPN密码已重置，密码为618@VPNpwd，请及时登录并修改为强密码。</body>"
local content = {
    to = "961756805@qq.com",
    subject = "测试邮件",
    typ = "html",
    content = data,
    attach = ""
}

email.send_obj(content)
```

## go 调用

```go
package main

import (
	"encoding/json"
	"github.com/rock-go/rock/lua"
)

/*
调用rock-email-go的组件须定义下面的Mail interface和Obj struct。
具体可参考rock-mailbait-go模块的调用方式。
*/

// Mail rock-email-go 模块实现了该接口。
type Mail interface {
	lua.LightUserDataIFace
	SendMail(interface{}) error
}

// Obj rock-email-go 发送时的数据结构，为了通用性，其它模块调用时，用json序列化为[]byte，rock-email-go会将[]byte反序列化。
type Obj struct {
	To          string   `json:"to"`
	Subject     string   `json:"subject"`
	Typ         string   `json:"typ"` // text, html
	Content     []byte   `json:"content"`
	Attachments []string `json:"attachments"` // 附件链接
}

// Test 要调用邮件发送模块的其它对象
type Test struct {
	mail Mail // 传入该模块的邮件对象
}

func (t *Test) Send() error {
	obj := Obj{
		To:          "xxxx@qq.com",
		Subject:     "test",
		Typ:         "html",
		Content:     []byte("test"),
		Attachments: "test.docx",
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = t.mail.SendMail(data)
	return err
}
```