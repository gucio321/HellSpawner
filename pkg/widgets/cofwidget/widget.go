package cofwidget

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/gucio321/HellSpawner/pkg/widgets"
)

const (
	layerListW                           = 64
	buttonWidthHeight                    = 15
	actionButtonW, actionButtonH         = 200, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	speedInputW                          = 40
)

type widget struct {
	id  giu.ID
	cof *d2cof.COF
}

// Create a new COF widget
func Create(
	state []byte,
	id string, cof *d2cof.COF,
) giu.Widget {
	result := &widget{
		id:  giu.ID(id),
		cof: cof,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding cof widget state: %v", err)
		}

		result.setState(s)
	}

	return result
}

// Build builds a cof viewer
func (p *widget) Build() {
	state := p.getState()

	// builds appropriate menu (depends on state)
	switch state.Mode {
	case modeViewer:
		p.makeViewerLayout().Build()
	case modeAddLayer:
		p.makeAddLayerLayout().Build()
	case modeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.TabBar().TabItems(
			giu.TabItem("Animation").Layout(p.makeAnimationTab(state)),
			giu.TabItem("Layer").Layout(p.makeLayerTab(state)),
			giu.TabItem("Priority").Layout(p.makePriorityTab(state)),
		),
	}
}

func (p *widget) makeAnimationTab(state *widgetState) giu.Layout {
	const (
		fmtFPS        = "FPS: %.1f"
		fmtDuration   = "Duration: %.2fms"
		fmtDirections = "Directions: %v"
		strSpeed      = "Speed: "
		maxSpeed      = 256
	)

	numDirs := p.cof.NumberOfDirections
	fps := speedToFPS(p.cof.Speed)
	duration := calculateDuration(p.cof)

	strLabelDirections := fmt.Sprintf(fmtDirections, numDirs)
	strLabelFPS := fmt.Sprintf(fmtFPS, fps)
	strLabelDuration := fmt.Sprintf(fmtDuration, duration)

	setSpeed := func() {
		if p.cof.Speed >= maxSpeed {
			p.cof.Speed = maxSpeed
		}
	}

	speedLabel := giu.Label(strSpeed)
	speedInput := widgets.MakeInputInt(
		speedInputW,
		&p.cof.Speed,
		setSpeed,
	)

	return giu.Layout{
		giu.Label(strLabelDirections),
		p.layoutAnimFrames(state),
		giu.Row(speedLabel, speedInput),
		giu.Label(strLabelFPS),
		giu.Label(strLabelDuration),
	}
}

func (p *widget) makeLayerTab(state *widgetState) giu.Layout {
	addLayerButtonID := fmt.Sprintf("Add a new layer...##%sAddLayer", p.id)
	addLayerButton := giu.Button(addLayerButtonID).Size(actionButtonW, actionButtonH)
	addLayerButton.OnClick(func() {
		p.createNewLayer()
	})

	if state.viewerState.layer == nil {
		return giu.Layout{addLayerButton}
	}

	layerStrings := make([]string, 0)
	for idx := range p.cof.CofLayers {
		layerStrings = append(layerStrings, strconv.Itoa(int(p.cof.CofLayers[idx].Type)))
	}

	currentLayerName := layerStrings[state.viewerState.LayerIndex]
	layerList := giu.Combo("", currentLayerName, layerStrings, &state.LayerIndex).ID("##" + p.id + "layer")
	layerList.Size(layerListW).OnChange(p.onUpdate)

	deleteLayerButtonID := fmt.Sprintf("Delete current layer...##%sDeleteLayer", p.id)
	deleteLayerButton := giu.Button(deleteLayerButtonID).Size(actionButtonW, actionButtonH)
	deleteLayerButton.OnClick(func() {
		const (
			strPrompt  = "Do you really want to remove this layer?"
			strMessage = "If you'll click YES, all data from this layer will be lost. Continue?"
		)

		fnYes := func() {
			p.deleteCurrentLayer(state.viewerState.LayerIndex)
			state.Mode = modeViewer
		}

		fnNo := func() {
			state.Mode = modeViewer
		}

		id := giu.ID(fmt.Sprintf("##%sDeleteLayerConfirm", p.id))
		state.viewerState.confirmDialog = widgets.NewPopUpConfirmDialog(id, strPrompt, strMessage, fnYes, fnNo)

		state.Mode = modeConfirm
	})

	layout := giu.Layout{
		giu.Row(giu.Label("Selected Layer: "), layerList),
		giu.Separator(),
		p.makeLayerLayout(),
		giu.Separator(),
		addLayerButton,
		deleteLayerButton,
	}

	return layout
}

func (p *widget) createNewLayer() {
	state := p.getState()

	state.Mode = modeAddLayer
}

