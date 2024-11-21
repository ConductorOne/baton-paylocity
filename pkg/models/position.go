package models

import "time"

type Position struct {
	ApprovedDate  *time.Time `json:"approvedDate"`
	EffectiveDate *time.Time `json:"effectiveDate"`
	ClosedDate    *time.Time `json:"closedDate"`
	Code          *string    `json:"code"`
	Title         *string    `json:"title"`
}

// Local Variables:
// go-tag-args: ("-transform" "camelcase")
// End:
