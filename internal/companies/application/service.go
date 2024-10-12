package application

import (
	"context"

	domain "github.com/erotokritosVall/xmapp/internal/companies/domain"
	pubsub "github.com/erotokritosVall/xmapp/internal/companies/pub_sub"
	"github.com/erotokritosVall/xmapp/internal/events"
	"github.com/rs/zerolog/log"
)

type companyService struct {
	repo      domain.CompanyRepository
	publisher *pubsub.PublisherManager
}

func NewService(repo domain.CompanyRepository, publisher *pubsub.PublisherManager) domain.CompanyService {
	return &companyService{
		repo:      repo,
		publisher: publisher,
	}
}

func (c *companyService) Delete(ctx context.Context, id string) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	toPublish := []events.DomainEvent{events.NewCompanyDeleted(id)}
	if err := c.publisher.PublishDomainEvents(ctx, toPublish); err != nil {
		log.Err(err).Msg("failed to publish CompanyDeleted event")
	}

	return nil
}

func (c *companyService) Insert(ctx context.Context, company *domain.Company) (*string, error) {
	if err := c.repo.Insert(ctx, company); err != nil {
		return nil, err
	}

	toPublish := []events.DomainEvent{events.NewCompanyCreated(company.Id)}
	if err := c.publisher.PublishDomainEvents(ctx, toPublish); err != nil {
		log.Err(err).Msg("failed to publish CompanyCreated event")
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

	if err := c.repo.Update(ctx, id, opts); err != nil {
		return err
	}

	toPublish := []events.DomainEvent{events.NewCompanyUpdated(id)}
	if err := c.publisher.PublishDomainEvents(ctx, toPublish); err != nil {
		log.Err(err).Msg("failed to publish CompanyUpdated event")
	}

	return nil
}
