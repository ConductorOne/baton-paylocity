package models

import "time"

type PositionCodeResponse struct {
	Total int
	Data  []Position
}

type Position struct {
	FlsaOvertimeExempt                 bool                               `json:"flsaOvertimeExempt"`
	Fte                                float64                            `json:"fte"`
	ApprovedDate                       time.Time                          `json:"approvedDate"`
	EffectiveDate                      time.Time                          `json:"effectiveDate"`
	ClosedDate                         time.Time                          `json:"closedDate"`
	SupervisorPosition                 bool                               `json:"supervisorPosition"`
	Code                               string                             `json:"code"`
	Title                              string                             `json:"title"`
	EEOClass                           EEOClass                           `json:"eeoClass"`
	WorkersCompensationCode            WorkersCompensationCode            `json:"workersCompensationCode"`
	CareerLevel                        CareerLevel                        `json:"careerLevel"`
	Families                           []PositionFamily                   `json:"families"`
	StandardOccupationalClassification StandardOccupationalClassification `json:"standardOccupationalClassification"`
}

type EEOClassCompany struct {
	CompanyID string `json:"companyId"`
	Active    bool   `json:"active"`
}

type EEOClass struct {
	EEOClassKey       int               `json:"eeoClassKey"`
	LegacyEEOClassID  int               `json:"legacyEEOClassId"`
	ClientID          int               `json:"clientId"`
	Code              string            `json:"code"`
	Description       string            `json:"description"`
	Active            bool              `json:"active"`
	PositionsAssigned int               `json:"positionsAssigned"`
	EEOClassCompanies []EEOClassCompany `json:"eeoClassCompanies"`
}

type WorkersCompensationCode struct {
	ClientID          int    `json:"clientId"`
	Code              string `json:"code"`
	Description       string `json:"description"`
	Active            bool   `json:"active"`
	PositionsAssigned int    `json:"positionsAssigned"`
}

type CareerLevel struct {
	CareerLevelKey    int    `json:"careerLevelKey"`
	ClientID          int    `json:"clientId"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PositionsAssigned int    `json:"positionsAssigned"`
	Active            bool   `json:"active"`
}

type PositionFamily struct {
	PositionFamilyKey int    `json:"positionFamilyKey"`
	ClientID          int    `json:"clientId"`
	Code              string `json:"code"`
	Active            bool   `json:"active"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PositionsAssigned int    `json:"positionsAssigned"`
}

type StandardOccupationalClassification struct {
	StandardOccupationalClassificationKey int    `json:"standardOccupationalClassificationKey"`
	Group                                 string `json:"group"`
	Code                                  string `json:"code"`
	Title                                 string `json:"title"`
}

// Local Variables:
// go-tag-args: ("-transform" "camelcase")
// End:
