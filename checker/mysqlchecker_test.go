package checker

import (
	"availability-checker/credentialprovider"
	"testing"

	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	dct "github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStruct struct {
	mock.Mock
}

func (m *mockStruct) Open(driverName, dataSourceName string) error {
	args := m.Called(driverName, dataSourceName)
	return args.Error(0)
}

func (m *mockStruct) Ping() error {
	return m.Called().Error(0)
}

func (m *mockStruct) Close() error {
	return m.Called().Error(0)
}

func TestMySQLChecker_Check(t *testing.T) {
	// Prepare the mock
	mockConn := new(mockStruct)
	mockConn.On("Open", "mysql", "mockuser:mockpassword@tcp(127.0.0.1:3306)/").Return(nil)
	mockConn.On("Ping").Return(nil)
	mockConn.On("Close").Return(nil)

	// Prepare the checker
	checker := MySQLChecker{
		Server:             "127.0.0.1",
		Port:               "3306",
		DBConnection:       mockConn,
		CredentialProvider: &credentialprovider.MockCredentialProvider{}, // assume you've mocked this too
	}

	// Call the method under test
	_, err := checker.Check()
	assert.Nil(t, err)

	// Assert that the mock was called as expected
	mockConn.AssertExpectations(t)
}

func (m *mockStruct) NewClient() error {
	return m.Called().Error(0)
}

func (m *mockStruct) ContainerList(ctx context.Context, opts types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]types.Container), args.Error(1)
}

func (m *mockStruct) ImagePull(ctx context.Context, image string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, image, opts)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockStruct) ContainerCreate(ctx context.Context, config *dct.Config, hostConfig *dct.HostConfig, containerName string) (dct.CreateResponse, error) {
	args := m.Called(ctx, config, hostConfig, containerName)
	return args.Get(0).(dct.CreateResponse), args.Error(1)
}

func (m *mockStruct) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

// Test the Fix() method
func TestMySQLChecker_Fix(t *testing.T) {
	ctx := context.TODO()

	mockDocker := new(mockStruct)
	mockDocker.On("NewClient").Return(nil)
	mockDocker.On("ContainerList", ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
	mockDocker.On("ImagePull", ctx, "mysql:latest", types.ImagePullOptions{}).Return(ioutil.NopCloser(strings.NewReader("")), nil)
	mockDocker.On("ContainerCreate", ctx, mock.AnythingOfType("*container.Config"), mock.AnythingOfType("*container.HostConfig"), "test-mysql").Return(dct.CreateResponse{ID: "testID"}, nil)
	mockDocker.On("ContainerStart", ctx, "testID", types.ContainerStartOptions{}).Return(nil)

	checker := MySQLChecker{
		containerClient: mockDocker,
		// initialize the other fields as necessary
	}

	err := checker.Fix()

	assert.Nil(t, err)
	mockDocker.AssertExpectations(t)
}
