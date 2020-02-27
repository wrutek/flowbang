package configmodels

// Configuration representation of current configuration
type Configuration struct {
	IssueRepoID        int    `yaml:"issue_repo_id"`
	WorkingRepoID      int    `yaml:"working_repo_id"`
	OauthToken         string `yaml:"oauth_token"`
	ProjectID          int    `yaml:"project_id"`
	TodoColumnID       int    `yaml:"todo_id"`
	InprogressColumnID int    `yaml:"inprogress_id"`
	DoneColumnID       int    `yaml:"done_id"`
}
