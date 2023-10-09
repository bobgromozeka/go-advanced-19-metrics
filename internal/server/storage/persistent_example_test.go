package storage

import (
	"context"
)

func ExamplePersistentStorage() {
	s := NewMemory()
	s = NewPersistenceStorage(
		s, PersistenceSettings{
			Path:     "./pers.json",
			Interval: 300,
			Restore:  false,
		},
	)

	s.SetGauge(context.Background(), "name", 1.111) // Will write this data into file after 300 seconds

	s = NewMemory()
	s = NewPersistenceStorage(
		s, PersistenceSettings{
			Path:     "./pers.json",
			Interval: 0,
			Restore:  false,
		},
	)

	s.SetGauge(context.Background(), "name", 1.111) // Will write this data immediately after update
}
