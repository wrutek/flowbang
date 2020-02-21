package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/wrutek/flowbang/config"
	"github.com/wrutek/flowbang/views/screen"
)

var CONFIG_DIR string = filepath.Join(".config", "flowbang")
var CONFIG_PATH string = filepath.Join(CONFIG_DIR, "flowbang.conf")

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Error: Please provide only one command to this program!")
		return
	}
	command := args[0]
	fmt.Println("Yours command is: ", command)

	switch command {
	case "configure":
		/* configure flowbang and prepare for first run */
		_, err := config.Configure(CONFIG_DIR, CONFIG_PATH)
		if err != nil {
			fmt.Println("Error: Flowbang configuration failed")
			panic(err)
		}
	case "get_todos":
		get_todos()
	case "work_on":
		screen.RenderChooseList()
		fmt.Println("Your are trying to work on not exisitng issue")
	default:
		fmt.Println("Error: wrong command: ", command)
	}
}

func get_todos() {
	fmt.Println("Getting todos")
}
