package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/wrutek/flowbang/config"
	"github.com/wrutek/flowbang/controllers"
)

// ConfigDir a path to place where all config file will be stored
var ConfigDir string = filepath.Join(".config", "flowbang")

// ConfigPath a path to main config file
var ConfigPath string = filepath.Join(ConfigDir, "flowbang.conf")

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
		err := config.Configure()
		if err != nil {
			fmt.Println("Error: Flowbang configuration failed")
			panic(err)
		}
	case "get_todos":
		getTodos()
	case "work_on":
		// screen.RenderChooseList()
		controllers.ChooseIssue()

		fmt.Println("Your are trying to work on not exisitng issue")
	default:
		fmt.Println("Error: wrong command: ", command)
	}
}

func getTodos() {
	fmt.Println("Getting todos")
}
