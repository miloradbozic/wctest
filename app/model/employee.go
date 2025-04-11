package model

type Employee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	ManagerID *int   `json:"manager_id"`
}
