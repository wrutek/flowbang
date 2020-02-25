package api

// ColumnItem struct representation of column
// in the json response from github
type ColumnItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetID get column ID
func (column ColumnItem) GetID() int {
	return column.ID
}

// GetName get column name
func (column ColumnItem) GetName() string {
	return column.Name
}

// GetFullName get column full_name
func (column ColumnItem) GetFullName() string {
	return column.Name
}

// ProjectItem struct representation of project
// in the json response from github
type ProjectItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetID get project ID
func (proj ProjectItem) GetID() int {
	return proj.ID
}

// GetName get project name
func (proj ProjectItem) GetName() string {
	return proj.Name
}

// GetFullName get project full_name
func (proj ProjectItem) GetFullName() string {
	return proj.Name
}

// RepoItem struct representation of repository
// in the json response from github
type RepoItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

// GetID get repository ID
func (resp RepoItem) GetID() int {
	return resp.ID
}

// GetName get repository name
func (resp RepoItem) GetName() string {
	return resp.Name
}

// GetFullName get repository full_name
func (resp RepoItem) GetFullName() string {
	return resp.FullName
}
