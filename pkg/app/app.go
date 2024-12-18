package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/abysswrapper"
	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsfiletypes"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
	"github.com/gucio321/HellSpawner/pkg/window/popup/aboutdialog"
	"github.com/gucio321/HellSpawner/pkg/window/popup/preferences"
	"github.com/gucio321/HellSpawner/pkg/window/popup/projectproperties"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/console"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/mpqexplorer"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/projectexplorer"
)

const (
	baseWindowTitle          = "HellSpawner"
	baseWindowW, baseWindowH = 1280, 720
	editorWindowDefaultX     = 320
	editorWindowDefaultY     = 30
	projectExplorerDefaultX  = 0
	projectExplorerDefaultY  = 25
	mpqExplorerDefaultX      = 30
	mpqExplorerDefaultY      = 30
	consoleDefaultX          = 10
	consoleDefaultY          = 500

	samplesPerSecond = 22050
	sampleDuration   = time.Second / 10

	autoSaveTimer = 120

	logFileSeparator = "-----%v-----\n"
	logFilePerms     = 0o644
)

const (
	baseFontSize          = 17
	fixedFontSize         = 15
	fixedSmallFontSize    = 12
	diabloRegularFontSize = 15
	diabloBoldFontSize    = 30
)

type editorConstructor func(
	config *config.Config,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte,
	x, y float32,
	project *hsproject.Project,
) (editor.Editor, error)

// App represents an app
type App struct {
	masterWindow *g.MasterWindow
	*Flags
	project      *hsproject.Project
	config       *config.Config
	abyssWrapper *abysswrapper.AbyssWrapper
	logFile      *os.File

	aboutDialog             *aboutdialog.AboutDialog
	preferencesDialog       *preferences.Dialog
	projectPropertiesDialog *projectproperties.Dialog

	projectExplorer *projectexplorer.ProjectExplorer
	mpqExplorer     *mpqexplorer.MPQExplorer
	console         *console.Console

	editors            []editor.Editor
	editorConstructors map[hsfiletypes.FileType]editorConstructor

	editorManagerMutex sync.RWMutex
	focusedEditor      editor.Editor

	fontFixed         *g.FontInfo
	fontFixedSmall    *g.FontInfo
	diabloBoldFont    *g.FontInfo
	diabloRegularFont *g.FontInfo

	showUsage   bool
	justStarted bool
}

// Create creates new app instance
func Create() (*App, error) {
	result := &App{
		Flags:              &Flags{},
		editors:            make([]editor.Editor, 0),
		editorConstructors: make(map[hsfiletypes.FileType]editorConstructor),
		abyssWrapper:       abysswrapper.Create(),
		justStarted:        true,
	}

	if shouldTerminate := result.parseArgs(); shouldTerminate {
		return nil, nil
	}

	result.config = config.Load(*result.Flags.optionalConfigPath)

	return result, nil
}

// Run runs an app instance
func (a *App) Run() (err error) {
	defer a.Quit() // force-close and save everything (in case of crash)

	// setting up the logging here, as opposed to inside of app.setup(),
	// because of the deferred call to logfile.Close()
	if a.config.LoggingToFile || *a.Flags.logFile != "" {
		path := a.config.LogFilePath
		if *a.Flags.logFile != "" {
			path = *a.Flags.logFile
		}

		a.logFile, err = os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFilePerms)
		if err != nil {
			logErr("Error opening log file at %s: %v", a.config.LogFilePath, err)
		}

		defer func() {
			if logErr := a.logFile.Close(); logErr != nil {
				log.Fatal(logErr)
			}
		}()
	}

	dialog.Init()
	a.setupMasterWindow()

	a.masterWindow.Run(a.render)

	return nil
}

func (a *App) render() {
	// unfortunately can't do that in Run as this requires imgui.MainViewport
	if a.justStarted {
		a.justStarted = false

		fmt.Println("start setup")
		err := a.setup()
		if err != nil {
			logErr("could not set up application: %v", err)
		}

		if a.config.OpenMostRecentOnStartup && len(a.config.RecentProjects) > 0 {
			err = a.loadProjectFromFile(a.config.RecentProjects[0])
			if err != nil {
				logErr("could not load most recent project on startup: %v", err)
			}
		}
	}

	switch a.config.ViewMode {
	case config.ViewModeLegacy:
		a.renderLegacy()
	case config.ViewModeStatic:
		a.renderStatic()
	}
}

func (a *App) renderLegacy() {
	g.MainMenuBar().Layout(a.menuLayout()).Build()

	a.renderEditors()
	a.renderWindows()

	g.Update()
}

func (a *App) renderStatic() {
	g.SingleWindowWithMenuBar().Layout(
		g.MenuBar().Layout(a.menuLayout()),
		renderWnd(a.preferencesDialog),
		renderWnd(a.aboutDialog),
		renderWnd(a.projectPropertiesDialog),
		g.SplitLayout(g.DirectionVertical, &a.config.StaticLayout.ProjectSplit,
			a.projectExplorer.GetLayout(),
			g.SplitLayout(g.DirectionVertical, &a.config.StaticLayout.MPQSplit,
				g.SplitLayout(g.DirectionHorizontal, &a.config.StaticLayout.ConsoleSplit,
					a.renderStaticEditors(),
					a.console.GetLayout(),
				).SplitRefType(g.SplitRefProc),
				a.mpqExplorer.GetLayout(),
			).SplitRefType(g.SplitRefProc),
		).SplitRefType(g.SplitRefProc),
	)
}

func logErr(fmtErr string, args ...interface{}) {
	log.Printf(fmtErr, args...)
	dialog.Message(fmtErr, args...).Error()
}