func (p *widget) makePriorityTab(state *widgetState) giu.Layout {
	if len(p.cof.Priority) == 0 {
		return giu.Layout{
			giu.Label("Nothing here"),
		}
	}

	directionStrings := make([]string, 0)
	for idx := range p.cof.Priority {
		directionStrings = append(directionStrings, fmt.Sprintf("%d", idx))
	}

	directionString := directionStrings[state.viewerState.DirectionIndex]
	directionList := giu.Combo("", directionString, directionStrings, &state.DirectionIndex).
		ID("##" + p.id + "dir")
	directionList.Size(layerListW).OnChange(p.onUpdate)

	frameStrings := make([]string, 0)
	for idx := range p.cof.Priority[state.DirectionIndex] {
		frameStrings = append(frameStrings, fmt.Sprintf("%d", idx))
	}

	frameString := frameStrings[state.FrameIndex]
	frameList := giu.Combo("", frameString, frameStrings, &state.FrameIndex).
		ID("##" + p.id + "frame")
	frameList.Size(layerListW).OnChange(p.onUpdate)

	const (
		strPrompt  = "Do you really want to remove this direction?"
		strMessage = "If you'll click YES, all data from this direction will be lost. Continue?"
	)

	duplicateButtonID := fmt.Sprintf("Duplicate current direction...##%sDuplicateDirection", p.id)
	duplicateButton := giu.Button(duplicateButtonID).Size(actionButtonW, actionButtonH)
	duplicateButton.OnClick(func() {
		p.duplicateDirection()
	})

	deleteButtonID := fmt.Sprintf("Delete current direction...##%sDeleteDirection", p.id)
	deleteButton := giu.Button(deleteButtonID).Size(actionButtonW, actionButtonH)
	deleteButton.OnClick(func() {
		fnYes := func() {
			p.deleteCurrentDirection()

			state.Mode = modeViewer
		}

		fnNo := func() {
			state.Mode = modeViewer
		}

		popupID := giu.ID(fmt.Sprintf("%sDeleteLayerConfirm", p.id))

		state.confirmDialog = widgets.NewPopUpConfirmDialog(popupID, strPrompt, strMessage, fnYes, fnNo)
		state.Mode = modeConfirm
	})

	return giu.Layout{
		giu.Row(
			giu.Label("Direction: "), directionList,
			giu.Label("Frame: "), frameList,
		),
		giu.Separator(),
		p.makeDirectionLayout(),
		duplicateButton,
		deleteButton,
	}
}

// the layout ends up looking like this:
// Frames (x6):  <- 10 ->
// you use the arrows to set the number of frames per direction
func (p *widget) layoutAnimFrames(state *widgetState) *giu.RowWidget {
	numFrames := p.cof.FramesPerDirection
	numDirs := p.cof.NumberOfDirections

	strLabel := "Frames:"
	if numDirs > 1 {
		strLabel = fmt.Sprintf("Frames (x%v):", numDirs)
	}

	fnDecrease := func() {
		p.cof.FramesPerDirection = max(p.cof.FramesPerDirection-1, 0)
	}

	fnIncrease := func() {
		p.cof.FramesPerDirection++
	}

	label := giu.Label(strLabel)

	leftButtonID := giu.ID(fmt.Sprintf("##%sDecreaseFramesPerDirection", p.id))
	rightButtonID := giu.ID(fmt.Sprintf("##%sIncreaseFramesPerDirection", p.id))

	left := widgets.MakeImageButton(leftButtonID, buttonWidthHeight, buttonWidthHeight, state.textures.left, fnDecrease)
	frameCount := giu.Label(fmt.Sprintf("%d", numFrames))
	right := widgets.MakeImageButton(rightButtonID, buttonWidthHeight, buttonWidthHeight, state.textures.right, fnIncrease)

	return giu.Row(label, left, frameCount, right)
}

func (p *widget) onUpdate() {
	state := p.getState()

	clone := p.cof.CofLayers[state.LayerIndex]
	state.viewerState.layer = &clone

	giu.Context.SetState(p.id, state)
}

func (p *widget) makeLayerLayout() giu.Layout {
	state := p.getState()

	if state.viewerState.layer == nil {
		p.onUpdate()
	}

	layerName := state.viewerState.layer.Type.Name()

	strType := fmt.Sprintf("Type: %s (%s)", state.viewerState.layer.Type, layerName)
	strShadow := fmt.Sprintf("Shadow: %t", state.viewerState.layer.Shadow > 0)
	strSelectable := fmt.Sprintf("Selectable: %t", state.viewerState.layer.Selectable)
	strTransparent := fmt.Sprintf("Transparent: %t", state.viewerState.layer.Transparent)

	effect := state.viewerState.layer.DrawEffect.String()

	strEffect := fmt.Sprintf("Draw Effect: %s", effect)

	weapon := state.viewerState.layer.WeaponClass.String()

	strWeaponClass := fmt.Sprintf("Weapon Class: (%s) %s", state.viewerState.layer.WeaponClass, weapon)

	return giu.Layout{
		giu.Label(strType),
		giu.Label(strShadow),
		giu.Label(strSelectable),
		giu.Label(strTransparent),
		giu.Label(strEffect),
		giu.Label(strWeaponClass),
	}
}

