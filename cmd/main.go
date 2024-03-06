package main

import (
	"equinox/internal/routers"
)

func main() {
	r := routers.SetupRouter()
	r.Run("localhost:8080")
}
