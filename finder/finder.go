/*
Copyright © 2024 PACLabs
*/
package finder

/*
 *   Simple code for iterating over a command line file specification that may or may
 *   not include wildcards and looking for .raygun files
 */

import (
	"os"
	"path/filepath"
	"raygun/log"
)

type Finder struct {
	extension string
}

/*
 *  By default, the extension is .raygun
 */
func NewFinder(pExtension string) *Finder {

	finder := &Finder{extension: pExtension}

	return finder

}

/*
 *  This is a two step process - we get a set of command line arguments.  Each one may
 *  be a file, a directory or a glob.
 *
 *  We work through these files and directories, looking for .raygun files, and build
 *  out a complete list of every .raygun file we find.
 */
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

			//
			// this could potentially be overruled by a flag --skip-on-file-error
			// or something like that.  TBD
			//
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

	log.Verbose("Directories to search for %s files: %v", f.extension, directories)

	for _, dir := range directories {

		file_info_list, err := os.ReadDir(dir)

		if err != nil {
			log.Verbose("Can't open directory %s, skipping", dir)
			continue
		}

		for _, file_info := range file_info_list {

			// this might be a place for another flag, indicating some sort of recursive
			// directory search. The current structure wouldn't allow this, so we'll have
			// to rewrite this section if it needs to support recursive searches
			//
			if file_info.IsDir() {
				log.Debug("Skipping subdirectories of %s", dir)
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
