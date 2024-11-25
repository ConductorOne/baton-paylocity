package connector

import (
	"fmt"
	"time"

	"github.com/conductorone/baton-paylocity/pkg/models"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

func employees2users(employees []models.Employee, parentResourceID *v2.ResourceId) ([]*v2.Resource, error) {
	var users []*v2.Resource

	for _, employee := range employees {
		user, err := employee2user(employee, parentResourceID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func employee2user(employee models.Employee, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := make(map[string]interface{})

	profile["first_name"] = employee.Info.FirstName
	profile["middle_name"] = employee.Info.MiddleName
	profile["last_name"] = employee.Info.LastName
	profile["status"] = employee.CurrentStatus.Status
	profile["type"] = employee.Position.EmployeeType
	profile["job_title"] = employee.Info.JobTitle
	profile["is_supervisor"] = employee.Info.IsSupervisor
	profile["supervisor"] = employee.Info.Supervisor
	profile["hire_date"] = employee.Info.HireDate.Format(time.RFC3339)
	profile["adjusted_seniority_date"] = employee.Info.AdjustedSeniorityDate.Format(time.RFC3339)
	profile["shift"] = employee.Info.Shift
	profile["eligible_for_hire"] = employee.Info.EligibleForRehire
	profile["position_code"] = employee.Position.PositionCode
	profile["pay_group"] = employee.Info.PayGroup
	profile["annual_salary"] = fmt.Sprintf("%f", employee.CurrentPayRate.AnnualSalary)
	profile["salary"] = fmt.Sprintf("%f", employee.CurrentPayRate.Salary)
	profile["pay_frequency"] = employee.CurrentPayRate.PayFrequency
	profile["pay_grade"] = employee.CurrentPayRate.PayGrade
	profile["pay_type"] = employee.CurrentPayRate.PayType
	profile["base_rate"] = employee.CurrentPayRate.BaseRate
	profile["rate_per"] = employee.CurrentPayRate.RatePer
	profile["default_hours"] = fmt.Sprintf("%f", employee.CurrentPayRate.DefaultHours)
	profile["location"] = fmt.Sprintf("%s, %s, %s %s", employee.Info.Address.Country, employee.Info.Address.State, employee.Info.Address.City, employee.Info.Address.PostalCode)

	traitOptions := []resource.UserTraitOption{
		resource.WithUserProfile(profile),
		resource.WithEmail(employee.Info.PersonalEmail, true),
		resource.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
	}

	ret, err := resource.NewUserResource(
		employee.DisplayName,
		userResourceType,
		employee.ID,
		traitOptions,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot make user resource from employee «%s» (id «%s»), error: %w", employee.DisplayName, employee.ID, err)
	}

	return ret, nil
}
