package companies

import (
	"context"
	"errors"
	"time"

	"github.com/erotokritosVall/xmapp/pkg/util"
	"github.com/google/uuid"
)

const (
	companyNameMaxLength        = 15
	companyDescriptionMaxLength = 3000
)

type Company struct {
	Id             string      `bson:"id"`
	Name           string      `bson:"name"`
	Description    string      `bson:"description"`
	EmployeeAmount int         `bson:"employee_amount"`
	Registered     bool        `bson:"registered"`
	Type           CompanyType `bson:"type"`
	CreatedAt      *time.Time  `bson:"created_at"`
	UpdatedAt      *time.Time  `bson:"updated_at"`
}

func New(name, description string, employeeAmount int, registered bool, cType CompanyType) (*Company, error) {
	c := &Company{
		Id:             uuid.NewString(),
		Name:           name,
		Description:    description,
		EmployeeAmount: employeeAmount,
		Registered:     registered,
		Type:           cType,
		CreatedAt:      util.New(time.Now().UTC()),
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Company) Update(opts *CompanyUpdateOptions) error {
	if c == nil {
		return errors.New("company is nil")
	}

	if opts.Name != nil {
		c.Name = *opts.Name
	}

	if opts.Description != nil {
		c.Description = *opts.Description
	}

	if opts.EmployeeAmount != nil {
		c.EmployeeAmount = *opts.EmployeeAmount
	}

	if opts.Type != nil {
		c.Type = *opts.Type
	}

	now := util.New(time.Now().UTC())
	c.UpdatedAt = now
	opts.UpdatedAt = now

	return c.validate()
}

func (c *Company) validate() error {
	if c == nil {
		return errors.New("company is nil")
	}

	if util.IsEmptyOrWhitespace(c.Name) {
		return ErrCompanyNameEmpty
	}

	if len(c.Name) > companyNameMaxLength {
		return ErrCompanyNameTooLong
	}

	if len(c.Description) > companyDescriptionMaxLength {
		return ErrCompanyDescriptionTooLong
	}

	if c.EmployeeAmount < 1 {
		return ErrCompanyEmployeesEmpty
	}

	if !c.Type.IsACompanyType() {
		return ErrCompanyTypeInvalid
	}

	return nil
}

type CompanyUpdateOptions struct {
	Name           *string
	Description    *string
	EmployeeAmount *int
	Registered     *bool
	Type           *CompanyType
	UpdatedAt      *time.Time
}

type CompanyRepository interface {
	Read(ctx context.Context, id string) (*Company, error)
	Insert(ctx context.Context, company *Company) error
	Update(ctx context.Context, id string, opts *CompanyUpdateOptions) error
	Delete(ctx context.Context, id string) error
}

type CompanyService interface {
	Read(ctx context.Context, id string) (*Company, error)
	Insert(ctx context.Context, company *Company) (*string, error)
	Update(ctx context.Context, id string, opts *CompanyUpdateOptions) error
	Delete(ctx context.Context, id string) error
}
