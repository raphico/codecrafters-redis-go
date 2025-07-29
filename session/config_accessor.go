package session

// ConfigAccessor defines the minimal interface session needs from the server
// This ensures session depends only on behavior, not concrete types, avoiding import cycles
type ConfigAccessor interface {
	GetConfig() map[string]string
	GetRDBPath() string
}