func (p *widget) makeDirectionLayout() giu.Layout {
	const (
		strRenderOrderLabel = "Render Order (first to last):"
		fmtIncreasePriority = "LayerPriorityUp_%d"
		fmtDecreasePriority = "LayerPriorityDown_%d"
		fmtLayerLabel       = "%d: %s"
	)

	state := p.getState()

	frames := p.cof.Priority[state.DirectionIndex]
	layers := frames[int(state.FrameIndex)%len(frames)]

	// increase / decrease callback function providers, based on layer index
	makeIncPriorityFn := func(idx int) func() {
		return func() {
			if idx <= 0 {
				return
			}

			list := &p.cof.Priority[state.DirectionIndex][state.FrameIndex]
			(*list)[idx-1], (*list)[idx] = (*list)[idx], (*list)[idx-1]
		}
	}

	makeDecPriorityFn := func(idx int) func() {
		return func() {
			list := &p.cof.Priority[state.DirectionIndex][state.FrameIndex]

			if idx >= len(*list)-1 {
				return
			}

			(*list)[idx], (*list)[idx+1] = (*list)[idx+1], (*list)[idx]
		}
	}

	// each layer line looks like:
	// <- -> 0: Name
	// the left/right buttons use the callbacks created by the previous funcs for index=0
	buildLayerPriorityRow := func(idx int) {
		currentIdx := idx

		strIncPri := giu.ID(fmt.Sprintf(fmtIncreasePriority, currentIdx))
		strDecPri := giu.ID(fmt.Sprintf(fmtDecreasePriority, currentIdx))

		fnIncPriority := makeIncPriorityFn(currentIdx)
		fnDecPriority := makeDecPriorityFn(currentIdx)

		increasePriority := widgets.MakeImageButton(strIncPri, buttonWidthHeight, buttonWidthHeight, state.textures.up, fnIncPriority)
		decreasePriority := widgets.MakeImageButton(strDecPri, buttonWidthHeight, buttonWidthHeight, state.textures.down, fnDecPriority)

		strLayerName := layers[idx].Name()
		strLayerLabel := fmt.Sprintf(fmtLayerLabel, idx, strLayerName)

		layerNameLabel := giu.Label(strLayerLabel)

		giu.Row(increasePriority, decreasePriority, layerNameLabel).Build()
	}

	// finally, a func that we can pass to giu.Custom
	buildLayerRows := func() {
		for idx := range layers {
			buildLayerPriorityRow(idx)
		}
	}

	return giu.Layout{
		giu.Label(strRenderOrderLabel),
		giu.Custom(buildLayerRows),
	}
}

func (p *widget) makeAddLayerLayout() giu.Layout {
	state := p.getState()

	// available is a list of available (not currently used) composite types
	available := make([]d2enum.CompositeType, 0)

	for i := d2enum.CompositeTypeHead; i < d2enum.CompositeTypeMax; i++ {
		contains := false

		for _, j := range p.cof.CofLayers {
			if i == j.Type {
				contains = true

				break
			}
		}

		if !contains {
			available = append(available, i)
		}
	}

	compositeTypeList := make([]string, len(available))
	for n, i := range available {
		compositeTypeList[n] = i.String() + " (" + i.Name() + ")"
	}

	drawEffectList := make([]string, d2enum.DrawEffectNone+1)
	for i := d2enum.DrawEffectPctTransparency25; i <= d2enum.DrawEffectNone; i++ {
		drawEffectList[i] = strconv.Itoa(int(i)) + " (" + i.String() + ")"
	}

	weaponClassList := make([]string, d2enum.WeaponClassTwoHandToHand+1)
	for i := d2enum.WeaponClassNone; i <= d2enum.WeaponClassTwoHandToHand; i++ {
		weaponClassList[i] = i.String() + " (" + i.Name() + ")"
	}

	return giu.Layout{
		giu.Label("Select new COF's Layer parameters:"),
		giu.Separator(),
		giu.Row(
			giu.Label("Type: "),
			giu.Combo("", compositeTypeList[state.newLayerFields.LayerType],
				compositeTypeList, &state.newLayerFields.LayerType).Size(bigListW).ID(
				"##"+p.id+"AddLayerType",
			),
		),
		giu.Row(
			giu.Label("Shadow: "),
			widgets.MakeCheckboxFromByte("##"+p.id+"AddLayerShadow", &state.newLayerFields.Shadow),
		),
		giu.Row(
			giu.Label("Selectable: "),
			giu.Checkbox("", &state.newLayerFields.Selectable).ID("##"+p.id+"AddLayerSelectable"),
		),
		giu.Row(
			giu.Label("Transparent: "),
			giu.Checkbox("", &state.newLayerFields.Transparent).ID("##"+p.id+"AddLayerTransparent"),
		),
		giu.Row(
			giu.Label("Draw effect: "),
			giu.Combo("", drawEffectList[state.newLayerFields.DrawEffect],
				drawEffectList, &state.newLayerFields.DrawEffect).Size(bigListW).ID(
				"##"+p.id+"AddLayerDrawEffect",
			),
		),
		giu.Row(
			giu.Label("Weapon class: "),
			giu.Combo("", weaponClassList[state.newLayerFields.WeaponClass],
				weaponClassList, &state.newLayerFields.WeaponClass).Size(bigListW).ID(
				"##"+p.id+"AddLayerWeaponClass",
			),
		),
		giu.Separator(),
		p.makeSaveCancelButtonRow(available, state),
	}
}

