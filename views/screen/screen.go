package screen

import (
	"fmt"
	"regexp"
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
			renderSeparator(message)
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
	renderSeparator(question)

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

// ProgressBar render progress bar for given atributes
// You need to close terminal manually here and open it.
// `bar.Init(...)`
// `bar.Update(...)`
// `bar.Close()`
type ProgressBar struct {
	Title  string
	Desc   string
	index  int
	Length int
	SWidth int
}

// Init initialize progress bar
func (bar *ProgressBar) Init() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	bar.SWidth, _ = term.Size()
	bar.SWidth -= 2

}

// Close progress bar
func (bar *ProgressBar) Close() {
	term.Close()
}

// Update status of progress bar
func (bar *ProgressBar) Update(index int) {
	bar.SWidth, _ = term.Size()
	bar.SWidth -= 2
	reset()
	fmt.Println(bar.Title)
	renderSeparator(bar.Title)
	screenIndex := float32(index) * (float32(bar.SWidth) / float32(bar.Length))
	fmt.Printf("[%s%s]\n", strings.Repeat("#", int(screenIndex)), strings.Repeat(".", bar.SWidth-int(screenIndex)))
	fmt.Println(bar.Desc)
}

func reset() {
	term.Sync() // cosmestic purpose
}

func renderSeparator(msg string) {
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(msg)))
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
			prefix = "[\x1b[31m*\x1b[0m]"
		}
		fmt.Println(prefix, stringElipsis(item.GetFullName(), 4))
	}
}

// helper functions

func stringElipsis(in string, prefixLen int) (out string) {
	additionalPrefix := 0

	// We need to ignore vt100 escape characters from sting
	re, _ := regexp.Compile(`\033\[[0-9]+m`)
	matches := re.FindAllString(in, -1)
	if matches != nil {
		for _, match := range matches {
			additionalPrefix += utf8.RuneCountInString(match)
		}
	}

	out = in
	length, _ := term.Size()
	length -= prefixLen - additionalPrefix
	runes := []rune(out)
	if len(runes) > length {
		out = string(runes[:length-3]) + "..."
	}
	return
}
