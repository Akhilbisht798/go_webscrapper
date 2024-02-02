package main

import (
	"fmt"
	"os"
)

func putValueInFile(file string) {
	file += ".csv"
	fo, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = fo.WriteString("Product,Rating,Price\n")
	if err != nil {
		panic(err)
	}
	for _, val := range result {
		_, err := fo.WriteString(val)
		if err != nil {
			continue
		}
	}
}
