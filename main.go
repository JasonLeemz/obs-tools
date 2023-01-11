package main

import (
	"fmt"
	"obs-tools/tools"
)

func main() {
	//err := tools.BatchConvert(os.Args[1], os.Args[2])
	//if err != nil {
	//	fmt.Println(err)
	//}

	err := tools.Push()
	fmt.Println(err)
}
