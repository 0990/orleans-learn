package catalog

type ActivationState int

const (
	Create ActivationState = iota
	Activating
	Valid
	Deactivating
	Invalid // 5
	FailedToActivate
)
