package disk

import (
	"os"

	"github.com/rs/zerolog/log"
)

const (
	dir      = "data/"
	filepath = dir + "data.rdb"
)

func Save(data []byte) error {
	log.Info().Str("filepath", filepath).Msg("saving data to disk")

	// rwxrwxrwx
	if _, err := os.Stat(dir); os.IsNotExist(err) {
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

func Load() ([]byte, error) {
	log.Info().Str("filepath", filepath).Msg("loading data from disk")

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Error().Err(err).Msg("failed to load data from disk")
	}
	return data, err
}
