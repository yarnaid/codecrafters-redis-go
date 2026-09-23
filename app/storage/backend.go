package storage

type Backend interface {
	Get(key string) (Value, bool)
	Set(key string, v *Value)
	Delete(key string)
}
