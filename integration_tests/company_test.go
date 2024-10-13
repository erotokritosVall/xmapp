package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/erotokritosVall/xmapp/integration_tests/util"
	companiesApp "github.com/erotokritosVall/xmapp/internal/companies/application"
	companiesDomain "github.com/erotokritosVall/xmapp/internal/companies/domain"
	companiesInfra "github.com/erotokritosVall/xmapp/internal/companies/infrastructure"
	mg "github.com/erotokritosVall/xmapp/pkg/mongo"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	clientTimeout = 5 * time.Second
)

type insertCompanyRespone struct {
	Data  string `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type CompanyTestSuit struct {
	suite.Suite
	ctx         context.Context
	testConfig  *util.TestsConfig
	db          *mongo.Database
	companyRepo companiesDomain.CompanyRepository
	client      *http.Client
	authToken   string
}

func (suite *CompanyTestSuit) SetupSuite() {
	suite.ctx = context.Background()

	cfg := &util.TestsConfig{}

	err := envconfig.Process("", cfg)
	suite.NoError(err, "failed to process envconfig")
	suite.testConfig = cfg

	db, err := mg.New(suite.ctx, suite.testConfig.MongoConfig)
	suite.NoError(err, "failed to initialize mongo")
	suite.db = db

	repo := companiesInfra.New(db)
	suite.companyRepo = repo

	suite.client = &http.Client{
		Timeout: clientTimeout,
	}

	token, err := util.Login(suite.ctx, suite.testConfig)
	suite.NoError(err, "failed to login")

	if token != nil {
		suite.authToken = *token
	}
}

func (suite *CompanyTestSuit) TearDownSuite() {
}

func (suite *CompanyTestSuit) SetupTest() {
}

func (suite *CompanyTestSuit) TearDownTest() {
}

func (suite *CompanyTestSuit) TestInsertCompany() {
	req := companiesApp.InsertCompanyRequest{
		Name:           util.RandString(10),
		Description:    "testDescription",
		EmployeeAmount: 20,
		Registered:     true,
		Type:           1,
	}

	body, err := json.Marshal(&req)
	suite.NoError(err, "failed to marshal InsertCompanyRequest")

	reader := bytes.NewReader(body)
	url := fmt.Sprintf("http://%s:%s/v1/companies", suite.testConfig.AppHost, suite.testConfig.AppPort)
	httpReq, err := http.NewRequestWithContext(suite.ctx, http.MethodPost, url, reader)
	suite.NoError(err, "failed to create http requuest")

	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", suite.authToken))

	resp, err := suite.client.Do(httpReq)
	suite.NoError(err, "failed to send http request")
	defer resp.Body.Close()

	suite.Assert().Equal(http.StatusCreated, resp.StatusCode, "response status should be 201")

	respBody := &insertCompanyRespone{}
	err = json.NewDecoder(resp.Body).Decode(respBody)
	suite.NoError(err, "failed to unmarshal insertCompanyRespone")

	suite.Assert().NotEmpty(respBody.Data)

	company, err := suite.companyRepo.Read(suite.ctx, respBody.Data)
	suite.NoError(err, "failed to read company from repository")
	suite.Assert().NotNil(company, "company should not be nil")
	suite.Assert().Equal(req.Name, company.Name)
	suite.Assert().Equal(req.Description, company.Description)
	suite.Assert().Equal(req.EmployeeAmount, company.EmployeeAmount)
	suite.Assert().Equal(req.Type, int(company.Type))
}

func TestCompanies(t *testing.T) {
	suite.Run(t, new(CompanyTestSuit))
}
