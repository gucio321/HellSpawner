package ds1widget

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math/d2vector"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2path"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/gucio321/HellSpawner/pkg/widgets"
)

const (
	layerDeleteButtonSize                = 24
	inputIntW                            = 40
	filePathW                            = 200
	deleteButtonSize                     = 15
	actionButtonW, actionButtonH         = 170, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	imageW, imageH                       = 32, 32
)

type widget struct {
	id                  giu.ID
	ds1                 *d2ds1.DS1
	deleteButtonTexture *giu.Texture
}

// Create creates a new ds1 viewer
func Create(id string, ds1 *d2ds1.DS1, dbt *giu.Texture, state []byte) giu.Widget {
	result := &widget{
		id:                  giu.ID(id),
		ds1:                 ds1,
		deleteButtonTexture: dbt,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding ds1 widget state: %v", err)
		}

		result.setState(s)
	}

	return result
}

// Build builds widget - implements giu.Widget
func (p *widget) Build() {
	state := p.getState()

	switch state.Mode {
	case widgetModeViewer:
		p.makeViewerLayout().Build()
	case widgetModeAddFile:
		p.makeAddFileLayout().Build()
	case widgetModeAddObject:
		p.makeAddObjectLayout().Build()
	case widgetModeAddPath:
		p.makeAddPathLayout().Build()
	case widgetModeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

// creates standard viewer/editor layout
func (p *widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	tabs := []*giu.TabItemWidget{
		giu.TabItem("Files").Layout(p.makeFilesLayout()),
		giu.TabItem("Objects").Layout(p.makeObjectsLayout(state)),
		giu.TabItem("Tiles").Layout(p.makeTilesTabLayout(state)),
	}

	if len(p.ds1.SubstitutionGroups) > 0 {
		tabs = append(tabs, giu.TabItem("Substitutions").Layout(p.makeSubstitutionsLayout(state)))
	}

	return giu.Layout{
		p.makeDataLayout(),
		giu.Separator(),
		giu.TabBar().TabItems(tabs...),
	}
}

// makeDataLayout creates basic data layout
// used in p.makeViewerLayout
func (p *widget) makeDataLayout() giu.Layout {
	version := int32(p.ds1.Version())

	state := p.getState()

	w, h := int32(p.ds1.Width()), int32(p.ds1.Height())
	l := giu.Layout{
		giu.Row(
			giu.Label("Version: "),
			giu.InputInt(&version).Size(inputIntW).OnChange(func() {
				state.confirmDialog = widgets.NewPopUpConfirmDialog(
					"##"+p.id+"confirmVersionChange",
					"Are you sure, you want to change DS1 Version?",
					"This value is used while decoding and encoding ds1 file\n"+
						"Please check github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1/ds1_version.go\n"+
						"to get more informations what does version determinates.\n\n"+
						"Continue?",
					func() {
						p.ds1.SetVersion(int(version))
						state.Mode = widgetModeViewer
					},
					func() {
						state.Mode = widgetModeViewer
					},
				)
				state.Mode = widgetModeConfirm
			}),
		),
		// giu.Label(fmt.Sprintf("Size: %d x %d tiles", p.ds1.Width, p.ds1.Height)),
		giu.Label("Size:"),
		giu.Row(
			giu.Label("\tWidth: "),
			giu.InputInt(&w).Size(inputIntW).OnChange(func() {
				state.confirmDialog = widgets.NewPopUpConfirmDialog(
					"##"+p.id+"confirmWidthChange",
					"Are you really sure, you want to change size of DS1 tiles?",
					"This will affect all your tiles in Tile tab.\n"+
						"Continue?",
					func() {
						p.ds1.SetWidth(int(w))
						state.Mode = widgetModeViewer
					},
					func() {
						state.Mode = widgetModeViewer
					},
				)
				state.Mode = widgetModeConfirm
			}),
		),
		giu.Row(
			giu.Label("\tHeight: "),
			giu.InputInt(&h).Size(inputIntW).OnChange(func() {
				state.confirmDialog = widgets.NewPopUpConfirmDialog(
					"##"+p.id+"confirmWidthChange",
					"Are you really sure, you want to change size of DS1 tiles?",
					"This will affect all your tiles in Tile tab.\n"+
						"Continue?",
					func() {
						p.ds1.SetHeight(int(h))
						state.Mode = widgetModeViewer
					},
					func() {
						state.Mode = widgetModeViewer
					},
				)
				state.Mode = widgetModeConfirm
			}),
		),
		giu.Label(fmt.Sprintf("Substitution Type: %d", p.ds1.SubstitutionType)),
		giu.Separator(),
		giu.Label("Number of"),
		giu.Label(fmt.Sprintf("\tWall Layers: %d", len(p.ds1.Walls))),
		giu.Label(fmt.Sprintf("\tFloor Layers: %d", len(p.ds1.Floors))),
		giu.Label(fmt.Sprintf("\tShadow Layers: %d", len(p.ds1.Shadows))),
		giu.Label(fmt.Sprintf("\tSubstitution Layers: %d", len(p.ds1.Substitutions))),
	}

	return l
}

// makeFilesLayout creates files list
// used in p.makeViewerLayout (files tab)
func (p *widget) makeFilesLayout() giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	// iterating using the value should not be a big deal as
	// we only expect a handful of strings in this slice.
	for n, str := range p.ds1.Files {
		currentIdx := n

		l = append(l, giu.Layout{
			giu.Row(
				widgets.MakeImageButton(
					deleteButtonSize, deleteButtonSize,
					p.deleteButtonTexture,
					func() {
						p.ds1.Files = append(p.ds1.Files[:currentIdx], p.ds1.Files[currentIdx+1:]...)
					},
				),
				giu.Label(str),
			),
		})
	}

	return giu.Layout{
		l,
		giu.Separator(),
		giu.Button("").ID("Add File##"+p.id+"AddFile").Size(actionButtonW, actionButtonH).OnClick(func() {
			state.Mode = widgetModeAddFile
		}),
	}
}

func (p *widget) makeAddFileLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Label("File path:"),
		giu.InputText(&state.NewFilePath).Size(filePathW),
		giu.Separator(),
		giu.Row(
			giu.Button("").ID("Add##"+p.id+"addFileAdd").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				p.ds1.Files = append(p.ds1.Files, state.NewFilePath)
				state.Mode = widgetModeViewer
			}),
			giu.Button("").ID("Cancel##"+p.id+"addFileCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.Mode = widgetModeViewer
			}),
		),
	}
}

