package catalog

var NewGrainInstance func() any

func RegisterNewGrain(NewGrain func() any) {
	NewGrainInstance = NewGrain
}

func CreateNewGrainInstance() any {
	return NewGrainInstance()
}
