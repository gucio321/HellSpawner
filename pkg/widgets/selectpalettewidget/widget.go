package selectpalettewidget

import (
	"github.com/gucio321/HellSpawner/pkg/app/config"
	"github.com/gucio321/HellSpawner/pkg/window/popup"
	"log"
	"path/filepath"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsfiletypes"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/mpqexplorer"
	"github.com/gucio321/HellSpawner/pkg/window/toolwindow/projectexplorer"
)

const (
	paletteSelectW, paletteSelectH = 400, 600
	actionButtonW, actionButtonH   = 200, 30
)

// SelectPaletteWidget represents an pop-up MPQ explorer, when we're
// selectin DAT palette
type SelectPaletteWidget struct {
	mpqExplorer     *mpqexplorer.MPQExplorer
	projectExplorer *projectexplorer.ProjectExplorer
	id              string
	saveCB          func(colors *[256]d2interface.Color)
	closeCB         func()
}

// NewSelectPaletteWidget creates a select palette widget
func NewSelectPaletteWidget(
	id string,
	project *hsproject.Project,
	cfg *config.Config,
	saveCB func(colors *[256]d2interface.Color),
	closeCB func(),
) *SelectPaletteWidget {
	result := &SelectPaletteWidget{
		id:      id,
		saveCB:  saveCB,
		closeCB: closeCB,
	}

	callback := func(path *common.PathEntry) {
		bytes, bytesErr := path.GetFileBytes()
		if bytesErr != nil {
			log.Print(bytesErr)

			return
		}

		ft, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &bytes)
		if err != nil {
			log.Print(err)

			return
		}

		if ft == hsfiletypes.FileTypePalette {
			// load new palette:
			paletteData, err := path.GetFileBytes()
			if err != nil {
				log.Print(err)
			}

			palette, err := d2dat.Load(paletteData)
			if err != nil {
				log.Print(err)
			}

			colors := palette.GetColors()

			saveCB(&colors)
			closeCB()
		}
	}

	mpqExplorer, err := mpqexplorer.Create(callback, cfg, 0, 0)
	if err != nil {
		log.Print(err)
	}

	mpqExplorer.SetProject(project)

	result.mpqExplorer = mpqExplorer

	projectExplorer, err := projectexplorer.Create(callback, 0, 0)
	if err != nil {
		log.Print(err)
	}

	projectExplorer.SetProject(project)

	result.projectExplorer = projectExplorer

	return result
}

// Build builds a widget
func (p *SelectPaletteWidget) Build() {
	// always true (we don't use this feature in this case
	isOpen := true
	giu.Layout{
		popup.New("##" + p.id + "popUpSelectPalette").IsOpen(&isOpen).Layout(giu.Layout{
			giu.Child().Size(paletteSelectW, paletteSelectH).Layout(giu.Layout{
				p.projectExplorer.GetProjectTreeNodes(),
				giu.Layout(p.mpqExplorer.GetMpqTreeNodes()),
				giu.Separator(),
				giu.Button("Don't use any palette##"+p.id+"selectPaletteDonotUseAny").
					Size(actionButtonW, actionButtonH).
					OnClick(func() {
						p.saveCB(nil)
						p.closeCB()
					}),
				giu.Button("Exit##"+p.id+"selectPaletteExit").
					Size(actionButtonW, actionButtonH).
					OnClick(func() {
						p.closeCB()
					}),
			}),
		}),
	}.Build()

	if !isOpen {
		p.closeCB()
	}
}