// makeObjectsLayout creates Objects info tab
// used in p.makeViewerLayout (in Objects tab)
func (p *widget) makeObjectsLayout(state *widgetState) giu.Layout {
	numObjects := int32(len(p.ds1.Objects))

	l := giu.Layout{}

	if numObjects > 1 {
		l = append(l, giu.SliderInt(&state.Object, 0, numObjects-1).Label("Object Index"))
	}

	if numObjects > 0 {
		l = append(l, p.makeObjectLayout(state))
	} else {
		line := giu.Row(
			giu.Label("No Objects."),
			giu.ImageWithFile("hsassets/images/shrug.png").Size(imageW, imageH),
		)

		l = append(l, line)
	}

	l = append(
		l,
		giu.Separator(),
		giu.Row(
			giu.Button("").ID("Add new Object...##"+p.id+"AddObject").Size(actionButtonW, actionButtonH).OnClick(func() {
				state.Mode = widgetModeAddObject
			}),
			giu.Button("").ID("Add path to this Object...##"+p.id+"AddPath").Size(actionButtonW, actionButtonH).OnClick(func() {
				state.Mode = widgetModeAddPath
			}),
			widgets.MakeImageButton(
				layerDeleteButtonSize, layerDeleteButtonSize,
				p.deleteButtonTexture,
				func() {
					p.ds1.Objects = append(p.ds1.Objects[:state.Object], p.ds1.Objects[state.Object+1:]...)
				},
			),
		),
	)

	return l
}

