package main

import (
	"equinox/internal/routers"
	"fmt"
)

func LaunchRouter(host string, port int) {
	r := routers.SetupRouter()
	addr := fmt.Sprintf("%s:%d", host, port)
	r.Run(addr)
}

func main() {
	LaunchRouter("localhost", 8080)
}
