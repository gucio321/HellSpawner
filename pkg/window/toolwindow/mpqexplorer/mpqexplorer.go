// Package mpqexplorer contains an implementation of a MPQ archive explorer,
// which displays the archive contents as a tree.
package mpqexplorer

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gucio321/HellSpawner/pkg/app/config"
	"github.com/gucio321/HellSpawner/pkg/app/state"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/widgets"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow"
)

const (
	mainWindowW, mainWindowH = 300, 400
)

// FileSelectedCallback represents file selected callback
type FileSelectedCallback func(path *common.PathEntry)

var _ toolwindow.ToolWindow = (*MPQExplorer)(nil)

// MPQExplorer represents a mpq explorer
type MPQExplorer struct {
	*toolwindow.ToolWindowBase
	config               *config.Config
	project              *hsproject.Project
	fileSelectedCallback FileSelectedCallback
	nodeCache            []g.Widget

	filesToOverwrite []fileToOverwrite
}

type fileToOverwrite struct {
	Path string
	Data []byte
}

// Create creates a new explorer
func Create(fileSelectedCallback FileSelectedCallback, cfg *config.Config, x, y float32) (*MPQExplorer, error) {
	result := &MPQExplorer{
		ToolWindowBase:       toolwindow.New("MPQ Explorer", state.ToolWindowTypeMPQExplorer, x, y),
		fileSelectedCallback: fileSelectedCallback,
		config:               cfg,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result, nil
}

// SetProject sets explorer's project
func (m *MPQExplorer) SetProject(project *hsproject.Project) {
	m.project = project
}

// Build builds an explorer
func (m *MPQExplorer) Build() {
	m.IsOpen(&m.Visible).
		Size(mainWindowW, mainWindowH).
		Layout(m.GetLayout())
}

func (m *MPQExplorer) GetLayout() g.Widget {
	if m.project == nil {
		return g.Label("No project loaded...")
	}

	needToShowOverwritePrompt := len(m.filesToOverwrite) > 0
	if needToShowOverwritePrompt {
		return g.Layout{
			g.PopupModal("Overwrite File?").IsOpen(&needToShowOverwritePrompt).Layout(g.Layout{
				g.Label("File at " + m.filesToOverwrite[0].Path + " already exists. Overwrite?"),
				g.Row(
					g.Button("Overwrite").OnClick(func() {
						success := hsutil.CreateFileAtPath(m.filesToOverwrite[0].Path, m.filesToOverwrite[0].Data)
						if success {
							m.project.InvalidateFileStructure()
						}
						m.filesToOverwrite = m.filesToOverwrite[1:]
					}),
					g.Button("Cancel").OnClick(func() {
						m.filesToOverwrite = m.filesToOverwrite[1:]
					}),
				),
			}),
		}
	}

	return g.Child().
		Border(false).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(m.GetMpqTreeNodes()...)
}

// GetMpqTreeNodes returns mpq tree
func (m *MPQExplorer) GetMpqTreeNodes() []g.Widget {
	if m.nodeCache != nil {
		return m.nodeCache
	}

	wg := sync.WaitGroup{}
	result := make([]g.Widget, len(m.project.AuxiliaryMPQs))
	wg.Add(len(m.project.AuxiliaryMPQs))

	for mpqIndex := range m.project.AuxiliaryMPQs {
		go func(idx int) {
			fullPath := filepath.Join(m.config.AuxiliaryMpqPath, m.project.AuxiliaryMPQs[idx])

			mpq, err := d2mpq.FromFile(fullPath)
			if err != nil {
				log.Printf("failed to load mpq: %s", fullPath)
			}

			if mpq != nil {
				nodes := m.project.GetMPQFileNodes(mpq, m.config)
				result[idx] = m.renderNodes(nodes)
			}

			wg.Done()
		}(mpqIndex)
	}

	wg.Wait()

	m.nodeCache = result

	return result
}

func (m *MPQExplorer) renderNodes(pathEntry *common.PathEntry) g.Widget {
	if !pathEntry.IsDirectory {
		id := generatePathEntryID(pathEntry)

		return g.Layout{
			g.Selectable(pathEntry.Name + id),
			widgets.OnDoubleClick(func() { m.fileSelectedCallback(pathEntry) }),
			g.ContextMenu().Layout(g.Layout{
				g.Selectable("Copy to Project").OnClick(func() {
					m.copyToProject(pathEntry)
				}),
			}),
		}
	}

	nodes := make([]g.Widget, len(pathEntry.Children))
	common.SortPaths(pathEntry)

	wg := sync.WaitGroup{}
	wg.Add(len(pathEntry.Children))

	for childIdx := range pathEntry.Children {
		go func(idx int) {
			nodes[idx] = m.renderNodes(pathEntry.Children[idx])

			wg.Done()
		}(childIdx)
	}

	wg.Wait()

	return g.TreeNode(pathEntry.Name).Layout(nodes...)
}

func (m *MPQExplorer) copyToProject(pathEntry *common.PathEntry) {
	data, err := pathEntry.GetFileBytes()
	if err != nil {
		log.Printf("failed to read file %s when copying to project: %s", pathEntry.FullPath, err)
		return
	}

	pathToFile := pathEntry.FullPath
	if strings.HasPrefix(pathEntry.FullPath, "data") {
		// strip "data" from the beginning of the path if it exists
		pathToFile = pathToFile[4:]
	}

	pathToFile = path.Join(m.project.GetProjectFileContentPath(), pathToFile)
	pathToFile = strings.ReplaceAll(pathToFile, "\\", "/")

	if _, err := os.Stat(pathToFile); err == nil {
		// file already exists
		fileInfo := fileToOverwrite{
			Path: pathToFile,
			Data: data,
		}

		m.filesToOverwrite = append(m.filesToOverwrite, fileInfo)

		return
	}

	success := hsutil.CreateFileAtPath(pathToFile, data)
	if success {
		m.project.InvalidateFileStructure()
	}
}

func generatePathEntryID(pathEntry *common.PathEntry) string {
	return "##MPQExplorerNode_" + pathEntry.FullPath
}

// Reset resets the explorer
func (m *MPQExplorer) Reset() {
	m.nodeCache = nil
}
