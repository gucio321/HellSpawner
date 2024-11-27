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

	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsds1editor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsdt1editor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsfonttableeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hspalettemapeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsstringtableeditor"

	"github.com/gucio321/HellSpawner/hsassets"
	"github.com/gucio321/HellSpawner/hscommon/hsenum"
	"github.com/gucio321/HellSpawner/hscommon/hsfiletypes"
	"github.com/gucio321/HellSpawner/hscommon/hsutil"
	"github.com/gucio321/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/gucio321/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/gucio321/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsanimdataeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hscofeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsdc6editor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsdcceditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hsfonteditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hspaletteeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hssoundeditor"
	"github.com/gucio321/HellSpawner/hswindow/hseditor/hstexteditor"
	"github.com/gucio321/HellSpawner/hswindow/hstoolwindow/hsconsole"
	"github.com/gucio321/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/gucio321/HellSpawner/hswindow/hstoolwindow/hsprojectexplorer"
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
	const bitSize = 64

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
	a.editorConstructors[hsfiletypes.FileTypeText] = hstexteditor.Create
	a.editorConstructors[hsfiletypes.FileTypeAudio] = hssoundeditor.Create
	a.editorConstructors[hsfiletypes.FileTypePalette] = hspaletteeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeAnimationData] = hsanimdataeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDC6] = hsdc6editor.Create
	a.editorConstructors[hsfiletypes.FileTypeDCC] = hsdcceditor.Create
	a.editorConstructors[hsfiletypes.FileTypeCOF] = hscofeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeFont] = hsfonteditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDT1] = hsdt1editor.Create
	a.editorConstructors[hsfiletypes.FileTypePL2] = hspalettemapeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLStringTable] = hsstringtableeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLFontTable] = hsfonttableeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDS1] = hsds1editor.Create
}

func (a *App) setupMainMpqExplorer() error {
	window, err := hsmpqexplorer.Create(a.openEditor, a.config, mpqExplorerDefaultX, mpqExplorerDefaultY)
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
	about, err := hsaboutdialog.Create(a.TextureLoader, a.diabloRegularFont, a.diabloBoldFont, a.fontFixedSmall)
	if err != nil {
		return fmt.Errorf("error creating an about dialog: %w", err)
	}

	a.aboutDialog = about
	a.projectPropertiesDialog = hsprojectpropertiesdialog.Create(a.TextureLoader, a.onProjectPropertiesChanged)
	a.preferencesDialog = hspreferencesdialog.Create(a.onPreferencesChanged, a.masterWindow.SetBgColor)

	return nil
}

// please note, that this steps will not affect app language
// it will only load an appropriate glyph ranges for
// displayed text (e.g. for string/font table editors)
func (a *App) setupFonts() {
	font := hsassets.FontNotoSansRegular

	switch a.config.Locale {
	// glyphs supported by default
	case hsenum.LocaleEnglish, hsenum.LocaleGerman,
		hsenum.LocaleFrench, hsenum.LocaleItalien,
		hsenum.LocaleSpanish, hsenum.LocalePolish:
		// noop
	case hsenum.LocaleChineseTraditional:
		font = hsassets.FontSourceHanSerif
	case hsenum.LocaleKorean:
		font = hsassets.FontSourceHanSerif
	}

	g.Context.FontAtlas.SetDefaultFontFromBytes(font, baseFontSize)

	// please note, that the following fonts will not use
	// previously generated glyph ranges.
	// they'll have a default range
	a.fontFixed = g.Context.FontAtlas.AddFontFromBytes("fixed font", hsassets.FontCascadiaCode, fixedFontSize)
	a.fontFixedSmall = g.Context.FontAtlas.AddFontFromBytes("small fixed font", hsassets.FontCascadiaCode, fixedSmallFontSize)
	a.diabloRegularFont = g.Context.FontAtlas.AddFontFromBytes("diablo regular", hsassets.FontDiabloRegular, diabloRegularFontSize)
	a.diabloBoldFont = g.Context.FontAtlas.AddFontFromBytes("diablo bold", hsassets.FontDiabloBold, diabloBoldFontSize)
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