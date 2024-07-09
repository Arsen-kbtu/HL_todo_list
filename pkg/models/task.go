package models

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title" validate:"required,max=200"`
	ActiveAt string `json:"activeAt" validate:"required,datetime=2006-01-02"`
	Done     bool   `json:"done"`
}
