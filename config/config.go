// Package config is responsible for creating and or just opening config file
// and save user configuration of this tool to this file in yaml format
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wrutek/flowbang/config/configmodels"
	"github.com/wrutek/flowbang/settings"
	"github.com/wrutek/flowbang/views/api"
	"github.com/wrutek/flowbang/views/screen"
	"gopkg.in/yaml.v2"
)

// Configure one of the main functions of flobang. Configure a whole system
// get from the user: oauth token, projects and repositories to github
func Configure() (err error) {
processThisCommand:
	for {
		answare := screen.AskQuestion("This will override your current configuration. Are you sure? [y/n]")
		switch answare {
		case "Y", "y":
			break processThisCommand
		case "n", "N":
			return fmt.Errorf("you choose not to continue configuration")
		}
	}
	file, err := createConfigFile(settings.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	token := screen.AskQuestion("Please provide us an oauth token to you github profile:")

	headers := map[string]string{
		"Authorization": "token " + token,
	}
	var repos []api.RepoItem
	err = api.RawRequest("GET", "user/repos", headers, nil, &repos)
	if err != nil {
		panic(err)
	}

	var items []screen.Item
	for _, repo := range repos {
		items = append(items, screen.Item(repo))
	}
	workingRepo := screen.ChooseList(items, "Select working repository")
	issueRepo := screen.ChooseList(items, "Select board repository")

	headers["Accept"] = "application/vnd.github.inertia-preview+json"
	var projects []api.ProjectItem
	err = api.RawRequest("GET", "repos/"+issueRepo.GetFullName()+"/projects", headers, nil, &projects)
	if err != nil {
		panic(err)
	}
	items = nil
	for _, project := range projects {
		items = append(items, screen.Item(project))
	}
	project := screen.ChooseList(items, "Select project you will be working on")

	var columns []api.ColumnItem
	err = api.RawRequest("GET", fmt.Sprintf("projects/%d/columns", project.GetID()), headers, nil, &columns)
	if err != nil {
		panic(err)
	}
	items = nil
	for _, column := range columns {
		items = append(items, screen.Item(column))
	}
	todoColumn := screen.ChooseList(items, "Select a \"To Do\" column")
	inProgressColumn := screen.ChooseList(items, "Select an \"In Progress\" column")
	doneColumn := screen.ChooseList(items, "Select a \"Done\" column")

	cfg := configmodels.Configuration{
		IssueRepoID:        issueRepo.GetID(),
		WorkingRepoID:      workingRepo.GetID(),
		OauthToken:         token,
		ProjectID:          project.GetID(),
		TodoColumnID:       todoColumn.GetID(),
		InprogressColumnID: inProgressColumn.GetID(),
		DoneColumnID:       doneColumn.GetID(),
	}
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return
}

func createConfigFile(configFileName string) (*os.File, error) {
	// User is required to fetch home directory
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return nil, err
	}

	configPath := filepath.Join(currDir, configFileName)

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return os.Create(configPath)
	} else if err != nil {
		return nil, err
	}
	return os.OpenFile(configPath, os.O_WRONLY, 600)
}
