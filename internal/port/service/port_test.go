package service_test

import (
	"context"
	"github.com/boriszhilko/ports-service/internal/port/service"
	"io"
	"strings"
	"testing"

	"github.com/boriszhilko/ports-service/internal/port/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	svc      service.Service
	input    *mockInput
	repo     *mockRepository
	testData = `{"AEAJM": {"name": "Ajman","city": "Ajman","country": "United Arab Emirates","alias": [],"regions": []}}`
)

type mockInput struct {
	mock.Mock
}

func (m *mockInput) GetData() io.Reader {
	args := m.Called()
	return args.Get(0).(io.Reader)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateOrUpdate(ctx context.Context, ports []domain.Port) error {
	args := m.Called(ctx, ports)
	return args.Error(0)
}

func (m *mockRepository) Get(ctx context.Context, id string) (domain.Port, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Port), args.Error(1)
}

func setup() {
	var err error
	input = new(mockInput)
	repo = new(mockRepository)
	svc, err = service.NewPortService(input, repo)
	if err != nil {
		panic(err)
	}
}

func TestNewPortService_RepositoryRequired(t *testing.T) {
	// GIVEN
	setup()
	// WHEN
	_, err := service.NewPortService(input, nil)
	// THEN
	assert.EqualError(t, err, service.ErrRepositoryRequired.Error())
}

func TestNewPortService_InputRequired(t *testing.T) {
	// GIVEN
	setup()
	// WHEN
	_, err := service.NewPortService(nil, repo)
	// THEN
	assert.EqualError(t, err, service.ErrInputRequired.Error())
}

func TestPortService_CreateOrUpdatePorts(t *testing.T) {
	// GIVEN
	setup()

	reader := strings.NewReader(testData)
	input.On("GetData").Return(reader)

	repo.On("CreateOrUpdate", mock.Anything, mock.Anything).Return(nil)
	repo.On("Get", mock.Anything, mock.Anything).Return(domain.Port{}, nil)

	// WHEN
	err := svc.CreateOrUpdatePorts(context.Background())
	// THEN
	assert.NoError(t, err)
	input.AssertExpectations(t)
	repo.AssertExpectations(t)
}
