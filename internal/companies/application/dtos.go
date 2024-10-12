package application

import (
	"time"

	domain "github.com/erotokritosVall/xmapp/internal/companies/domain"
)

type insertCompanyRequest struct {
	Name           string `json:"name" validate:"required,max=15"`
	Description    string `json:"description" validate:"omitempty,max=3000"`
	EmployeeAmount int    `json:"employeeAmount" validate:"gt=0"`
	Registered     bool   `json:"registered"`
	Type           int    `json:"type" validate:"oneof=0 1 2 3"`
}

type updateCompanyRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	EmployeeAmount *int    `json:"employeeAmount"`
	Registered     *bool   `json:"registered"`
	Type           *int    `json:"type"`
}

type readCompanyDto struct {
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	EmployeeAmount int        `json:"employeeAmount"`
	Registered     bool       `json:"registered"`
	Type           int        `json:"type"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

func (r *insertCompanyRequest) toDomain() (*domain.Company, error) {
	return domain.New(r.Name, r.Description, r.EmployeeAmount, r.Registered, domain.CompanyType(r.Type))
}

func (r *updateCompanyRequest) toDomain() *domain.CompanyUpdateOptions {
	return &domain.CompanyUpdateOptions{
		Name:           r.Name,
		Description:    r.Description,
		EmployeeAmount: r.EmployeeAmount,
		Registered:     r.Registered,
		Type:           (*domain.CompanyType)(r.Type),
	}
}

func companyToReadDto(c *domain.Company) *readCompanyDto {
	return &readCompanyDto{
		Id:             c.Id,
		Name:           c.Name,
		Description:    c.Description,
		EmployeeAmount: c.EmployeeAmount,
		Registered:     c.Registered,
		Type:           int(c.Type),
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
