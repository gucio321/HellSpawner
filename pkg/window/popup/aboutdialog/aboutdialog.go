// Package aboutdialog provides the "About" window implementation, which shows information about hellspawner.
package aboutdialog

import (
	"fmt"
	"os"
	"strings"

	"github.com/gucio321/HellSpawner/pkg/app/assets"
	"github.com/gucio321/HellSpawner/pkg/window"

	g "github.com/AllenDang/giu"
	"github.com/jaytaylor/html2text"
	"github.com/russross/blackfriday"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/window/popup"
)

const (
	mainWindowW, mainWindowH = 256, 256
	mainLayoutW, mainLayoutH = 500, -1
)

const (
	white = 0xffffffff
)

var _ window.Renderable = &AboutDialog{}

// AboutDialog represents about dialog
type AboutDialog struct {
	*popup.Dialog
	titleFont   *g.FontInfo
	regularFont *g.FontInfo
	fixedFont   *g.FontInfo
	credits     string
	license     string
	readme      string
	logo        *g.Texture
}

// Create creates a new AboutDialog
func Create(regularFont, titleFont, fixedFont *g.FontInfo) (*AboutDialog, error) {
	result := &AboutDialog{
		Dialog:      popup.New("About HellSpawner"),
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}

	common.LoadTexture(assets.HellSpawnerLogo, func(t *g.Texture) {
		result.logo = t
	})

	var err error

	var data []byte

	if data, err = os.ReadFile("LICENSE"); err != nil {
		data = nil
	}

	result.license = string(data)

	if data, err = os.ReadFile("CONTRIBUTORS"); err != nil {
		data = nil
	}

	result.credits = string(data)

	if data, err = os.ReadFile("README.md"); err != nil {
		data = nil
	}

	// convert output md to html
	html := blackfriday.MarkdownBasic(data)
	// convert html to text
	text, err := html2text.FromString(string(html), html2text.Options{PrettyTables: true})
	if err != nil {
		return result, fmt.Errorf("error converting HTML to text: %w", err)
	}

	// set string's max length
	const maxColumns = 70
	text = strings.Join(hsutil.SplitIntoLinesWithMaxWidth(text, maxColumns), "\n")
	result.readme = text

	return result, nil
}

// Build build an about dialog
func (a *AboutDialog) Build() {
	a.IsOpen(&a.Visible).Layout(a.GetLayout()).Build()
}

func (a *AboutDialog) GetLayout() g.Widget {
	colorWhite := hsutil.Color(white)

	return g.Layout{
		g.Row(
			g.Image(a.logo).Size(mainWindowW, mainWindowH),
			g.Child().Size(mainLayoutW, mainLayoutH).Layout(
				g.Style().SetColor(g.StyleColorText, colorWhite).To(
					g.Label("HellSpawner").Font(a.titleFont),
					g.Label("The OpenDiablo 2 Toolset").Font(a.regularFont),
					g.Label("Local Build").Font(a.fixedFont),
				),
				g.Separator(),
				g.TabBar().Flags(g.TabBarFlagsNoCloseWithMiddleMouseButton).TabItems(
					g.TabItem("README##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline(&a.readme).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
					g.TabItem("Credits##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline(&a.credits).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
					g.TabItem("Licenses##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline(&a.license).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
				),
			),
		),
	}
}