// makeObjectLayout creates informations about single Object
// used in p.makeObjectsLayout
func (p *widget) makeObjectLayout(state *widgetState) giu.Layout {
	if objIdx := int(state.Object); objIdx >= len(p.ds1.Objects) {
		state.ds1Controls.Object = int32(len(p.ds1.Objects) - 1)
		p.setState(state)
	} else if objIdx < 0 {
		state.ds1Controls.Object = 0
		p.setState(state)
	}

	obj := &p.ds1.Objects[int(state.ds1Controls.Object)]

	l := giu.Layout{
		giu.Row(
			giu.Label("Type: "),
			widgets.MakeInputInt(
				inputIntW,
				&obj.Type,
				nil,
			),
		),
		giu.Row(
			giu.Label("ID: "),
			widgets.MakeInputInt(
				inputIntW,
				&obj.ID,
				nil,
			),
		),
		giu.Label("Position (tiles): "),
		giu.Row(
			giu.Label("\tX: "),
			widgets.MakeInputInt(
				inputIntW,
				&obj.X,
				nil,
			),
		),
		giu.Row(
			giu.Label("\tY: "),
			widgets.MakeInputInt(
				inputIntW,
				&obj.Y,
				nil,
			),
		),
		giu.Row(
			giu.Label("Flags: 0x"),
			widgets.MakeInputInt(
				inputIntW,
				&obj.Flags,
				nil,
			),
		),
	}

	if len(obj.Paths) > 0 {
		const spacerHeight = 16

		vspace := giu.Dummy(1, spacerHeight)
		l = append(l, vspace, p.makePathLayout(state, obj))
	}

	return l
}

func (p *widget) makeAddObjectLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Row(
			giu.Label("Type: "),
			giu.InputInt(&state.addObjectState.ObjType).Size(inputIntW),
		),
		giu.Row(
			giu.Label("ID: "),
			giu.InputInt(&state.addObjectState.ObjID).Size(inputIntW),
		),
		giu.Row(
			giu.Label("X: "),
			giu.InputInt(&state.addObjectState.ObjX).Size(inputIntW),
		),
		giu.Row(
			giu.Label("Y: "),
			giu.InputInt(&state.addObjectState.ObjY).Size(inputIntW),
		),
		giu.Row(
			giu.Label("Flags: "),
			giu.InputInt(&state.addObjectState.ObjFlags).Size(inputIntW),
		),
		giu.Separator(),
		giu.Row(
			giu.Button("").ID("Save##"+p.id+"AddObjectSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				newObject := d2ds1.Object{
					Type:  int(state.addObjectState.ObjType),
					ID:    int(state.addObjectState.ObjID),
					X:     int(state.addObjectState.ObjX),
					Y:     int(state.addObjectState.ObjY),
					Flags: int(state.addObjectState.ObjFlags),
				}

				p.ds1.Objects = append(p.ds1.Objects, newObject)

				state.Mode = widgetModeViewer
			}),
			giu.Button("").ID("Cancel##"+p.id+"AddObjectCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.Mode = widgetModeViewer
			}),
		),
	}
}

// makePathLayout creates paths table
// used in p.makeObjectLayout
func (p *widget) makePathLayout(state *widgetState, obj *d2ds1.Object) giu.Layout {
	rowWidgets := make([]*giu.TableRowWidget, 0)

	rowWidgets = append(rowWidgets, giu.TableRow(
		giu.Label("Index"),
		giu.Label("Position"),
		giu.Label("Action"),
		giu.Label(""),
	))

	for idx := range obj.Paths {
		currentIdx := idx
		x, y := obj.Paths[idx].Position.X(), obj.Paths[idx].Position.Y()
		rowWidgets = append(rowWidgets, giu.TableRow(
			giu.Label(fmt.Sprintf("%d", idx)),
			giu.Label(fmt.Sprintf("(%d, %d)", int(x), int(y))),
			giu.Label(fmt.Sprintf("%d", obj.Paths[idx].Action)),
			widgets.MakeImageButton(
				deleteButtonSize, deleteButtonSize,
				p.deleteButtonTexture,
				func() {
					p.ds1.Objects[state.Object].Paths = append(p.ds1.Objects[state.Object].Paths[:currentIdx],
						p.ds1.Objects[state.Object].Paths[currentIdx+1:]...)
				},
			),
		))
	}

	return giu.Layout{
		giu.Label("Path Points:"),
		giu.Table().FastMode(true).Rows(rowWidgets...),
	}
}

