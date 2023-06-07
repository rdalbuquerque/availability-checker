package credentialprovider

type MockCredentialProvider struct {
}

func (c *MockCredentialProvider) Authenticate() error {
	return nil
}

func (c *MockCredentialProvider) GetCredentials(checker string) (user, password string, err error) {
	return "mockuser", "mockpassword", nil
}