func (a *App) createEditor(path *common.PathEntry, state []byte, x, y, w, h float32) {
	data, err := path.GetFileBytes()
	if err != nil {
		const fmtErr = "Could not load file: %v"

		logErr(fmtErr, err)

		return
	}

	fileType, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &data)
	if err != nil {
		const fmtErr = "Error reading file type: %v"

		logErr(fmtErr, err)

		return
	}

	if a.editorConstructors[fileType] == nil {
		const fmtErr = "Error opening editor: %v"

		logErr(fmtErr, err)

		return
	}

	editor, err := a.editorConstructors[fileType](a.config, path, state, &data, x, y, a.project)
	if err != nil {
		const fmtErr = "Error creating editor: %v"

		logErr(fmtErr, err)

		return
	}

	editor.Size(w, h)

	a.editors = append(a.editors, editor)
	editor.Show()
	editor.BringToFront()
}

func (a *App) openEditor(path *common.PathEntry) {
	a.editorManagerMutex.RLock()

	uniqueID := path.GetUniqueID()
	for idx := range a.editors {
		if a.editors[idx].GetID() == uniqueID {
			a.editors[idx].BringToFront()
			a.editorManagerMutex.RUnlock()

			return
		}
	}

	a.editorManagerMutex.RUnlock()

	// since we sue multiviewport, we need to get base position of the main window - we want an editor
	// inside window
	basePos := imgui.MainViewport().Pos()

	// w, h = 0, because we're creating a new editor,
	// width and height aren't saved, so we give 0 and
	// editors without AutoResize flag sets w, h to default
	a.editorManagerMutex.Lock()
	a.createEditor(path, nil, editorWindowDefaultX+basePos.X, editorWindowDefaultY+basePos.Y, 0, 0)
	a.editorManagerMutex.Unlock()
}

func (a *App) loadProjectFromFile(file string) error {
	project, err := hsproject.LoadFromFile(file)
	if err != nil {
		return fmt.Errorf("could not load project from file %s, %w", file, err)
	}

	err = project.ValidateAuxiliaryMPQs(a.config)
	if err != nil {
		return fmt.Errorf("could not validate aux mpq's, %w", err)
	}

	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()

	err = a.reloadAuxiliaryMPQs()
	if err != nil {
		return err
	}

	a.projectExplorer.SetProject(a.project)
	a.mpqExplorer.SetProject(a.project)

	a.CloseAllOpenWindows()

	if state, ok := a.config.ProjectStates[a.project.GetProjectFilePath()]; ok {
		a.RestoreAppState(state)
	} else {
		// if we don't have a state saved for this project, just open the project explorer
		a.projectExplorer.Show()
	}

	return nil
}

func (a *App) updateWindowTitle() {
	if a.project == nil {
		a.masterWindow.SetTitle(baseWindowTitle)
		return
	}

	a.masterWindow.SetTitle(baseWindowTitle + " - " + a.project.ProjectName)
}

func (a *App) toggleMPQExplorer() {
	a.mpqExplorer.ToggleVisibility()
}

func (a *App) onProjectPropertiesChanged(project *hsproject.Project) {
	a.project = project
	if err := a.project.Save(); err != nil {
		logErr("could not save project properties after changing, %s", err)
	}

	a.mpqExplorer.SetProject(a.project)
	a.updateWindowTitle()

	if err := a.reloadAuxiliaryMPQs(); err != nil {
		logErr("could not reload aux mpq's after changing project properties, %s", err)
	}
}

func (a *App) onPreferencesChanged(cfg *config.Config) {
	a.config = cfg
	if err := a.config.Save(); err != nil {
		logErr("after changing preferences, %s", err)
	}

	if a.project == nil {
		return
	}

	if err := a.reloadAuxiliaryMPQs(); err != nil {
		logErr("after changing preferences, %s", err)
	}
}

func (a *App) reloadAuxiliaryMPQs() error {
	if err := a.project.ReloadAuxiliaryMPQs(a.config); err != nil {
		return fmt.Errorf("could not reload aux mpq's in project, %w", err)
	}

	a.mpqExplorer.Reset()

	return nil
}

func (a *App) toggleProjectExplorer() {
	a.projectExplorer.ToggleVisibility()
}

func (a *App) closeActiveEditor() {
	for _, editor := range a.editors {
		if editor.HasFocus() {
			// don't call Cleanup here. the Render loop will call Cleanup when it notices that this editor isn't visible
			editor.SetVisible(false)
			return
		}
	}
}

func (a *App) closePopups() {
	a.projectPropertiesDialog.Cleanup()
	a.aboutDialog.Cleanup()
	a.preferencesDialog.Cleanup()
}

func (a *App) toggleConsole() {
	a.console.ToggleVisibility()
}

// CloseAllOpenWindows closes all opened windows
func (a *App) CloseAllOpenWindows() {
	a.closePopups()
	a.projectExplorer.Cleanup()
	a.mpqExplorer.Cleanup()
	a.focusedEditor = nil

	for _, editor := range a.editors {
		editor.Cleanup()
	}
}

// Save saves app state
func (a *App) Save() {
	if a.project != nil {
		a.config.ProjectStates[a.project.GetProjectFilePath()] = a.State()
	}

	if err := a.config.Save(); err != nil {
		logErr("failed to save config: %s", err)
		return
	}

	if a.focusedEditor != nil {
		a.focusedEditor.Save()
	}
}

// Quit quits the app
func (a *App) Quit() {
	if a.abyssWrapper.IsRunning() {
		_ = a.abyssWrapper.Kill()
	}

	a.Save()

	a.CloseAllOpenWindows()
}
