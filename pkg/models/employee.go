package models

import "time"

type EmployeesResponse struct {
	TotalCount int        `json:"totalCount"`
	Employees  []Employee `json:"employees"`
}

type Employee struct {
	ID              string        `json:"id"`
	CompanyID       string        `json:"companyId"`
	RelantionshipID string        `json:"relantionshipId"`
	LastName        string        `json:"lastName"`
	DisplayName     string        `json:"displayName"`
	Status          string        `json:"status"`
	StatusType      string        `json:"statusType"`
	CurrentStatus   CurrentStatus `json:"currentStatus"`
	Info            Information   `json:"info"`
}

type CurrentStatus struct {
	Status           string    `json:"status"`
	StatusType       string    `json:"statusType"`
	EffectiveDate    time.Time `json:"effectiveDate"`
	ChangeReason     string    `json:"changeReason"`
	ChangeReasonCode int       `json:"changeReasonCode"`
}

type Information struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	DisplayName string    `json:"displayName"`
	MiddleName  string    `json:"middleName"`
	HireDate    time.Time `json:"hireDate"`
	JobTitle    string    `json:"jobTitle"`
}

// Local Variables:
// go-tag-args: ("-transform" "camelcase")
// End:
