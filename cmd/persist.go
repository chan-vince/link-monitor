package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func PathCreate(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			log.Errorf(err.Error())
			log.Fatalf("Could not create directory: %s", path)
		}
	}
	log.Infof("Created path: %s", path)
	return nil
}

func Save(content string, path string, filename string) error {
	err := PathCreate(path)
	if err != nil {
		log.Error("Failed to save")
		return err
	}
	fp := filepath.Join(path, filename)
	f, err := os.Create(fp)
	defer f.Close()
	_, err = f.WriteString(content + "\n")
	if err != nil {
		log.Error("Failed to write to file")
		return err
	}
	f.Sync()
	return nil
}