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

	err = api.Get(fmt.Sprintf("projects/columns/%d/cards", cfg.TodoColumnID), &cards, headers)
	if err != nil {
		return
	}
	var issue issueItem
	for _, card := range cards {
		// Fetch related issue/pull request to this card
		err = api.Get(card.ContentURL, &issue, headers)
		if err != nil {
			return
		}
		// Get all users that are assigned to this ticket
		assignees := ""
		for _, assignee := range issue.Assignees {
			assignees = assignees + fmt.Sprintf("[%s]", assignee.Login)
		}

		card.Name = fmt.Sprintf("%s %s", issue.Name, assignees)
		items = append(items, screen.Item(card))
	}
	screen.ChooseList(items, "Select a ticket on which you want to work on")
	fmt.Println(items)

	return
}
