package main

import (
	"fmt"
	"github.com/JasonLeemz/obs-tools/tools"
)

func main() {
	err := tools.Push()
	fmt.Println(err)

}
