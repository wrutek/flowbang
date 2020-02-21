// package config is responsible for creating and or just opening config file
// and save user configuration of this tool to this file in yaml format
package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/wrutek/flowbang/views/api"
	"github.com/wrutek/flowbang/views/screen"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	IssueRepoID   int    `yaml:"issue_repo_id"`
	WorkingRepoID int    `yaml:"working_repo_id"`
	OauthToken    string `yaml:"oauth_token"`
}

func Configure(dirPath string, filePath string) (file *os.File, err error) {
	// get from the user: username, password, projects and repositories to github
	file, err = createConfigFile(dirPath, filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO: those values should be taken from user input.
	// 		 For now it's just dummy values
	token := screen.AskQuestion("Please provide us an oauth token to you github profile:")
	headers := map[string]string{
		"Authorization": "token " + token,
	}
	api.RawRequest("GET", "user/repos", headers, nil)
	// req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// client := &http.Client{}
	// req.Header.Add("Authorization", "token "+token)
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.CloseIdleConnections()

	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	// respBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(respBody)

	cfg := Configuration{
		IssueRepoID:   1234,
		WorkingRepoID: 4321,
		OauthToken:    token,
	}
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(cfg)
	if err != nil {
		return nil, err
	}

	return
}

// func GetConfiguration() conf *Configuration {

// }

func createConfigFile(dirPath string, filePath string) (*os.File, error) {
	// User is required to fetch home directory
	cUser, err := user.Current()
	if err != nil {
		fmt.Println("ERROR: Couldn't find currently logged in user.")
		return nil, err
	}

	config_dir := filepath.Join(cUser.HomeDir, dirPath)
	config_path := filepath.Join(cUser.HomeDir, filePath)

	// Create config directory if does not exists
	_, err = os.Stat(config_dir)
	if os.IsNotExist(err) {
		os.MkdirAll(config_dir, 0711)
	} else if err != nil {
		// If there is different error from IsNotExists return error
		fmt.Println("Error: ", err)
		return nil, err
	}

	_, err = os.Stat(config_path)
	if os.IsNotExist(err) {
		return os.Create(config_path)
	} else if err != nil {
		return nil, err
	}
	return os.Open(config_path)
}
