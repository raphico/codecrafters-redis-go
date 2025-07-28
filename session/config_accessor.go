package session

type ConfigAccessor interface {
	GetConfig() map[string]string
}