// makeTilesTabLayout creates tiles layout (tile x, y)
func (p *widget) makeTilesTabLayout(state *widgetState) giu.Layout {
	l := giu.Layout{}

	tx, ty := int(state.TileX), int(state.TileY)

	if ty < 0 {
		state.ds1Controls.TileY = 0
		p.setState(state)
	}

	if tx < 0 {
		state.ds1Controls.TileX = 0
		p.setState(state)
	}

	numRows := p.ds1.Height()
	if numRows == 0 {
		return l
	}

	if ty >= numRows {
		state.ds1Controls.TileY = int32(numRows - 1)
		p.setState(state)
	}

	if numCols := p.ds1.Width(); tx >= numCols {
		state.ds1Controls.TileX = int32(numCols - 1)
		p.setState(state)
	}

	tx, ty = int(state.TileX), int(state.TileY)

	l = append(
		l,
		giu.SliderInt(&state.ds1Controls.TileX, 0, int32(p.ds1.Width()-1)).Label("Tile X"),
		giu.SliderInt(&state.ds1Controls.TileY, 0, int32(p.ds1.Height()-1)).Label("Tile Y"),
		giu.TabBar().TabItems(
			p.makeTilesGroupLayout(state, tx, ty, d2ds1.FloorLayerGroup),
			p.makeTilesGroupLayout(state, tx, ty, d2ds1.WallLayerGroup),
			p.makeTilesGroupLayout(state, tx, ty, d2ds1.ShadowLayerGroup),
			p.makeTilesGroupLayout(state, tx, ty, d2ds1.SubstitutionLayerGroup),
		),
	)

	return l
}

// makeTilesGroupLayout creates a tileS group layout
// used in makeTilesTabLayout
func (p *widget) makeTilesGroupLayout(state *widgetState, x, y int, t d2ds1.LayerGroupType) *giu.TabItemWidget {
	l := giu.Layout{}
	group := p.ds1.GetLayersGroup(t)
	numRecords := len(*group)

	// this is a pointer to appropriate record index
	var recordIdx *int32
	// addCb is a callback for layer-add button
	var addCb func(int32)
	// delCb is a callback for layer-delete button
	var deleteCb func(int32)

	// sets "everything" ;-)
	switch t {
	case d2ds1.FloorLayerGroup:
		recordIdx = &state.Tile.Floor
		addCb = p.addFloor
		deleteCb = p.deleteFloor
	case d2ds1.WallLayerGroup:
		recordIdx = &state.Tile.Wall
		addCb = p.addWall
		deleteCb = p.deleteWall
	case d2ds1.ShadowLayerGroup:
		recordIdx = &state.Tile.Shadow
	case d2ds1.SubstitutionLayerGroup:
		recordIdx = &state.Tile.Sub
	}

	var addBtn *giu.ButtonWidget
	if addCb != nil {
		addBtn = giu.Button("").ID("Add "+giu.ID(t.String())+" ##"+p.id+"addButton").
			Size(actionButtonW, actionButtonH).
			OnClick(func() { addCb(*recordIdx) })
	}

	var deleteBtn giu.Widget
	if deleteCb != nil {
		deleteBtn = widgets.MakeImageButton(
			layerDeleteButtonSize, layerDeleteButtonSize,
			p.deleteButtonTexture,
			func() {
				deleteCb(*recordIdx)
			},
		)
	}

	if numRecords > 0 {
		// checks, if record index is correct
		if int(*recordIdx) >= numRecords {
			*recordIdx = int32(numRecords - 1)

			p.setState(state)
		} else if *recordIdx < 0 {
			*recordIdx = 0

			p.setState(state)
		}

		if numRecords > 1 {
			l = append(l, giu.SliderInt(recordIdx, 0, int32(numRecords-1)).Label(t.String()))
		}

		l = append(l, p.makeTileLayout((*group)[*recordIdx].Tile(x, y), t))
	}

	return giu.TabItem(t.String()).Layout(giu.Layout{
		l,
		giu.Separator(),
		giu.Custom(func() {
			var l giu.Layout
			if btn := addBtn; btn != nil {
				l = append(l, btn)
			}
			if btn := deleteBtn; btn != nil && numRecords > 0 {
				l = append(l, btn)
			}
			giu.Row(l...).Build()
		}),
	})
}

