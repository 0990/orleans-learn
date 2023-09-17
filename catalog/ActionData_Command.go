package catalog

import "time"

type CommandDeactivate struct {
}

type CommandActivate struct {
}

type CommandDelay struct {
	Duration time.Duration
}

type CommandUnregisterFromCatalog struct {
}