func (p *widget) makeSaveCancelButtonRow(available []d2enum.CompositeType, state *widgetState) *giu.RowWidget {
	fnSave := func() {
		newCofLayer := &d2cof.CofLayer{
			Type:        available[state.newLayerFields.LayerType],
			Shadow:      state.newLayerFields.Shadow,
			Selectable:  state.newLayerFields.Selectable,
			Transparent: state.newLayerFields.Transparent,
			DrawEffect:  d2enum.DrawEffect(state.newLayerFields.DrawEffect),
			WeaponClass: d2enum.WeaponClass(state.newLayerFields.WeaponClass),
		}

		p.cof.CofLayers = append(p.cof.CofLayers, *newCofLayer)

		p.cof.NumberOfLayers++

		for dirIdx := range p.cof.Priority {
			for frameIdx := range p.cof.Priority[dirIdx] {
				p.cof.Priority[dirIdx][frameIdx] = append(p.cof.Priority[dirIdx][frameIdx], newCofLayer.Type)
			}
		}

		// this sets layer index to just added layer
		//nolint:gosec // this is for giu and has to be int32.
		state.viewerState.LayerIndex = int32(p.cof.NumberOfLayers - 1)
		state.viewerState.layer = newCofLayer
		state.Mode = modeViewer
	}

	fnCancel := func() {
		state.Mode = modeViewer
	}

	buttonSave := giu.Button("Save##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(fnSave)
	buttonCancel := giu.Button("Cancel##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(fnCancel)

	return giu.Row(buttonSave, buttonCancel)
}

func (p *widget) deleteCurrentLayer(index int32) {
	p.cof.NumberOfLayers--

	newPriority := make([][][]d2enum.CompositeType, p.cof.NumberOfDirections)

	for dn := range p.cof.Priority {
		newPriority[dn] = make([][]d2enum.CompositeType, p.cof.FramesPerDirection)
		for fn := range p.cof.Priority[dn] {
			newPriority[dn][fn] = make([]d2enum.CompositeType, p.cof.NumberOfLayers)

			for ln := range p.cof.Priority[dn][fn] {
				if p.cof.CofLayers[index].Type != p.cof.Priority[dn][fn][ln] {
					newPriority[dn][fn] = append(newPriority[dn][fn], p.cof.Priority[dn][fn][ln])
				}
			}
		}
	}

	p.cof.Priority = newPriority

	newLayers := make([]d2cof.CofLayer, 0)

	for n, i := range p.cof.CofLayers {
		if n != int(index) {
			newLayers = append(newLayers, i)
		}
	}

	p.cof.CofLayers = newLayers

	state := p.getState()

	if state.viewerState.LayerIndex != 0 {
		state.viewerState.LayerIndex--
	}
}

func (p *widget) duplicateDirection() {
	state := p.getState()

	idx := state.viewerState.DirectionIndex

	p.cof.NumberOfDirections++

	p.cof.Priority = append(p.cof.Priority, p.cof.Priority[idx])

	//nolint:gosec // this is for giu and has to be int32.
	state.DirectionIndex = int32(len(p.cof.Priority) - 1)
}

func (p *widget) deleteCurrentDirection() {
	state := p.getState()

	index := state.viewerState.DirectionIndex

	p.cof.NumberOfDirections--

	newPriority := make([][][]d2enum.CompositeType, 0)

	for n, i := range p.cof.Priority {
		if n != int(index) {
			newPriority = append(newPriority, i)
		}
	}

	p.cof.Priority = newPriority
}
