package screen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	term "github.com/nsf/termbox-go"
)

// Item is an interface for every input that is viewable
type Item interface {
	GetID() int
	GetName() string
	GetFullName() string
}

// ChooseList shows screen with a
// list of items that can be choosen
func ChooseList(items []Item, message string) Item {
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
		position, msg, isBreak, _ = fetchMenuPosition(position, maxPosition-1)
		if isBreak == true {
			// Close list menu
			break renderMenuLoop
		}
		prevPosition = position

	}
	reset()
	return items[prevPosition]
}

// AskQuestion simple question->answare screen
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

func fetchMenuPosition(actualPos, maxPos int) (pos int, msg *string, isBreak bool, err error) {
	isBreak = false
	switch ev := term.PollEvent(); ev.Type {
	case term.EventKey:
		switch ev.Key {
		case term.KeyArrowUp:
			if actualPos > 0 {
				return actualPos - 1, nil, false, nil
			}
			return maxPos, nil, false, nil
		case term.KeyArrowDown:
			if actualPos < maxPos {
				return actualPos + 1, nil, false, nil
			}
			return 0, nil, false, nil
		case term.KeyEsc:
			// Close list screen
			return actualPos, nil, true, nil
		case term.KeyEnter:
			return actualPos, nil, true, nil
		default:
			msg := fmt.Sprintf("Try to use only arrrows: %c", ev.Ch)
			return actualPos, &msg, isBreak, nil

		}
	case term.EventError:
		panic(ev.Err)
	}
	return 0, nil, isBreak, nil
}

func renderList(list []Item, position int) {
	var prefix string
	for i, item := range list {
		prefix = "[ ]"
		if i == position {
			prefix = "[*]"
		}
		fmt.Println(prefix, item.GetFullName())
	}
}
