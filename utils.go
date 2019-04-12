package main

import "fmt"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func log(args... interface{}) {
	fmt.Println(args)
}
