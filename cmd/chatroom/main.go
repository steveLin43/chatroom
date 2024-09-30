package main

import (
	"chatroom/global"
	"chatroom/server"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s

`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner, addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
