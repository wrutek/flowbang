package api

type ColumnItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (column ColumnItem) GetId() int {
	return column.ID
}

func (column ColumnItem) GetName() string {
	return column.Name
}

func (column ColumnItem) GetFullName() string {
	return column.Name
}

type ProjectItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (proj ProjectItem) GetId() int {
	return proj.ID
}

func (proj ProjectItem) GetName() string {
	return proj.Name
}

func (proj ProjectItem) GetFullName() string {
	return proj.Name
}

type RepoItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func (resp RepoItem) GetId() int {
	return resp.ID
}

func (resp RepoItem) GetName() string {
	return resp.Name
}

func (resp RepoItem) GetFullName() string {
	return resp.FullName
}
