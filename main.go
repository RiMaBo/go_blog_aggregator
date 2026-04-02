package main

import (
	"internal/config"

	"fmt"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = cfg.SetUser("rimabo")
	if err != nil {
		fmt.Println(err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)
}
