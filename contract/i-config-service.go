package contract

type IConfigService interface {
	Get(key string, value interface{}) error
	GetStruct(value interface{}) error
	Has(key string) (bool, error)
	HasStruct() (bool, error)
}
