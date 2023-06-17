package credentialprovider

type CredentialProvider interface {
	Authenticate() error
	GetCredentials(checkerType string) (user, password string, err error)
}
