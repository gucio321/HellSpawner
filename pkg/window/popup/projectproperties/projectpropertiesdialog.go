// Package projectproperties contains project properties dialog's data
package projectproperties

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gucio321/HellSpawner/pkg/app/assets"
	"github.com/gucio321/HellSpawner/pkg/app/config"
	"github.com/gucio321/HellSpawner/pkg/window"

	"github.com/gucio321/HellSpawner/pkg/common"

	g "github.com/AllenDang/giu"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/window/popup"
)

const (
	mainWindowW, mainWindowH   = 300, 200
	mpqSelectW, mpqSelectH     = 300, 250
	mpqGroupW, mpqGroupH       = 0, 180
	imgBtnW, imgBtnH           = 16, 16
	dummyW, dummyH             = 8, 0
	inputTextSize              = 250
	descriptionW, descriptionH = inputTextSize, 100
)

var _ window.Renderable = &Dialog{}

// Dialog represent project properties' dialog
type Dialog struct {
	*popup.Dialog

	removeIconTexture          *g.Texture
	upIconTexture              *g.Texture
	downIconTexture            *g.Texture
	project                    hsproject.Project
	config                     *config.Config
	onProjectPropertiesChanged func(project *hsproject.Project)
	auxMPQs, auxMPQNames       []string
	mpqsToAdd                  []int

	mpqSelectDialogVisible bool
}

// Create creates a new project properties' dialog
func Create(onProjectPropertiesChanged func(project *hsproject.Project)) *Dialog {
	result := &Dialog{
		Dialog:                     popup.New("Project Properties"),
		onProjectPropertiesChanged: onProjectPropertiesChanged,
		mpqSelectDialogVisible:     false,
	}

	common.LoadTexture(assets.DeleteIcon, func(texture *g.Texture) {
		result.removeIconTexture = texture
	})

	common.LoadTexture(assets.UpArrowIcon, func(texture *g.Texture) {
		result.upIconTexture = texture
	})

	common.LoadTexture(assets.DownArrowIcon, func(texture *g.Texture) {
		result.downIconTexture = texture
	})

	return result
}

// Show shows project properties dialog
func (p *Dialog) Show(project *hsproject.Project, cfg *config.Config) {
	p.config = cfg
	p.project = *project
	p.auxMPQs = cfg.GetAuxMPQs()
	p.auxMPQNames = make([]string, len(p.auxMPQs))

	for idx := range p.auxMPQNames {
		p.auxMPQNames[idx] = filepath.Base(p.auxMPQs[idx])
	}

	p.mpqsToAdd = make([]int, 0)

	p.Dialog.Show()
}

// Build builds a dialog
//
//nolint:gocognit,funlen,gocyclo // no need to change
func (p *Dialog) Build() {
	p.IsOpen(&p.mpqSelectDialogVisible).Layout(p.GetLayout()).Build()
}

