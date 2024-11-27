// Package preferences contains preferences dialog data
package preferences

import (
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/common/hsenum"
	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/config"
	"github.com/gucio321/HellSpawner/pkg/window/popup"
)

const (
	mainWindowW, mainWindowH = 320, 200
	textboxSize              = 245
	btnW, btnH               = 30, 0
)

// Dialog represents preferences dialog
type Dialog struct {
	*popup.Dialog

	config             *config.Config
	onConfigChanged    func(config *config.Config)
	windowColorChanger func(c color.Color)
	restartPrompt      bool
}

// Create creates a new preferences dialog
func Create(onConfigChanged func(config *config.Config), windowColorChanger func(c color.Color)) *Dialog {
	result := &Dialog{
		Dialog:             popup.New("Preferences"),
		onConfigChanged:    onConfigChanged,
		windowColorChanger: windowColorChanger,
		restartPrompt:      false,
	}
	result.Visible = false

	return result
}

// Build builds a preferences dialog
func (p *Dialog) Build() {
	locales := make([]string, 0)
	for i := hsenum.LocaleEnglish; i <= hsenum.LocalePolish; i++ {
		locales = append(locales, i.String())
	}

	locale := int32(p.config.Locale)

	p.IsOpen(&p.Visible).Layout(
		g.Child().Size(mainWindowW, mainWindowH).Layout(
			g.Label("Auxiliary MPQ Path"),
			g.Row(
				g.InputText(&p.config.AuxiliaryMpqPath).Size(textboxSize).Flags(g.InputTextFlagsReadOnly),
				g.Button("...##AppPreferencesAuxMPQPathBrowse").Size(btnW, btnH).OnClick(p.onBrowseAuxMpqPathClicked),
			),
			g.Separator(),
			g.Label("External MPQ listfile Path"),
			g.Row(
				g.InputText(&p.config.ExternalListFile).Size(textboxSize).Flags(g.InputTextFlagsReadOnly),
				g.Button("...##AppPreferencesListfilePathBrowse").Size(btnW, btnH).OnClick(p.onBrowseExternalListfileClicked),
			),
			g.Separator(),
			g.Label("Abyss Engine Path"),
			g.Row(
				g.InputText(&p.config.AbyssEnginePath).Size(textboxSize).Flags(g.InputTextFlagsReadOnly),
				g.Button("...##AppPreferencesAbyssEnginePathBrowse").Size(btnW, btnH).OnClick(p.onBrowseAbyssEngineClicked),
			),
			g.Separator(),
			g.Checkbox("Open most recent project on start-up", &p.config.OpenMostRecentOnStartup),
			g.Separator(),
			g.Checkbox("Save log output in a log file", &p.config.LoggingToFile),
			g.Custom(func() {
				if p.config.LoggingToFile {
					g.Layout{
						g.Label("Log file path"),
						g.Row(
							g.InputText(&p.config.LogFilePath).Size(textboxSize).Flags(g.InputTextFlagsReadOnly),
							g.Button("...##AppPreferencesLogFilePathBrowse").Size(btnW, btnH).OnClick(p.onBrowseLogFilePathClicked),
						),
					}.Build()
				}
			}),
			g.Separator(),
			g.Custom(func() {
				if !p.restartPrompt {
					return
				}

				g.Layout{
					g.Label("WARNING: to introduce there changes"),
					g.Label("you need to restart HellSpawner"),
				}.Build()
			}),
			g.Combo("locale", locales[locale], locales, &locale).OnChange(func() {
				p.restartPrompt = true
				p.config.Locale = hsenum.Locale(locale)
			}),
			g.Separator(),
			g.Label("Background color:"),
			g.Row(
				g.ColorEdit("##BackgroundColor", &p.config.BGColor).
					Flags(g.ColorEditFlagsNoAlpha).OnChange(func() {
					p.windowColorChanger(p.config.BGColor)
				}),
				g.Button("Default##BackgroundColorDefault").OnClick(func() {
					p.config.BGColor = hsutil.Color(config.DefaultBGColor)
					p.windowColorChanger(p.config.BGColor)
				}),
			),
		),
		g.Row(
			g.Button("Save##AppPreferencesSave").OnClick(p.onSaveClicked),
			g.Button("Cancel##AppPreferencesCancel").OnClick(p.onCancelClicked),
		),
	).Build()
}

// Show switch preferences dialog to visible state
func (p *Dialog) Show(cfg *config.Config) {
	p.Dialog.Show()

	p.config = cfg
}

func (p *Dialog) onBrowseAuxMpqPathClicked() {
	path, err := dialog.Directory().Browse()
	if err != nil || path == "" {
		return
	}

	p.config.AuxiliaryMpqPath = path
}

func (p *Dialog) onBrowseExternalListfileClicked() {
	path := dialog.File()
	path.Filter("Text file", "txt")
	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	p.config.ExternalListFile = filePath
}

func (p *Dialog) onBrowseLogFilePathClicked() {
	path, err := dialog.Directory().Browse()
	if err != nil || path == "" {
		return
	}

	p.config.LogFilePath = path
}

func (p *Dialog) onSaveClicked() {
	p.onConfigChanged(p.config)
	p.Visible = false
}

func (p *Dialog) onCancelClicked() {
	p.Visible = false
}

func (p *Dialog) onBrowseAbyssEngineClicked() {
	path := dialog.File()

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	p.config.AbyssEnginePath = filePath
}
