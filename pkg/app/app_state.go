package app

import (
	"encoding/json"
	"log"

	state "github.com/gucio321/HellSpawner/pkg/app/state"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow"
)

// State creates a new app state
func (a *App) State() state.AppState {
	appState := state.AppState{
		ProjectPath:   a.project.GetProjectFilePath(),
		EditorWindows: []state.EditorState{},
		ToolWindows:   []state.ToolWindowState{},
	}

	for _, editor := range a.editors {
		appState.EditorWindows = append(appState.EditorWindows, editor.State())
	}

	appState.ToolWindows = append(
		appState.ToolWindows,
		a.mpqExplorer.State(),
		a.projectExplorer.State(),
		a.console.State(),
	)

	return appState
}

// RestoreAppState restores an app state
func (a *App) RestoreAppState(appState state.AppState) {
	for _, toolState := range appState.ToolWindows {
		var tool toolwindow.ToolWindow

		switch toolState.Type {
		case state.ToolWindowTypeConsole:
			tool = a.console
		case state.ToolWindowTypeMPQExplorer:
			tool = a.mpqExplorer
		case state.ToolWindowTypeProjectExplorer:
			tool = a.projectExplorer
		default:
			continue
		}

		tool.Pos(toolState.PosX, toolState.PosY)
		tool.SetVisible(toolState.Visible)
		tool.Size(toolState.Width, toolState.Height)
	}

	for _, editorState := range appState.EditorWindows {
		var path common.PathEntry

		err := json.Unmarshal(editorState.Path, &path)
		if err != nil {
			log.Print("failed to restore editor: ", err)
			continue
		}

		go a.createEditor(&path, editorState.Encoded, editorState.PosX, editorState.PosY, editorState.Width, editorState.Height)
	}
}