func (p *Dialog) GetLayout() g.Widget {
	canSave := strings.TrimSpace(p.project.ProjectName) != ""

	if !p.mpqSelectDialogVisible {
		p.IsOpen(&p.Visible).Layout(
			g.Row(
				g.Child().Size(mpqSelectW, mpqSelectH).Layout(
					g.Label("Project Name:"),
					g.InputText(&p.project.ProjectName).Size(inputTextSize),
					g.Label("Description:"),
					g.InputTextMultiline(&p.project.Description).Size(descriptionW, descriptionH),
					g.Label("Author:"),
					g.InputText(&p.project.Author).Size(inputTextSize),
				),
				g.Child().Size(mpqSelectW, mpqSelectH).Layout(
					g.Label("Auxiliary MPQs:"),
					g.Child().Border(false).Size(mpqGroupW, mpqGroupH).Layout(
						g.Custom(func() {
							imgui.PushStyleColorVec4(imgui.ColButton, imgui.Vec4{})
							imgui.PushStyleColorVec4(imgui.ColBorder, imgui.Vec4{})
							imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{})

							for idx := range p.project.AuxiliaryMPQs {
								currentIdx := idx

								if idx >= len(p.project.AuxiliaryMPQs) {
									break
								}

								g.Row(
									g.ImageButton(p.removeIconTexture).Size(imgBtnW, imgBtnH).OnClick(func() {
										copy(p.project.AuxiliaryMPQs[currentIdx:], p.project.AuxiliaryMPQs[currentIdx+1:])
										p.project.AuxiliaryMPQs = p.project.AuxiliaryMPQs[:len(p.project.AuxiliaryMPQs)-1]
									}),
									g.ImageButton(p.downIconTexture).Size(imgBtnW, imgBtnH).OnClick(func() {
										if currentIdx < len(p.project.AuxiliaryMPQs)-1 {
											p.project.AuxiliaryMPQs[currentIdx],
												p.project.AuxiliaryMPQs[currentIdx+1] = p.project.AuxiliaryMPQs[currentIdx+1],
												p.project.AuxiliaryMPQs[currentIdx]
										}
									}),
									g.ImageButton(p.upIconTexture).Size(imgBtnW, imgBtnH).OnClick(func() {
										if currentIdx > 0 {
											p.project.AuxiliaryMPQs[currentIdx-1],
												p.project.AuxiliaryMPQs[currentIdx] = p.project.AuxiliaryMPQs[currentIdx],
												p.project.AuxiliaryMPQs[currentIdx-1]
										}
									}),
									g.Dummy(dummyW, dummyH),
									g.Label(p.project.AuxiliaryMPQs[idx]),
								).Build()
							}

							//nolint:mnd // we've pushed twice and we need to pop twice
							imgui.PopStyleColorV(2)
							imgui.PopStyleVar()
						}),
					),
					g.Button("Add Auxiliary MPQ...##ProjectPropertiesAddAuxMpq").OnClick(p.onAddAuxMpqClicked),
				),
			),
			g.Row(
				g.Custom(func() {
					const halfOpacity = 0.5

					if !canSave {
						imgui.PushStyleVarFloat(imgui.StyleVarAlpha, halfOpacity)
					}
				}),
				g.Button("Save##DialogSave").OnClick(p.onSaveClicked),
				g.Custom(func() {
					if !canSave {
						imgui.PopStyleVar()
					}
				}),
				g.Button("Cancel##DialogCancel").OnClick(p.onCancelClicked),
			),
		)
	}

	return g.Layout{
		g.Child().Size(mainWindowW, mainWindowH).Layout(
			g.Custom(func() {
				addMPQ := func(i int) {
					p.mpqsToAdd = append(p.mpqsToAdd, i)
				}
				removeMPQ := func(i int) {
					for n, idx := range p.mpqsToAdd {
						if i == idx {
							p.mpqsToAdd = append(p.mpqsToAdd[:n], p.mpqsToAdd[n+1:]...)
						}
					}
				}

				isInMpqList := func(i int) bool {
					for _, idx := range p.mpqsToAdd {
						if i == idx {
							return true
						}
					}

					return false
				}

				const listItemHeight = 20

				// list of `Selectable widgets`;
				for i, mpq := range p.auxMPQNames {
					isSelected := isInMpqList(i)
					g.Row(
						g.Checkbox(
							"##"+"ProjectPropertiesSelectAuxMPQDialogCheckbox"+strconv.Itoa(i),
							&isSelected,
						).OnChange(func() {
							// opposite, because giu.Checkbox already changed this value
							if !isSelected {
								removeMPQ(i)
							} else {
								addMPQ(i)
							}
						}),
						g.Selectable(mpq+"##"+"ProjectPropertiesSelectAuxMPQDialogIdx"+strconv.Itoa(i)).
							Selected(isSelected).
							Size(mainWindowW, listItemHeight).OnClick(func() {
							if isSelected {
								removeMPQ(i)
							} else {
								addMPQ(i)
							}
						}),
					).Build()
				}
			}),
		),
		g.Row(
			g.Button("Add Selected...##ProjectPropertiesSelectAuxMPQDialogAddSelected").OnClick(func() {
				// checks if aux MPQs list isn't empty
				if len(p.auxMPQs) > 0 {
					for _, idx := range p.mpqsToAdd {
						p.addAuxMpq(p.auxMPQs[idx])
					}

					p.onProjectPropertiesChanged(&p.project)
				}

				p.mpqSelectDialogVisible = false
			}),
			g.Button("Cancel##ProjectPropertiesSelectAuxMPQDialogCancel").OnClick(func() {
				p.mpqSelectDialogVisible = false
			}),
		),
	}
}

func (p *Dialog) onSaveClicked() {
	if strings.TrimSpace(p.project.ProjectName) == "" {
		return
	}

	p.onProjectPropertiesChanged(&p.project)
	p.Visible = false
}

func (p *Dialog) onCancelClicked() {
	p.Visible = false
}

func (p *Dialog) onAddAuxMpqClicked() {
	p.mpqSelectDialogVisible = true
}

func (p *Dialog) addAuxMpq(mpqPath string) {
	relPath, err := filepath.Rel(p.config.AuxiliaryMpqPath, mpqPath)
	if err != nil {
		log.Print(err)
		return
	}

	for idx := range p.project.AuxiliaryMPQs {
		if p.project.AuxiliaryMPQs[idx] == relPath {
			return
		}
	}

	p.project.AuxiliaryMPQs = append(p.project.AuxiliaryMPQs, relPath)
}
