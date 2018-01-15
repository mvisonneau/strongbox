package config

// Config : A type that handles the configuration of the app
type Config struct {
	Vault struct {
		Address  string
		Token    string
		RoleID   string
		SecretID string
	}

	State struct {
		Location string
	}

	Log struct {
		Level  string
		Format string
	}
}
