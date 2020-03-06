package controllers

import (
	"fmt"

	"github.com/wrutek/flowbang/config/configgetter"
	"github.com/wrutek/flowbang/views/api"
	"github.com/wrutek/flowbang/views/screen"
)

// CardItem represents single card on github projects board
type CardItem struct {
	ID         int    `json:"id"`
	ContentURL string `json:"content_url"`
	Name       string
}

// GetID get card ID
func (card CardItem) GetID() int {
	return card.ID
}

// GetName get card name
func (card CardItem) GetName() string {
	return card.Name
}

// GetFullName get card full_name
func (card CardItem) GetFullName() string {
	return card.Name
}

type issueItem struct {
	Name      string `json:"title"`
	Number    int    `json:"number"`
	Assignees []struct {
		Login string `json:"login"`
	} `json:"assignees"`
}

// ChooseIssue fetch all issues from project
// column where issues are ready to be `in progress`
func ChooseIssue() (err error) {

	cfg, err := configgetter.GetConfiguration()

	headers := map[string]string{
		"Accept": "application/vnd.github.inertia-preview+json",
	}
	var cards []CardItem
	var items []screen.Item

	fmt.Println("Fetching cards from the board...")
	err = api.Get(fmt.Sprintf("projects/columns/%d/cards", cfg.TodoColumnID), &cards, headers)
	if err != nil {
		return
	}
	fmt.Println("Cards fetched.")
	fmt.Print("fetching issues for all cards...\n")
	var issue issueItem
	progressBar := screen.ProgressBar{
		Title:  "Fetching issues for board cards",
		Desc:   "",
		Length: len(cards),
	}
	progressBar.Init()
	for i, card := range cards {
		// Fetch related issue/pull request to this card
		progressBar.Update(i + 1)
		err = api.Get(card.ContentURL, &issue, headers)
		if err != nil {
			return
		}
		// Get all users that are assigned to this ticket
		assignees := ""
		for _, assignee := range issue.Assignees {
			assignees = assignees + fmt.Sprintf("<\x1b[33m%s\x1b[0m> ", assignee.Login)
		}
		fmt.Print("\n")
		issueNumber := fmt.Sprintf("\x1b[34m#%d\x1b[0m", issue.Number)

		card.Name = fmt.Sprintf("%s %s%s", issueNumber, assignees, issue.Name)
		items = append(items, screen.Item(card))
	}
	progressBar.Close()
	item := screen.ChooseList(items, "Select a ticket on which you want to work on")
	fmt.Println(item)

	return
}
