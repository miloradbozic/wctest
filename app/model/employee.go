package model

import "fmt"

type Employee struct {
	FirstName string 	`json:"firstName"`
	LastName  string 	`json:"lastName"`
	Title     string 	`json:"title"`
	ReportsTo *Employee `json:"reportsTo"`
}

func (e *Employee) FullName() string {
	return fmt.Sprintf("%s %s", e.FirstName, e.LastName)
}

func (e *Employee) ReportsToName() string {
	if e.ReportsTo == nil {
		return ""
	}
	return e.ReportsTo.FullName()
}
