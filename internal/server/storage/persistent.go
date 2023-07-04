package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type PersistenceSettings struct {
	Path     string
	Interval uint
	Restore  bool
}

type PersistentStorage struct {
	Storage
	persistencePath string
	syncPersisting  bool
}

func WithPersistence(config PersistenceSettings) func(Storage) Storage {
	return func(storage Storage) Storage {
		ps := &PersistentStorage{
			Storage:         storage,
			persistencePath: config.Path,
		}

		if config.Path != "" {
			if config.Restore {
				data, restoreErr := restoreFrom(config.Path)
				if restoreErr != nil {
					log.Println("Could not restore data from file: ", restoreErr)
				} else {
					ps.SetMetrics(data)
				}
			}

			if config.Interval == 0 {
				ps.syncPersisting = true
			} else {
				go func() {
					ticker := time.Tick(time.Second * time.Duration(config.Interval))
					for range ticker {
						ps.persist()
					}
				}()
			}
		}

		return ps
	}
}

func (s *PersistentStorage) SetGauge(name string, value float64) float64 {
	res := s.Storage.SetGauge(name, value)

	if s.syncPersisting {
		s.persist()
	}

	return res
}

func (s *PersistentStorage) AddCounter(name string, value int64) int64 {
	res := s.Storage.AddCounter(name, value)

	if s.syncPersisting {
		s.persist()
	}

	return res
}

func (s *PersistentStorage) persist() {
	err := persistToPath(s.persistencePath, Metrics{
		Gauge:   s.GetAllGaugeMetrics(),
		Counter: s.GetAllCounterMetrics(),
	})
	if err != nil {
		fmt.Println("Error during syncing storage data: ", err)
	}
}

func persistToPath(path string, data Metrics) error {
	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	if writeErr := os.WriteFile(path, jsonData, 0666); writeErr != nil {
		return writeErr
	}

	return nil
}

func restoreFrom(filepath string) (Metrics, error) {
	data := Metrics{}

	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return data, err
	}

	unmarshalErr := json.Unmarshal(jsonData, &data)
	if unmarshalErr != nil {
		return data, unmarshalErr
	}

	return data, nil
}
