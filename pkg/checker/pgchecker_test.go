package checker

import (
	"availability-checker/pkg/credentialprovider"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
