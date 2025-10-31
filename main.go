package main

import (
	"fmt"

	"github.com/jkk290/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		return
	}
	if err := cfg.SetUser("jared"); err != nil {
		fmt.Printf("error setting current user: %v", err)
		return
	}
	updatedCfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading updated config: %v\n", err)
		return
	}
	fmt.Println(updatedCfg.DBURL)
	fmt.Println(updatedCfg.CurrentUserName)
}
