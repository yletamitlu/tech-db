package additions

type AdditionRepository interface {
	Status() (uint64, uint64, uint64, uint64, error)
	Clear() error
}
