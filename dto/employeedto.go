package dto

import (
	"time"
)

type Employee struct {
	EmpID         string    `json:"empid,omitempty" bson:"empid,omitempty"`
	Designation   string    `json:"designation,omitempty" bson:"designation,omitempty"`
	ProjectName   string    `json:"projectname,omitempty" bson:"projectname,omitempty"`
	Salary        float64   `json:"salary,omitempty" bson:"salary,omitempty"`
	Experience    float64   `json:"experience,omitempty" bson:"experience,omitempty"`
	Doj           string    `json:"doj,omitempty" bson:"doj,omitempty"`
	DateOfJoining time.Time `json:"dateofjoining,omitempty" bson:"dateofjoining,omitempty"`
	Page          int       `json:"page,omitempty" bson:"page,omitempty"`
	Size          int       `json:"size,omitempty" bson:"size,omitempty"`
}
