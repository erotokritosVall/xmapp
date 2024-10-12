package infrastructure

import (
	"context"

	domain "github.com/erotokritosVall/xmapp/internal/companies/domain"
	"github.com/erotokritosVall/xmapp/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "companies"
)

type companyRepository struct {
	db *mongo.Database
}

func New(db *mongo.Database) domain.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

func (c *companyRepository) Delete(ctx context.Context, id string) error {
	where := bson.M{
		"id": id,
	}

	_, err := c.db.Collection(collectionName).DeleteOne(ctx, where)

	return err
}

func (c *companyRepository) Insert(ctx context.Context, company *domain.Company) error {
	_, err := c.db.Collection(collectionName).InsertOne(ctx, company)

	return err
}

func (c *companyRepository) Read(ctx context.Context, id string) (*domain.Company, error) {
	where := bson.M{
		"id": id,
	}

	company := &domain.Company{}
	if err := c.db.Collection(collectionName).FindOne(ctx, where).Decode(company); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}

		return nil, err
	}

	return company, nil
}

func (c *companyRepository) Update(ctx context.Context, id string, opts *domain.CompanyUpdateOptions) error {
	fields := parseUpdateOpts(opts)
	if len(fields) == 0 {
		return nil
	}

	where := bson.M{
		"id": id,
	}

	update := bson.M{
		"$set": fields,
	}

	_, err := c.db.Collection(collectionName).UpdateOne(ctx, where, update)

	return err
}

func parseUpdateOpts(opts *domain.CompanyUpdateOptions) bson.M {
	res := bson.M{}

	if opts.Name != nil {
		res["name"] = *opts.Name
	}

	if opts.Description != nil {
		res["description"] = *opts.Description
	}

	if opts.EmployeeAmount != nil {
		res["employee_amount"] = *opts.EmployeeAmount
	}

	if opts.Registered != nil {
		res["registered"] = *opts.Registered
	}

	if opts.Type != nil {
		res["type"] = *opts.Type
	}

	if len(res) > 0 {
		res["updated_at"] = opts.UpdatedAt
	}

	return res
}
