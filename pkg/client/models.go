package client

import (
	"fmt"
	"strings"
	"time"
)

type EmployeesResponse struct {
	Employees []*User `json:"employees"`
	NextToken string  `json:"nextToken,omitempty"`
}

type User struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"displayName"`
	Status      string          `json:"statusType"`
	Info        InfoPayload     `json:"info"`
	Position    PositionPayload `json:"position"`
}

type InfoPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"personalEmail"`
	JobTitle  string `json:"jobTitle"`
}

type PositionPayload struct {
	PositionCode string `json:"positionCode"`
	EmployeeType string `json:"employeeType"`
	Department   string `json:"costCenter1"`
}

type Position struct {
	Code          string    `json:"code"`
	Title         string    `json:"title"`
	EffectiveDate time.Time `json:"effectiveDate"`
	ClosedDate    time.Time `json:"closedDate"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Fields  string `json:"fields"`
}

type PaylocityErrorResponse []ErrorResponse

func (e PaylocityErrorResponse) Message() string {
	if len(e) == 0 {
		return "no error details provided"
	}
	var parts []string
	for _, err := range e {
		parts = append(parts, fmt.Sprintf("code: %s, message: %s, field: %s", err.Code, err.Message, err.Fields))
	}
	return strings.Join(parts, " | ")
}
