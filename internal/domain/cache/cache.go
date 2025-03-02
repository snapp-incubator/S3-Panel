package cache

type ServerCache interface {
	Get(targetKey string) (string, error)
	Set(targetKey, targetValue string) error
}
