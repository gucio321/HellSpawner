// Package hsaboutdialog contains about dialog's data
package hsaboutdialog

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gucio321/giu"
	g "github.com/gucio321/giu"
	"github.com/jaytaylor/html2text"
	"github.com/russross/blackfriday"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

const (
	mainWindowW, mainWindowH = 256, 256
	mainLayoutW, mainLayoutH = 500, -1
)

const (
	white = 0xffffffff
)

// AboutDialog represents about dialog
type AboutDialog struct {
	*hsdialog.Dialog
	titleFont,
	regularFont,
	fixedFont *giu.FontInfo
	credits string
	license string
	readme  string
	logo    *g.Texture
}

// Create creates a new AboutDialog
func Create(textureLoader hscommon.TextureLoader, regularFont, titleFont, fixedFont *giu.FontInfo) (*AboutDialog, error) {
	result := &AboutDialog{
		Dialog:      hsdialog.New("About HellSpawner"),
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}

	textureLoader.CreateTextureFromFile(hsassets.HellSpawnerLogo, func(t *g.Texture) {
		result.logo = t
	})

	var err error

	var data []byte

	if data, err = ioutil.ReadFile("LICENSE"); err != nil {
		data = nil
	}

	result.license = string(data)

	if data, err = ioutil.ReadFile("CONTRIBUTORS"); err != nil {
		data = nil
	}

	result.credits = string(data)

	if data, err = ioutil.ReadFile("README.md"); err != nil {
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
	text = strings.Join(hsutil.SplitIntoLinesWithMaxWidth(text, 70), "\n")
	result.readme = text

	return result, nil
}

// Build build an about dialog
func (a *AboutDialog) Build() {
	colorWhite := hsutil.Color(white)
	a.IsOpen(&a.Visible).Layout(
		g.Row(
			g.Image(a.logo).Size(mainWindowW, mainWindowH),
			g.Child("AboutHellSpawnerLayout").Size(mainLayoutW, mainLayoutH).Layout(
				g.Style().SetColor(g.StyleColorText, colorWhite).To(
					g.Label("HellSpawner").Font(a.titleFont),
					g.Label("The OpenDiablo 2 Toolset").Font(a.regularFont),
					g.Label("Local Build").Font(a.fixedFont),
				),
				g.Separator(),
				g.TabBar("AboutHellSpawnerTabBar").Flags(g.TabBarFlagsNoCloseWithMiddleMouseButton).Layout(
					g.TabItem("README##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline("##AboutHellSpawnerReadme", &a.readme).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
					g.TabItem("Credits##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline("##AboutHellSpawnerCredits", &a.credits).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
					g.TabItem("Licenses##AboutHellSpawner").Layout(
						g.Style().SetFont(a.fixedFont).To(
							g.InputTextMultiline("##AboutHellSpawnerLicense", &a.license).
								Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						),
					),
				),
			),
		),
	).Build()
}
