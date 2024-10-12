package companies

import "errors"

var (
	ErrCompanyNameEmpty          = errors.New("Company name cannot be empty")
	ErrCompanyNameTooLong        = errors.New("Company name cannot contain more than 15 characters")
	ErrCompanyDescriptionTooLong = errors.New("Company description cannot contain more than 3000 characters")
	ErrCompanyEmployeesEmpty     = errors.New("Company must contain at least 1 employee")
	ErrCompanyTypeInvalid        = errors.New("Company type is invalid")
)
