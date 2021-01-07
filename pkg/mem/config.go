package mem

type MemcachedConfig struct {
	Server   string
	Username string
	Password string
}

func NewConfig(server, username, password string) *MemcachedConfig {
	return &MemcachedConfig{
		Server:   server,
		Username: username,
		Password: password,
	}
}

type Handler struct {
	MemcachedConfig *MemcachedConfig
}

func NewHandler(server, username, password string) *Handler {
	return &Handler{
		MemcachedConfig: &MemcachedConfig{
			Server:   server,
			Username: username,
			Password: password,
		},
	}
}
