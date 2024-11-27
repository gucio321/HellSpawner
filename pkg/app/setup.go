package app

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"

	"github.com/OpenDiablo2/dialog"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	g "github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/window/editor/ds1"
	"github.com/gucio321/HellSpawner/pkg/window/editor/dt1"
	"github.com/gucio321/HellSpawner/pkg/window/editor/fonttable"
	"github.com/gucio321/HellSpawner/pkg/window/editor/palettemap"
	"github.com/gucio321/HellSpawner/pkg/window/editor/stringtable"

	"github.com/gucio321/HellSpawner/pkg/assets"
	"github.com/gucio321/HellSpawner/pkg/common/hsenum"
	"github.com/gucio321/HellSpawner/pkg/common/hsfiletypes"
	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/window/editor/animdata"
	"github.com/gucio321/HellSpawner/pkg/window/editor/cof"
	"github.com/gucio321/HellSpawner/pkg/window/editor/dc6"
	"github.com/gucio321/HellSpawner/pkg/window/editor/dcc"
	"github.com/gucio321/HellSpawner/pkg/window/editor/font"
	"github.com/gucio321/HellSpawner/pkg/window/editor/palette"
	"github.com/gucio321/HellSpawner/pkg/window/editor/sound"
	"github.com/gucio321/HellSpawner/pkg/window/editor/text"
	"github.com/gucio321/HellSpawner/pkg/window/popup/aboutdialog"
	"github.com/gucio321/HellSpawner/pkg/window/popup/preferences"
	"github.com/gucio321/HellSpawner/pkg/window/popup/projectproperties"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/hsconsole"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/mpqexplorer"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/hsprojectexplorer"
)

func (a *App) setup() (err error) {
	dialog.Init()

	a.setupMasterWindow()
	a.setupConsole()
	a.setupAutoSave()
	a.registerGlobalKeyboardShortcuts()
	a.registerEditors()

	err = a.setupAudio()
	if err != nil {
		return err
	}

	err = a.setupMainMpqExplorer()
	if err != nil {
		return err
	}

	err = a.setupProjectExplorer()
	if err != nil {
		return err
	}

	err = a.setupDialogs()
	if err != nil {
		return err
	}

	// we may have tried loading some textures already...
	a.TextureLoader.ProcessTextureLoadRequests()

	return nil
}

func (a *App) setupMasterWindow() {
	a.masterWindow = g.NewMasterWindow(baseWindowTitle, baseWindowW, baseWindowH, 0)
	a.setupFonts()

	bgColor := a.determineBackgroundColor()
	a.masterWindow.SetBgColor(bgColor)
}

func (a *App) determineBackgroundColor() color.RGBA {
	const bitSize = 32

	result := a.config.BGColor

	strBytes := []byte(*a.Flags.bgColor)
	numChars := len(strBytes)
	includesBase := strBytes[1] == 'x'

	base := 16
	if includesBase {
		base = 0
	}

	includesAlpha := false
	if includesBase && numChars >= len("0xRGGBBAA") {
		includesAlpha = true
	} else if !includesBase && numChars >= len("RGGBBAA") {
		includesAlpha = true
	}

	bg, err := strconv.ParseInt(*a.Flags.bgColor, base, bitSize)
	if err == nil {
		if !includesAlpha {
			bg <<= 8
		}

		//nolint:gosec // this complains about intager size, but this number was created by strconv.ParseInt with bitsize 32
		result = hsutil.Color(uint32(bg))
	}

	return result
}

func (a *App) setupAutoSave() {
	go func() {
		time.Sleep(autoSaveTimer * time.Second)
		a.Save()
	}()
}

func (a *App) registerEditors() {
	a.editorConstructors[hsfiletypes.FileTypeText] = text.Create
	a.editorConstructors[hsfiletypes.FileTypeAudio] = sound.Create
	a.editorConstructors[hsfiletypes.FileTypePalette] = palette.Create
	a.editorConstructors[hsfiletypes.FileTypeAnimationData] = animdata.Create
	a.editorConstructors[hsfiletypes.FileTypeDC6] = dc6.Create
	a.editorConstructors[hsfiletypes.FileTypeDCC] = dcc.Create
	a.editorConstructors[hsfiletypes.FileTypeCOF] = cof.Create
	a.editorConstructors[hsfiletypes.FileTypeFont] = font.Create
	a.editorConstructors[hsfiletypes.FileTypeDT1] = dt1.Create
	a.editorConstructors[hsfiletypes.FileTypePL2] = palettemap.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLStringTable] = stringtable.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLFontTable] = fonttable.Create
	a.editorConstructors[hsfiletypes.FileTypeDS1] = ds1.Create
}

func (a *App) setupMainMpqExplorer() error {
	window, err := mpqexplorer.Create(a.openEditor, a.config, mpqExplorerDefaultX, mpqExplorerDefaultY)
	if err != nil {
		return fmt.Errorf("error creating a MPQ explorer: %w", err)
	}

	a.mpqExplorer = window

	return nil
}

