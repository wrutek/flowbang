package screen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	term "github.com/nsf/termbox-go"
)

func RenderChooseList() {
	var msg *string
	list := []string{"wikt", "tor", "rut", "ka"}
	position := 0
	maxPosition := len(list)

	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

renderMenuLoop:
	for {
		reset()
		renderList(list, position)
		if msg != nil {
			fmt.Println(*msg)
		}
		position, msg, _ = fetchMenuPosition(position, maxPosition-1)
		if position == -1234 {
			break renderMenuLoop
		}

	}
}

func AskQuestion(question string) (response string) {
	response = ""
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	reset()
	fmt.Println(question)
	// write exactly this many of "-" as is letters in `question`
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(question)))

responseLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEnter:
				break responseLoop
			case term.KeySpace:
				fmt.Print(" ")
				response = response + " "
			default:
				fmt.Print(string(ev.Ch))
				response = response + string(ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
	return
}

func reset() {
	term.Sync() // cosmestic purpose
}

func fetchMenuPosition(actualPos, maxPos int) (pos int, msg *string, err error) {
	switch ev := term.PollEvent(); ev.Type {
	case term.EventKey:
		switch ev.Key {
		case term.KeyArrowUp:
			if actualPos > 0 {
				return actualPos - 1, nil, nil
			}
			return 0, nil, nil
		case term.KeyArrowDown:
			if actualPos < maxPos {
				return actualPos + 1, nil, nil
			}
			return maxPos, nil, nil
		case term.KeyEsc:
			return -1234, nil, nil
		default:
			msg := fmt.Sprintf("Try to use only arrrows: %c", ev.Ch)
			return actualPos, &msg, nil

		}
	case term.EventError:
		panic(ev.Err)
	}
	return 0, nil, nil
}

func renderList(list []string, position int) {
	var prefix string
	for i, menu := range list {
		prefix = "[ ]"
		if i == position {
			prefix = "[*]"
		}
		fmt.Println(prefix, menu)
	}
}
