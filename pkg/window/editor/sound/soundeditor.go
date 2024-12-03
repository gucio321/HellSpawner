// Package sound represents a soundEditor's window
package sound

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/gucio321/HellSpawner/pkg/common"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"github.com/gucio321/HellSpawner/pkg/widgets"
	"github.com/gucio321/HellSpawner/pkg/window/editor"

	g "github.com/AllenDang/giu"
)

const (
	mainWindowW, mainWindowH  = 300, 70
	progressIndicatorModifier = 60
	progressTimeModifier      = 22050
	btnSize                   = 20
)

// static check, to ensure, if sound editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a sound editor
type Editor struct {
	*editor.EditorBase

	streamer beep.StreamSeekCloser
	control  *beep.Ctrl
	format   beep.Format
	file     string
}

// Create creates a new sound editor
func Create(_ *config.Config,
	pathEntry *common.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	streamer, format, err := wav.Decode(bytes.NewReader(*data))
	if err != nil {
		return nil, fmt.Errorf("wav decode error: %w", err)
	}

	control := &beep.Ctrl{
		Streamer: beep.Loop(-1, streamer),
		Paused:   false,
	}

	result := &Editor{
		EditorBase: editor.New(pathEntry, x, y, project),
		file:       filepath.Base(pathEntry.FullPath),
		streamer:   streamer,
		control:    control,
		format:     format,
	}

	result.Path = pathEntry

	speaker.Play(result.control)

	return result, nil
}

// Build builds a sound editor
func (s *Editor) Build() {
	s.IsOpen(&s.Visible).
		Flags(g.WindowFlagsNoResize).
		Size(mainWindowW, mainWindowH).
		Layout(s.GetLayout())
}

func (s *Editor) GetLayout() g.Widget {
	isPlaying := !s.control.Paused

	secondsCurrent := s.streamer.Position() / progressTimeModifier
	secondsTotal := s.streamer.Len() / progressTimeModifier

	const progressBarHeight = 24 // px

	progress := float32(s.streamer.Position()) / float32(s.streamer.Len())

	return g.Row(
		widgets.PlayPauseButton(&isPlaying).
			OnPlayClicked(s.play).OnPauseClicked(s.stop).Size(btnSize, btnSize),
		g.ProgressBar(progress).Size(-1, progressBarHeight).
			Overlay(fmt.Sprintf("%d:%02d / %d:%02d",
				secondsCurrent/progressIndicatorModifier,
				secondsCurrent%progressIndicatorModifier,
				secondsTotal/progressIndicatorModifier,
				secondsTotal%progressIndicatorModifier,
			)),
	)
}

// Cleanup closes an editor
func (s *Editor) Cleanup() {
	speaker.Lock()
	s.control.Paused = true

	if err := s.streamer.Close(); err != nil {
		log.Print(err)
	}

	if s.HasChanges(s) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			s.Path.FullPath).YesNo(); shouldSave {
			s.Save()
		}
	}

	s.EditorBase.Cleanup()
	speaker.Unlock()
}

func (s *Editor) play() {
	speaker.Lock()
	s.control.Paused = false
	speaker.Unlock()
}

func (s *Editor) stop() {
	speaker.Lock()

	if s.control.Paused {
		if err := s.streamer.Seek(0); err != nil {
			log.Print(err)
			return
		}
	}

	s.control.Paused = true

	speaker.Unlock()
}

// UpdateMainMenuLayout updates mainMenu's layout to it contain soundEditor's options
func (s *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Sound Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			s.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// GenerateSaveData generates data to be saved
func (s *Editor) GenerateSaveData() []byte {
	// https://github.com/gucio321/HellSpawner/issues/181
	data, _ := s.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (s *Editor) Save() {
	s.EditorBase.Save(s)
}
