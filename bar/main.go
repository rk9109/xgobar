package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb"
)

func main() {
	X, err := xgb.NewConn()
	if err != nil {
		fmt.Println(err)
		return
	}

	bar := New(X)
	err = bar.Map()
	if err != nil {
		fmt.Println(err)
		return
	}

	// spin
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
