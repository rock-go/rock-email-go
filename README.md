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
    "这是一封测试邮件，请忽略" , " 这是第二段内容")
```