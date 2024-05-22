package disk

import (
	"os"

	"github.com/rs/zerolog/log"
)

const (
	dir      = "data/"
	filepath = dir + "data.rdb"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Saves data to disk.
func Save(data []byte) error {
	log.Info().Str("filepath", filepath).Msg("saving data to disk")

	// rwxrwxrwx
	if !exists(dir) {
		if err := os.Mkdir(dir, 0777); err != nil {
			log.Error().Err(err).Msg("failed to create data dir on disk")
			return err
		}
	}
	// rw-rw-rw
	err := os.WriteFile(filepath, data, 0666)
	if err != nil {
		log.Error().Err(err).Msg("failed to save data to disk")
	}
	return err
}

// Loads data from disk.
// If the data file is not found, bytes returned is `nil`. Be sure to handle this case!
func Load() ([]byte, error) {
	log.Info().Str("filepath", filepath).Msg("loading data from disk")

	if !exists(filepath) {
		log.Info().Str("filepath", filepath).Msg("data does not exist on disk")
		return nil, nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Error().Err(err).Msg("failed to load data from disk")
	}
	return data, err
}
