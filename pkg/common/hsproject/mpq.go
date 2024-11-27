package hsproject

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/config"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

// GetMPQFileNodes returns mpq's node
func (p *Project) GetMPQFileNodes(mpq d2interface.Archive, cfg *config.Config) *common.PathEntry {
	result := &common.PathEntry{
		Name:        filepath.Base(mpq.Path()),
		IsDirectory: true,
		Source:      common.PathEntrySourceMPQ,
		MPQFile:     mpq.Path(),
	}

	files, err := mpq.Listfile()
	if err != nil {
		files, err = p.searchForMpqFiles(mpq, cfg)
		if err != nil {
			return result
		}
	}

	pathNodes := make(map[string]*common.PathEntry)
	pathNodes[""] = result

	for idx := range files {
		elements := strings.FieldsFunc(files[idx], func(r rune) bool { return r == '\\' || r == '/' })

		path := ""

		for elemIdx := range elements {
			oldPath := path

			path += elements[elemIdx]
			if elemIdx < len(elements)-1 {
				path += `\`
			}

			if pathNodes[strings.ToLower(path)] == nil {
				pathNodes[strings.ToLower(path)] = &common.PathEntry{
					Name:        elements[elemIdx],
					FullPath:    path,
					Source:      common.PathEntrySourceMPQ,
					MPQFile:     mpq.Path(),
					IsDirectory: elemIdx < len(elements)-1,
				}

				pathNodes[strings.ToLower(oldPath)].Children = append(pathNodes[strings.ToLower(oldPath)].Children, pathNodes[strings.ToLower(path)])
			}
		}
	}

	common.SortPaths(result)

	return result
}

// searchForMpqFiles searches for files in MPQ's without listfiles using a list of known filenames
func (p *Project) searchForMpqFiles(mpq d2interface.Archive, cfg *config.Config) ([]string, error) {
	var files []string

	if cfg.ExternalListFile != "" {
		file, err := os.Open(cfg.ExternalListFile)
		if err != nil {
			return files, errors.New("couldn't open listfile")
		}

		defer func() {
			err := file.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileName := scanner.Text()
			if mpq.Contains(fileName) {
				files = append(files, fileName)
			}
		}

		if err := scanner.Err(); err != nil {
			return files, fmt.Errorf("error scanning for file: %w", err)
		}

		return files, nil
	}

	return files, nil
}