// makeTileLayout creates a single tile's layout
func (p *widget) makeTileLayout(record *d2ds1.Tile, t d2ds1.LayerGroupType) giu.Layout {
	// for substitutions, only unknown bytes should be displayed
	if t == d2ds1.SubstitutionLayerGroup {
		unknown32 := int32(record.Substitution)

		return giu.Layout{
			giu.Row(
				giu.Label("Substitute value: "),
				giu.InputInt(&unknown32).Size(inputIntW).OnChange(func() {
					record.Substitution = uint32(unknown32)
				}),
			),
		}
	}

	// common for shadows/walls/floors (like d2ds1.tileCommonFields)
	l := giu.Layout{
		giu.Row(
			giu.Label("Prop1: "),
			widgets.MakeInputInt(
				inputIntW,
				&record.Prop1,
				nil,
			),
		),
		giu.Row(
			giu.Label("Sequence: "),
			widgets.MakeInputInt(
				inputIntW,
				&record.Sequence,
				nil,
			),
		),
		giu.Row(
			giu.Label("Unknown1: "),
			widgets.MakeInputInt(
				inputIntW,
				&record.Unknown1,
				nil,
			),
		),
		giu.Row(
			giu.Label("Style: "),
			widgets.MakeInputInt(
				inputIntW,
				&record.Style,
				nil,
			),
		),
		giu.Row(
			giu.Label("Unknown2: "),
			widgets.MakeInputInt(
				inputIntW,
				&record.Unknown2,
				nil,
			),
		),
		giu.Row(
			giu.Label("Hidden: "),
			widgets.MakeCheckboxFromByte(
				"##"+p.id+"floorHidden",
				&record.HiddenBytes,
			),
		),
		giu.Row(
			giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		),
		giu.Row(
			giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
		),
	}

	switch t {
	case d2ds1.WallLayerGroup:
		l = append(l,
			giu.Row(
				giu.Label("Zero: "),
				widgets.MakeInputInt(
					inputIntW,
					&record.Zero,
					nil,
				),
			),
		)
	case d2ds1.FloorLayerGroup, d2ds1.ShadowLayerGroup:
		l = append(l,
			giu.Row(
				giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
			),
		)
	}

	return l
}

func (p *widget) makeSubstitutionsLayout(state *widgetState) giu.Layout {
	l := giu.Layout{}

	recordIdx := int(state.Subgroup)
	numRecords := len(p.ds1.SubstitutionGroups)

	if p.ds1.SubstitutionGroups == nil || numRecords == 0 {
		return l
	}

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.Subgroup = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.Subgroup = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt(&state.Subgroup, 0, int32(numRecords-1)).Label("Substitution"))
	}

	l = append(l, p.makeSubstitutionLayout(&p.ds1.SubstitutionGroups[recordIdx]))

	return l
}

func (p *widget) makeSubstitutionLayout(group *d2ds1.SubstitutionGroup) giu.Layout {
	l := giu.Layout{
		giu.Label(fmt.Sprintf("TileX: %d", group.TileX)),
		giu.Label(fmt.Sprintf("TileY: %d", group.TileY)),
		giu.Label(fmt.Sprintf("WidthInTiles: %d", group.WidthInTiles)),
		giu.Label(fmt.Sprintf("HeightInTiles: %d", group.HeightInTiles)),
		giu.Label(fmt.Sprintf("Unknown: 0x%x", group.Unknown)),
	}

	return l
}

func (p *widget) makeAddPathLayout() giu.Layout {
	state := p.getState()

	// https://github.com/OpenDiablo2/OpenDiablo2/issues/811
	// this list should be created like in COFWidget.makeAddLayerLayout
	actionsList := []string{"1", "2", "3"}

	return giu.Layout{
		giu.Row(
			giu.Label("Action: "),
			giu.Combo("",
				actionsList[state.addPathState.PathAction],
				actionsList, &state.addPathState.PathAction,
			).Size(bigListW).ID(
				"##"+p.id+"newPathAction",
			),
		),
		giu.Label("Vector:"),
		giu.Row(
			giu.Label("\tX: "),
			giu.InputInt(&state.addPathState.PathX).Size(inputIntW),
		),
		giu.Row(
			giu.Label("\tY: "),
			giu.InputInt(&state.addPathState.PathY).Size(inputIntW),
		),
		giu.Separator(),
		giu.Row(
			giu.Button("").ID("Save##"+p.id+"AddPathSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				p.addPath()
				state.Mode = widgetModeViewer
			}),
			giu.Button("").ID("Cancel##"+p.id+"AddPathCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.Mode = widgetModeViewer
			}),
		),
	}
}

func (p *widget) addPath() {
	state := p.getState()

	newPath := d2path.Path{
		// npc actions starts from 1
		Action: int(state.addPathState.PathAction) + 1,
		Position: d2vector.NewPosition(
			float64(state.addPathState.PathX),
			float64(state.addPathState.PathY),
		),
	}

	p.ds1.Objects[state.Object].Paths = append(p.ds1.Objects[state.Object].Paths, newPath)
}
