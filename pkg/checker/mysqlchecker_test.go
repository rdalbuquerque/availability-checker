package checker

import (
	"availability-checker/pkg/credentialprovider"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"context"
	"io"

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
	// Test cases
	testCases := []struct {
		name            string
		closeErr        error
		openErr         error
		pingErr         error
		expectedSuccess bool
		expectedErr     error
	}{
		{
			name:            "valid connection",
			closeErr:        nil,
			openErr:         nil,
			pingErr:         nil,
			expectedSuccess: true,
			expectedErr:     nil,
		},
		{
			name:            "connection open error",
			closeErr:        nil,
			openErr:         errors.New("connection open error"),
			pingErr:         nil,
			expectedSuccess: false,
			expectedErr:     errors.New("connection open error"),
		},
		{
			name:            "ping error",
			closeErr:        nil,
			openErr:         nil,
			pingErr:         errors.New("ping error"),
			expectedSuccess: false,
			expectedErr:     errors.New("ping error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare the mock SQLConnection
			mockConn := new(mockStruct)
			mockConn.On("Open", "mysql", "mockuser:mockpassword@tcp(127.0.0.1:3306)/").Return(tc.openErr)
			if tc.openErr == nil {
				mockConn.On("Ping").Return(tc.pingErr)
				mockConn.On("Close").Return(tc.closeErr)
			}

			// Prepare the checker
			checker := MySQLChecker{
				Server:             "127.0.0.1",
				Port:               "3306",
				DBConnection:       mockConn,
				CredentialProvider: &credentialprovider.MockCredentialProvider{},
			}

			// Call the method under test
			success, err := checker.Check()

			// Assert the result
			assert.Equal(t, tc.expectedSuccess, success)
			assert.Equal(t, tc.expectedErr, err)

			// Assert that the mock was called as expected
			mockConn.AssertExpectations(t)
		})
	}
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
	testCases := []struct {
		name           string
		newClientErr   error
		expectedContID string
		containerList  []types.Container
		expectedErr    error
	}{
		{
			name:           "no container",
			expectedContID: "testmysql",
			newClientErr:   nil,
			containerList:  nil,
			expectedErr:    nil,
		},
		{
			name:           "containers exists",
			newClientErr:   nil,
			expectedContID: "testmysql",
			containerList: []types.Container{
				{ID: "testmysql", Names: []string{"/test-mysql"}, Image: "mysql"},
			},
			expectedErr: nil,
		},
	}
	ctx := context.TODO()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDocker := new(mockStruct)
			mockDocker.On("NewClient").Return(nil)
			if tc.newClientErr == nil {
				mockDocker.On("Close").Return(nil)
			}
			mockDocker.On("ContainerList", ctx, types.ContainerListOptions{All: true}).Return(tc.containerList, nil)
			if tc.containerList == nil {
				mockDocker.On("ImagePull", ctx, "mysql:latest", types.ImagePullOptions{}).Return(ioutil.NopCloser(strings.NewReader("")), nil)
				mockDocker.On("ContainerCreate", ctx, mock.AnythingOfType("*container.Config"), mock.AnythingOfType("*container.HostConfig"), "test-mysql").Return(dct.CreateResponse{ID: tc.expectedContID}, nil)
				mockDocker.On("ContainerStart", ctx, tc.expectedContID, types.ContainerStartOptions{}).Return(nil)
			} else {
				mockDocker.On("ContainerStart", ctx, tc.expectedContID, types.ContainerStartOptions{}).Return(nil)
			}

			checker := MySQLChecker{
				ContainerClient: mockDocker,
				// initialize the other fields as necessary
			}

			err := checker.Fix()
			assert.Equal(t, tc.expectedErr, err)
			mockDocker.AssertExpectations(t)
		})
	}
}
