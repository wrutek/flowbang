package screen

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	term "github.com/nsf/termbox-go"
)

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
	if bar.SWidth == int(screenIndex) {
		fmt.Printf("[\x1b[32m%s\x1b[30m%s\x1b[0m]\n", strings.Repeat("#", int(screenIndex)), strings.Repeat(".", bar.SWidth-int(screenIndex)))
		time.Sleep(time.Second * 2)
		return
	}
	fmt.Printf("[\x1b[33m%s\x1b[30m%s\x1b[0m]\n", strings.Repeat("#", int(screenIndex)), strings.Repeat(".", bar.SWidth-int(screenIndex)))
	fmt.Println(bar.Desc)
}

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
	var header string
	position := 0
	isBreak := false

	page := 0
	screenSize := 20 // default size of one paginated screen
	paginated := createItemsPages(items, screenSize)
	numOfPages := len(paginated)

	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

renderMenuLoop:
	for {
		reset()
		header = fmt.Sprintf("<%d/%d> %s", page+1, numOfPages, message)

		if len(message) != 0 {
			fmt.Println(header)
			renderSeparator(header)
		}
		renderList(paginated[page], position)
		if msg != nil {
			fmt.Println(*msg)
		}
		position, page, msg, isBreak, _ = fetchMenuPosition(position, page, len(paginated[page])-1)
		if isBreak == true {
			// Close list menu
			break renderMenuLoop
		}
		if numOfPages <= page {
			page = numOfPages - 1
		} else if page < 0 {
			page = 0
		}
	}
	return paginated[page][position]
}

// AskQuestion simple question->answare screen
func AskQuestion(question string) (response string) {
	response = ""
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

responseLoop:
	for {
		reset()
		fmt.Println(question)
		renderSeparator(question)
		fmt.Println(response)

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
			case term.KeyBackspace:
				if len(response) != 0 {
					response = response[:len(response)-1]
				}
			case term.KeyBackspace2:
				if len(response) != 0 {
					response = response[:len(response)-1]
				}
			default:
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

func createItemsPages(items []Item, pageSize int) (pages [][]Item) {
	itemsLen := len(items)
	for i := 0; i < itemsLen; i += pageSize {
		screenEnd := i + pageSize
		if screenEnd > itemsLen {
			screenEnd = itemsLen
		}
		pages = append(pages, items[i:screenEnd])
	}
	return
}

func renderSeparator(msg string) {
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(msg)))
}

func fetchMenuPosition(actualPos, page, maxPos int) (pos, pgIndex int, msg *string, isBreak bool, err error) {
	isBreak = false
	var message string
	switch ev := term.PollEvent(); ev.Type {
	case term.EventKey:
		switch ev.Key {
		case term.KeyArrowUp:
			if actualPos > 0 {
				return actualPos - 1, page, nil, false, nil
			}
			return maxPos, page, nil, false, nil
		case term.KeyArrowDown:
			if actualPos < maxPos {
				return actualPos + 1, page, nil, false, nil
			}
			return 0, page, nil, false, nil
		case term.KeyArrowRight:
			return actualPos, page + 1, &message, false, nil
		case term.KeyArrowLeft:
			return actualPos, page - 1, &message, false, nil
		case term.KeyEsc:
			// Close list screen
			return actualPos, page, nil, true, nil
		case term.KeyEnter:
			return actualPos, page, nil, true, nil
		default:
			message = fmt.Sprintf("Try to use only arrrows: %c", ev.Ch)
			return actualPos, page, &message, isBreak, nil

		}
	case term.EventError:
		panic(ev.Err)
	}
	return 0, 0, nil, isBreak, nil
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
