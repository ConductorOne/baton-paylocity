package models

import "time"

type EmployeesResponse struct {
	TotalCount int        `json:"totalCount"`
	Employees  []Employee `json:"employees"`
}

type Employee struct {
	ID             string           `json:"id"`
	CompanyID      string           `json:"companyID"`
	RelationshipID string           `json:"relationshipID"`
	LastName       string           `json:"lastName"`
	DisplayName    string           `json:"displayName"`
	Status         string           `json:"status"`
	StatusType     string           `json:"statusType"`
	CurrentStatus  Status           `json:"currentStatus"`
	Info           Info             `json:"info"`
	CurrentPayRate PayRate          `json:"currentPayRate"`
	Position       EmployeePosition `json:"position"`
	FuturePayRates []FutureRate     `json:"futurePayRates"`
}

type Status struct {
	Status           string    `json:"status"`
	StatusType       string    `json:"statusType"`
	EffectiveDate    time.Time `json:"effectiveDate"`
	ChangeReason     string    `json:"changeReason"`
	ChangeReasonCode int       `json:"changeReasonCode"`
}

type Info struct {
	FirstName             string       `json:"firstName"`
	LastName              string       `json:"lastName"`
	DisplayName           string       `json:"displayName"`
	MiddleName            string       `json:"middleName"`
	PreferredName         string       `json:"preferredName"`
	Suffix                string       `json:"suffix"`
	Address               Address      `json:"address"`
	HomePhone             string       `json:"homePhone"`
	MobilePhone           string       `json:"mobilePhone"`
	PersonalEmail         string       `json:"personalEmail"`
	SSN                   string       `json:"ssn"`
	DateOfBirth           string       `json:"dateOfBirth"`
	MaritalStatus         string       `json:"maritalStatus"`
	EthnicityRace         string       `json:"ethnicityRace"`
	Gender                string       `json:"gender"`
	HireDate              time.Time    `json:"hireDate"`
	AdjustedSeniorityDate time.Time    `json:"adjustedSeniorityDate"`
	EligibleForRehire     bool         `json:"eligibleForRehire"`
	SupervisorCo          string       `json:"supervisorCo"`
	Supervisor            string       `json:"supervisor"`
	IsSupervisor          bool         `json:"isSupervisor"`
	ReviewerCo            string       `json:"reviewerCo"`
	Reviewer              string       `json:"reviewer"`
	JobTitle              string       `json:"jobTitle"`
	EEOClass              string       `json:"eeoClass"`
	WCC                   string       `json:"wcc"`
	Shift                 string       `json:"shift"`
	ClockBadge            string       `json:"clockBadge"`
	PayGroup              string       `json:"payGroup"`
	OTExempt              bool         `json:"otExempt"`
	Tipped                string       `json:"tipped"`
	MinWageExempt         bool         `json:"minWageExempt"`
	WorkLocation          WorkLocation `json:"workLocation"`
}

type Address struct {
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	County     string `json:"county"`
}

type WorkLocation struct {
	Address     Address `json:"address"`
	Phone       string  `json:"phone"`
	PhoneExt    string  `json:"phoneExt"`
	MobilePhone string  `json:"mobilePhone"`
	Email       string  `json:"email"`
}

type PayRate struct {
	BaseRate       float64   `json:"baseRate"`
	DefaultHours   float64   `json:"defaultHours"`
	Salary         float64   `json:"salary"`
	PayFrequency   string    `json:"payFrequency"`
	PayGrade       string    `json:"payGrade"`
	AnnualSalary   float64   `json:"annualSalary"`
	RatePer        string    `json:"ratePer"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	BeginCheckDate string    `json:"beginCheckDate"`
	IsAutoPay      bool      `json:"isAutoPay"`
	PayType        string    `json:"payType"`
}

type EmployeePosition struct {
	EffectiveDate          time.Time                `json:"effectiveDate"`
	ChangeReason           string                   `json:"changeReason"`
	CostCenter1            string                   `json:"costCenter1"`
	CostCenter2            string                   `json:"costCenter2"`
	CostCenter3            string                   `json:"costCenter3"`
	EmployeeType           string                   `json:"employeeType"`
	PositionCode           string                   `json:"positionCode"`
	PositionDescription    string                   `json:"positionDescription"`
	CareerLevelCode        string                   `json:"careerLevelCode"`
	CareerLevelDescription string                   `json:"careerLevelDescription"`
	PositionFamilies       []EmployeePositionFamily `json:"positionFamilies"`
}

type EmployeePositionFamily struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type FutureRate struct {
	Salary             float64   `json:"salary"`
	PayFrequency       string    `json:"payFrequency"`
	ChangeReason       string    `json:"changeReason"`
	AnnualSalary       float64   `json:"annualSalary"`
	EffectiveDate      time.Time `json:"effectiveDate"`
	BeginCheckDate     string    `json:"beginCheckDate"`
	IsAutoPay          bool      `json:"isAutoPay"`
	PayType            string    `json:"payType"`
	PayRateDescription string    `json:"payRateDescription"`
}

// Local Variables:
// go-tag-args: ("-transform" "camelcase")
// End:
