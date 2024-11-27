// Package projectexplorer provides a project explorer, for viewing project directories as trees.
package projectexplorer

import (
	"github.com/gucio321/HellSpawner/pkg/app/assets"
	"github.com/gucio321/HellSpawner/pkg/app/state"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/OpenDiablo2/dialog"

	"github.com/AllenDang/cimgui-go/imgui"

	g "github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsfiletypes"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/widgets"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow"
)

const (
	mainWindowW, mainWindowH = 300, 400
	btnW, btnH               = 16, 16
	popStyle                 = 2
	pushStyle                = 4
)

const (
	blackHalfOpacity = 0xffffff20
)

// FileSelectedCallback represents callback on project file selected
type FileSelectedCallback func(path *common.PathEntry)

// ProjectExplorer represents a project explorer
type ProjectExplorer struct {
	*toolwindow.ToolWindowBase

	project              *hsproject.Project
	fileSelectedCallback FileSelectedCallback
	nodeCache            map[string][]g.Widget
	refreshIconTexture   *g.Texture
}

// Create creates a new project explorer
func Create(textureLoader common.TextureLoader,
	fileSelectedCallback FileSelectedCallback,
	x, y float32,
) (*ProjectExplorer, error) {
	result := &ProjectExplorer{
		ToolWindowBase:       toolwindow.New("Project Explorer", state.ToolWindowTypeProjectExplorer, x, y),
		nodeCache:            make(map[string][]g.Widget),
		fileSelectedCallback: fileSelectedCallback,
	}

	result.Visible = false

	// some type of workaround ;-). SOmetimes we only want to get tree nodes (and don't need textures)
	if textureLoader != nil {
		textureLoader.CreateTextureFromFile(assets.ReloadIcon, func(texture *g.Texture) {
			result.refreshIconTexture = texture
		})
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result, nil
}

// SetProject sets explored project
func (m *ProjectExplorer) SetProject(project *hsproject.Project) {
	m.project = project
}

// Build builds explorer
func (m *ProjectExplorer) Build() {
	if m.project == nil {
		return
	}

	header := g.Row(
		m.makeRefreshButtonLayout(),
	)

	tree := g.Child().
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(m.GetProjectTreeNodes())

	m.IsOpen(&m.Visible).
		Layout(g.Layout{
			header,
			g.Separator(),
			tree,
		})
}

func (m *ProjectExplorer) makeRefreshButtonLayout() g.Layout {
	button := g.ImageButton(m.refreshIconTexture).
		Size(btnW, btnH).
		OnClick(func() {
			m.onRefreshProjectExplorerClicked()
		})

	const tooltipText = "Refresh the view from the filesystem."

	if m.project == nil {
		button.TintColor(hsutil.Color(blackHalfOpacity))
	}

	return g.Layout{
		g.Custom(func() {
			imgui.PushStyleColorVec4(imgui.ColButton, imgui.Vec4{})
			imgui.PushStyleColorVec4(imgui.ColBorder, imgui.Vec4{})
			imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{Y: pushStyle})
			imgui.PushIDStr("ProjectExplorerRefresh")
		}),

		button,

		g.Tooltip(tooltipText),

		g.Custom(func() {
			imgui.PopID()
			imgui.PopStyleVar()
			imgui.PopStyleColorV(popStyle)
		}),
	}
}

// GetProjectTreeNodes returns project tree
func (m *ProjectExplorer) GetProjectTreeNodes() g.Layout {
	if m.project == nil {
		return g.Layout{g.Label("No project loaded...")}
	}

	fileStructure, err := m.project.GetFileStructure()
	if err != nil {
		log.Print(err)
	}

	if fileStructure == nil {
		return g.Layout{g.Label("No file structure detected...")}
	}

	nodes, err := m.project.GetFileStructure()
	if err != nil {
		return g.Layout{g.Label(err.Error())}
	}

	return g.Layout{m.renderNodes(nodes)}
}

func (m *ProjectExplorer) onRefreshProjectExplorerClicked() {
	if m.project == nil {
		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onNewFontClicked(pathEntry *common.PathEntry) {
	if err := m.project.CreateNewFile(hsfiletypes.FileTypeFont, pathEntry); err != nil {
		log.Print(err)
	}
}

func (m *ProjectExplorer) renderNodes(pathEntry *common.PathEntry) g.Widget {
	if !pathEntry.IsDirectory {
		return m.createFileTreeItem(pathEntry)
	}

	// File items and empty dirs
	if len(pathEntry.Children) == 0 {
		return m.createDirectoryTreeItem(pathEntry, nil)
	}

	nodes := make([]g.Widget, len(pathEntry.Children))

	sortPaths(pathEntry)

	for idx := range pathEntry.Children {
		nodes[idx] = m.renderNodes(pathEntry.Children[idx])
	}

	return m.createDirectoryTreeItem(pathEntry, nodes)
}

func (m *ProjectExplorer) createFileTreeItem(pathEntry *common.PathEntry) g.Widget {
	id := "##ProjectExplorerNode_" + pathEntry.FullPath

	var layout g.Layout = make([]g.Widget, 0)

	if pathEntry.IsRenaming {
		layout = g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextWithHint("##RenameField_"+pathEntry.FullPath, "", &pathEntry.Name,
					imgui.InputTextFlagsAutoSelectAll|imgui.InputTextFlagsEnterReturnsTrue, nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry)
				}
			}),
		}
	} else {
		layout = append(layout,
			g.Selectable(pathEntry.Name+id),
			widgets.OnDoubleClick(func() { m.fileSelectedCallback(pathEntry) }),
		)
	}

	layout = append(layout,
		g.ContextMenu().Layout(g.Layout{
			g.MenuItem("Rename").OnClick(func() { m.onRenameFileClicked(pathEntry) }),
			g.MenuItem("Delete...").OnClick(func() { m.onDeleteFileClicked(pathEntry) }),
		}),
	)

	return layout
}

func (m *ProjectExplorer) createDirectoryTreeItem(pathEntry *common.PathEntry, layout g.Layout) g.Widget {
	id := pathEntry.Name + "##ProjectExplorerNode_" + pathEntry.FullPath

	if pathEntry.IsRenaming {
		return g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextWithHint("##RenameField_"+pathEntry.FullPath, "", &pathEntry.Name,
					imgui.InputTextFlagsAutoSelectAll|imgui.InputTextFlagsEnterReturnsTrue, nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry)
				}
			}),
		}
	}

	contextMenuLayout := g.Layout{
		g.Menu("New").Layout(g.Layout{
			g.MenuItem("Folder").OnClick(func() { m.onNewFolderClicked(pathEntry) }),
			g.Separator(),
			g.MenuItem("Font").OnClick(func() { m.onNewFontClicked(pathEntry) }),
			g.MenuItem("Font table (.tbl)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeTBLFontTable, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("String table (.tbl)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeTBLStringTable, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Animation data (.d2)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeAnimationData, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Animation (.cof)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeCOF, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Palette (.dat)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypePalette, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Palette transform (.pl2)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypePL2, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Map tile data (.ds1)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeDS1, pathEntry); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Map tile animation (.dt1)").OnClick(func() {
				if err := m.project.CreateNewFile(hsfiletypes.FileTypeDT1, pathEntry); err != nil {
					log.Print(err)
				}
			}),
		}),
	}

	if !pathEntry.IsRoot {
		contextMenuLayout = append(contextMenuLayout,
			g.Separator(),
			g.MenuItem("Rename").OnClick(func() { m.onRenameFileClicked(pathEntry) }),
			g.MenuItem("Delete Folder...").OnClick(func() { m.onDeleteFolderClicked(pathEntry) }),
		)
	}

	menuLayout := g.Layout{
		g.Custom(func() { imgui.PushIDStr(id) }),
		g.ContextMenu().Layout(contextMenuLayout),
		g.Custom(func() { imgui.PopID() }),
	}

	if layout == nil {
		return g.TreeNode(id).Layout(menuLayout)
	}

	return g.TreeNode(id).Layout(append(menuLayout, layout...))
}

func (m *ProjectExplorer) onDeleteFolderClicked(entry *common.PathEntry) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}

	if err := os.RemoveAll(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()

		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onDeleteFileClicked(entry *common.PathEntry) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}

	if err := os.Remove(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()

		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onRenameFileClicked(entry *common.PathEntry) {
	entry.OldName = entry.Name
	entry.IsRenaming = true
}

func (m *ProjectExplorer) onFileRenamed(entry *common.PathEntry) {
	if entry.Name == entry.OldName {
		entry.OldName = ""

		return
	}

	if entry.Name == "" {
		dialog.Message("Cannot rename file:\nFiles cannot have a blank name.").Error()

		entry.Name = entry.OldName

		entry.OldName = ""

		return
	}

	if filepath.Ext(entry.Name) == "" {
		entry.Name += filepath.Ext(entry.OldName)
	}

	if !strings.EqualFold(filepath.Ext(entry.OldName), filepath.Ext(entry.Name)) {
		dialog.Message("Cannot rename file:\nFile extension cannot be changed.").Error()

		entry.Name = entry.OldName

		entry.OldName = ""

		return
	}

	basePath := filepath.Dir(entry.FullPath)

	oldPath := filepath.Join(basePath, entry.OldName)

	newPath := filepath.Join(basePath, entry.Name)

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		dialog.Message("Cannot rename file:\nAlready exists.").Error()

		entry.Name = entry.OldName
		entry.OldName = ""

		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		dialog.Message("Could not rename file: %v", err).Error()

		entry.Name = entry.OldName
		entry.OldName = ""

		return
	}

	m.project.InvalidateFileStructure()
}

func logErr(fmtErr string, args ...interface{}) {
	log.Printf(fmtErr, args...)
	dialog.Message(fmtErr, args...).Error()
}

func (m *ProjectExplorer) onNewFolderClicked(pathEntry *common.PathEntry) {
	if err := m.project.CreateNewFolder(pathEntry); err != nil {
		logErr("%s", err)
	}
}

func sortPaths(rootPath *common.PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if rootPath.Children[i].IsDirectory == rootPath.Children[j].IsDirectory {
			var nameI, nameJ string

			if rootPath.Children[i].OldName != "" {
				nameI = rootPath.Children[i].OldName
			} else {
				nameI = rootPath.Children[i].Name
			}

			if rootPath.Children[j].OldName != "" {
				nameJ = rootPath.Children[j].OldName
			} else {
				nameJ = rootPath.Children[j].Name
			}

			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		}

		return rootPath.Children[i].IsDirectory && !rootPath.Children[j].IsDirectory
	})
}
