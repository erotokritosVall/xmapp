package application

import (
	"context"

	domain "github.com/erotokritosVall/xmapp/internal/companies/domain"
)

type companyService struct {
	repo domain.CompanyRepository
}

func NewService(repo domain.CompanyRepository) domain.CompanyService {
	return &companyService{
		repo: repo,
	}
}

func (c *companyService) Delete(ctx context.Context, id string) error {
	return c.repo.Delete(ctx, id)
}

func (c *companyService) Insert(ctx context.Context, company *domain.Company) (*string, error) {
	if err := c.repo.Insert(ctx, company); err != nil {
		return nil, err
	}

	return &company.Id, nil
}

func (c *companyService) Read(ctx context.Context, id string) (*domain.Company, error) {
	return c.repo.Read(ctx, id)
}

func (c *companyService) Update(ctx context.Context, id string, opts *domain.CompanyUpdateOptions) error {
	company, err := c.repo.Read(ctx, id)
	if err != nil {
		return err
	}

	if err := company.Update(opts); err != nil {
		return err
	}

	return c.repo.Update(ctx, id, opts)
}
