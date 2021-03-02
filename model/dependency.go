package model

import "database/sql"

type Dependency struct {
	ID       int `json:"id"`
	ParentID int `json:"parent_id"`
	ChildID  int `json:"child_id"`
}

type DependencyDB struct {
	DB *sql.DB
}

// func GetDependencies(parentID int, childID int) []Dependency {

// }
