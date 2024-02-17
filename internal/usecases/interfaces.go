//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=usecases
package usecases

type Cache interface {
	AddNewInstance() (int, error)
}
