package finder

import (
	"os"
	"path/filepath"
	"raygun/log"
)

type Finder struct {
	extension string
}

func NewFinder(pExtension string) *Finder {

	finder := &Finder{extension: pExtension}

	return finder

}

func (f Finder) FindTargets(entities []string) ([]string, error) {

	directories := make([]string, 0)
	target_files := make([]string, 0)

	for _, entity := range entities {

		files, err := filepath.Glob(entity)

		if err != nil {
			log.Error("findTestSuites Error: %s", err)
			return nil, err
		}

		for _, file := range files {

			file_info, err := os.Stat(file)

			if err != nil {
				log.Error("findTestSuites Error: %s", err)
				return nil, err
			}

			if file_info.IsDir() {
				directories = append(directories, file)
			} else if f.isTargetFile(file) {
				target_files = append(target_files, file)
			} else {
				// ignore all other files
			}

		}

	}

	log.Debug("Directories to search for Raygun files: %v", directories)

	for _, dir := range directories {

		file_info_list, err := os.ReadDir(dir)

		if err != nil {
			log.Verbose("Can't read %s, skipping", dir)
			continue
		}

		for _, file_info := range file_info_list {

			if file_info.IsDir() {
				log.Verbose("Skipping subdirectories of %s", dir)
			} else {

				if f.isTargetFile(file_info.Name()) {
					target_files = append(target_files, filepath.Join(dir, file_info.Name()))
				}

			}

		}

	}

	return target_files, nil

}

func (f Finder) isTargetFile(entity string) bool {
	return filepath.Ext(entity) == f.extension
}
