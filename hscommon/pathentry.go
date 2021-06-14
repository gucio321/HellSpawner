package hscommon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
)

// PathEntrySource represents the type of path entry.
type PathEntrySource int

const (
	// PathEntrySourceMPQ represents a PathEntry that is relative to a specific MPQ.
	PathEntrySourceMPQ PathEntrySource = iota

	// PathEntrySourceProject represents a PathEntry that is relative to the project.
	PathEntrySourceProject

	// PathEntryVirtual represents a PathEntry that is based on the composite view of
	// the project directory and all MPQs (Project first, then MPQs based on load order).
	PathEntryVirtual
)

// PathEntry defines a file/folder
type PathEntry struct {
	OldName     string          `json:"oldMame"`
	Name        string          `json:"name"`
	FullPath    string          `json:"fullPath"`
	MPQFile     string          `json:"mpqFile"`
	Children    []*PathEntry    `json:"children"`
	Source      PathEntrySource `json:"source"`
	IsDirectory bool            `json:"isDirectory"`
	IsRoot      bool            `json:"isRoot"`
	IsRenaming  bool            `json:"isRenaming"`
}

// GetUniqueID returns path's ID
func (p *PathEntry) GetUniqueID() string {
	return fmt.Sprintf("%d_%s_%s", p.Source, p.MPQFile, p.FullPath)
}

// GetFileBytes reads the file and returns the contents
func (p *PathEntry) GetFileBytes() ([]byte, error) {
	if p.Source == PathEntrySourceProject {
		if _, err := os.Stat(p.FullPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot get informations about file %s: %w", p.FullPath, err)
		}

		data, err := ioutil.ReadFile(p.FullPath)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		return data, nil
	}

	mpq, err := d2mpq.FromFile(p.MPQFile)
	if err != nil {
		return nil, fmt.Errorf("error loading file from MPQ: %w", err)
	}

	if mpq.Contains(p.FullPath) {
		data, err := mpq.ReadFile(p.FullPath)
		if err != nil {
			return data, fmt.Errorf("error reading file from mpq: %w", err)
		}

		return data, nil
	}

	return nil, errors.New("could not locate file in mpq")
}

// WriteFile overwrites the file with the given data
func (p *PathEntry) WriteFile(data []byte) error {
	if p.Source != PathEntrySourceProject {
		return errors.New("saving is only supported for files in project, cannot write to MPQs")
	}

	info, err := os.Stat(p.FullPath)
	if err != nil {
		return fmt.Errorf("cannot get informations about file %s: %w", p.FullPath, err)
	}

	err = ioutil.WriteFile(p.FullPath, data, info.Mode())
	if err != nil {
		return fmt.Errorf("cannot write to file at %s: %w", p.FullPath, err)
	}

	return nil
}
