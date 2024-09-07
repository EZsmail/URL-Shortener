package main

import (
	"fmt"
	"restapi/URL-Shortener/internal/config"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: log

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server:
}