func (a *App) setupProjectExplorer() error {
	x, y := float32(projectExplorerDefaultX), float32(projectExplorerDefaultY)

	window, err := hsprojectexplorer.Create(a.TextureLoader,
		a.openEditor, x, y)
	if err != nil {
		return fmt.Errorf("error creating a project explorer: %w", err)
	}

	a.projectExplorer = window

	return nil
}

func (a *App) setupAudio() error {
	sampleRate := beep.SampleRate(samplesPerSecond)
	bufferSize := sampleRate.N(sampleDuration)

	if err := speaker.Init(sampleRate, bufferSize); err != nil {
		return fmt.Errorf("could not initialize, %w", err)
	}

	return nil
}

func (a *App) setupConsole() {
	a.console = hsconsole.Create(a.fontFixed, consoleDefaultX, consoleDefaultY, a.logFile)

	log.SetFlags(log.Lshortfile)
	log.SetOutput(a.console)

	t := time.Now()
	y, m, d := t.Date()

	line := fmt.Sprintf("%d-%d-%d, %d:%d:%d", y, m, d, t.Hour(), t.Minute(), t.Second())
	log.Printf(logFileSeparator, line)
}

func (a *App) setupDialogs() error {
	// Register the dialogs
	about, err := aboutdialog.Create(a.TextureLoader, a.diabloRegularFont, a.diabloBoldFont, a.fontFixedSmall)
	if err != nil {
		return fmt.Errorf("error creating an about dialog: %w", err)
	}

	a.aboutDialog = about
	a.projectPropertiesDialog = projectproperties.Create(a.TextureLoader, a.onProjectPropertiesChanged)
	a.preferencesDialog = preferences.Create(a.onPreferencesChanged, a.masterWindow.SetBgColor)

	return nil
}

// please note, that this steps will not affect app language
// it will only load an appropriate glyph ranges for
// displayed text (e.g. for string/font table editors)
func (a *App) setupFonts() {
	font := assets.FontNotoSansRegular

	switch a.config.Locale {
	// glyphs supported by default
	case hsenum.LocaleEnglish, hsenum.LocaleGerman,
		hsenum.LocaleFrench, hsenum.LocaleItalien,
		hsenum.LocaleSpanish, hsenum.LocalePolish:
		// noop
	case hsenum.LocaleChineseTraditional:
		font = assets.FontSourceHanSerif
	case hsenum.LocaleKorean:
		font = assets.FontSourceHanSerif
	}

	g.Context.FontAtlas.SetDefaultFontFromBytes(font, baseFontSize)

	// please note, that the following fonts will not use
	// previously generated glyph ranges.
	// they'll have a default range
	a.fontFixed = g.Context.FontAtlas.AddFontFromBytes("fixed font", assets.FontCascadiaCode, fixedFontSize)
	a.fontFixedSmall = g.Context.FontAtlas.AddFontFromBytes("small fixed font", assets.FontCascadiaCode, fixedSmallFontSize)
	a.diabloRegularFont = g.Context.FontAtlas.AddFontFromBytes("diablo regular", assets.FontDiabloRegular, diabloRegularFontSize)
	a.diabloBoldFont = g.Context.FontAtlas.AddFontFromBytes("diablo bold", assets.FontDiabloBold, diabloBoldFontSize)
}

func (a *App) registerGlobalKeyboardShortcuts() {
	a.masterWindow.RegisterKeyboardShortcuts(
		g.WindowShortcut{Key: g.KeyN, Modifier: g.ModControl + g.ModShift, Callback: a.onNewProjectClicked},
		g.WindowShortcut{Key: g.KeyO, Modifier: g.ModControl, Callback: a.onOpenProjectClicked},
		g.WindowShortcut{Key: g.KeyS, Modifier: g.ModControl, Callback: a.Save},
		g.WindowShortcut{Key: g.KeyP, Modifier: g.ModAlt, Callback: a.onFilePreferencesClicked},
		g.WindowShortcut{Key: g.KeyQ, Modifier: g.ModAlt, Callback: a.Quit},
		g.WindowShortcut{Key: g.KeyF1, Modifier: g.ModNone, Callback: a.onHelpAboutClicked},

		g.WindowShortcut{Key: g.KeyW, Modifier: g.ModControl, Callback: a.closeActiveEditor},
		g.WindowShortcut{Key: g.KeyEscape, Modifier: g.ModNone, Callback: func() { a.closePopups(); a.closeActiveEditor() }},

		g.WindowShortcut{Key: g.KeyM, Modifier: g.ModControl + g.ModShift, Callback: a.toggleMPQExplorer},
		g.WindowShortcut{Key: g.KeyP, Modifier: g.ModControl + g.ModShift, Callback: a.toggleProjectExplorer},
		g.WindowShortcut{Key: g.KeyC, Modifier: g.ModControl + g.ModShift, Callback: a.toggleConsole},
	)
}
