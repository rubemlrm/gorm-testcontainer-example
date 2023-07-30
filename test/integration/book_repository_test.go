package integration

import (
	"context"
	"gorm-test/internal/models"
	"gorm-test/internal/repositories"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BookRepositoryTestSuite struct {
	testcontainer *TestContainer
	suite.Suite
	DB *gorm.DB
}

func (s *BookRepositoryTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := StartContainer(ctx)
	s.Require().NoError(err)

	s.testcontainer = psqlContainer

	err = RunMigrations(s.testcontainer.DSN)
	s.Require().NoError(err)

	s.DB, err = gorm.Open(postgres.Open(psqlContainer.DSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	s.Require().NoError(err)

}

func (s *BookRepositoryTestSuite) TestWriteAndReadWithSuccess() {

	fakerData := models.Book{}
	err := faker.FakeData(&fakerData)
	s.Require().NoError(err)

	repository := repositories.NewRepository(s.DB)
	err = repository.Insert(&fakerData)
	result, err2 := repository.Get(fakerData.ID)

	assert.Equal(s.T(), fakerData.Title, result.Title)
	assert.NoError(s.T(), err2, nil)
	assert.NoError(s.T(), err, nil)
}

func (s *BookRepositoryTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	s.Require().NoError(s.testcontainer.Terminate(ctx))
}

func TestBookRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookRepositoryTestSuite))
}
