package checker

import (
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/winsvcmngr"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgresChecker_Check(t *testing.T) {
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
			mockConn.On("Open", "postgres", "host=127.0.0.1 port=3306 dbname=postgres user=mockuser password=mockpassword sslmode=disable connect_timeout=10").Return(tc.openErr)
			if tc.openErr == nil {
				mockConn.On("Ping").Return(tc.pingErr)
				mockConn.On("Close").Return(tc.closeErr)
			}

			// Prepare the checker
			checker := PostgresChecker{
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

type mockServiceMngr struct {
	mock.Mock
}

func (mngr *mockServiceMngr) Connect() error {
	args := mngr.Called()
	return args.Error(0)
}

func (mngr *mockServiceMngr) Disconnect() error {
	args := mngr.Called()
	return args.Error(0)
}

func (mngr *mockServiceMngr) OpenService(name string) (winsvcmngr.WinSvc, error) {
	args := mngr.Called(name)
	return args.Get(0).(*MockService), args.Error(1)
}

type MockService struct {
	mock.Mock
}

func (svc *MockService) Close() error {
	args := svc.Called()
	return args.Error(0)
}

func (svc *MockService) Start() error {
	args := svc.Called()
	return args.Error(0)
}

func TestPostgresChecker_Fix(t *testing.T) {
	testCases := []struct {
		name        string
		connectErr  error
		openSvcErr  error
		startSvcErr error
		expectedErr error
	}{
		{
			name:        "successfull fix",
			connectErr:  nil,
			openSvcErr:  nil,
			startSvcErr: nil,
			expectedErr: nil,
		},
		{
			name:        "connection error",
			connectErr:  errors.New("error connecting"),
			openSvcErr:  nil,
			startSvcErr: nil,
			expectedErr: errors.New("error connecting"),
		},
		{
			name:        "open service error",
			connectErr:  nil,
			openSvcErr:  errors.New("error opening service"),
			startSvcErr: nil,
			expectedErr: errors.New("error opening service"),
		},
		{
			name:        "starting service error",
			connectErr:  nil,
			openSvcErr:  nil,
			startSvcErr: errors.New("error starting service"),
			expectedErr: errors.New("error starting service"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockWinSvcMngr := new(mockServiceMngr)
			mockSvc := new(MockService)
			mockWinSvcMngr.On("Connect").Return(tc.connectErr)
			if tc.connectErr == nil {
				mockWinSvcMngr.On("Disconnect").Return(nil)
				mockWinSvcMngr.On("OpenService", "postgresql-x64-15").Return(mockSvc, tc.openSvcErr)
				if tc.openSvcErr == nil {
					mockSvc.On("Close").Return(nil)
					mockSvc.On("Start").Return(tc.startSvcErr)
				}
			}

			checker := PostgresChecker{WinSvcMngr: mockWinSvcMngr}

			err := checker.Fix()
			assert.Equal(t, tc.expectedErr, err)
			mockWinSvcMngr.AssertExpectations(t)
		})
	}
}

func TestPostgresChecker_Name(t *testing.T) {
	checker := PostgresChecker{
		Server: "testserver",
		Port:   "1433",
	}

	name := checker.Name()
	assert.Equal(t, fmt.Sprintf("Postgres: %s:%s", checker.Server, checker.Port), name)
}

func TestPostgresChecker_IsFixable(t *testing.T) {
	checker := PostgresChecker{}
	fixable := checker.IsFixable()
	assert.True(t, fixable)
}
