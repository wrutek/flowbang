package screen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	term "github.com/nsf/termbox-go"
)

type ScreenItem interface {
	GetId() int
	GetName() string
	GetFullName() string
}

func ChooseList(items []ScreenItem, message string) ScreenItem {
	var msg *string
	position := 0
	prevPosition := 0
	maxPosition := len(items)
	isBreak := false

	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

renderMenuLoop:
	for {
		reset()

		if len(message) != 0 {
			fmt.Println(message)
		}
		renderList(items, position)
		if msg != nil {
			fmt.Println(*msg)
		}
		position, msg, _, isBreak = fetchMenuPosition(position, maxPosition-1)
		if isBreak == true {
			// Close list menu
			break renderMenuLoop
		}
		prevPosition = position

	}
	reset()
	return items[prevPosition]
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
				// Quit question screen on enter
				break responseLoop
			case term.KeySpace:
				// space is not uncer ev.CH event
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

func fetchMenuPosition(actualPos, maxPos int) (pos int, msg *string, err error, is_breake bool) {
	is_breake = false
	switch ev := term.PollEvent(); ev.Type {
	case term.EventKey:
		switch ev.Key {
		case term.KeyArrowUp:
			if actualPos > 0 {
				return actualPos - 1, nil, nil, false
			}
			return maxPos, nil, nil, false
		case term.KeyArrowDown:
			if actualPos < maxPos {
				return actualPos + 1, nil, nil, false
			}
			return 0, nil, nil, false
		case term.KeyEsc:
			// Close list screen
			return actualPos, nil, nil, true
		case term.KeyEnter:
			return actualPos, nil, nil, true
		default:
			msg := fmt.Sprintf("Try to use only arrrows: %c", ev.Ch)
			return actualPos, &msg, nil, is_breake

		}
	case term.EventError:
		panic(ev.Err)
	}
	return 0, nil, nil, is_breake
}

func renderList(list []ScreenItem, position int) {
	var prefix string
	for i, item := range list {
		prefix = "[ ]"
		if i == position {
			prefix = "[*]"
		}
		fmt.Println(prefix, item.GetFullName())
	}
}
