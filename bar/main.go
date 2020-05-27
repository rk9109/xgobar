package main

import "fmt"

func main() {
	//
	bar, err := NewBar()
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	err = bar.Map()
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	err = bar.Draw()
	if err != nil {
		fmt.Println(err)
		return
	}
}
