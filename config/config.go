package config

type Config struct {
	Vault struct {
		Address string
		Token   string
	}

  State struct {
    Location string
  }

  Log struct {
		Level  string
		Format string
	}
}
