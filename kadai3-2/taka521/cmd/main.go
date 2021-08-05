package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/taka521/gopherdojo-studyroom/kadai3-2/taka521/pdown"
)

func main() {
	flag.Parse()
	dir := flag.Arg(0)
	url := flag.Arg(1)

	d := pdown.New()
	if err := d.Run(pdown.Input{URL: url, DownDir: dir}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("✨ Downlod successflly !!")
}
