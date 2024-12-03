package app

import (
	g "github.com/AllenDang/giu"
	"github.com/gucio321/HellSpawner/pkg/window"
)

func renderWnd(d window.Renderable) g.Widget {
	return g.Custom(func() {
		if d.IsVisible() {
			d.Build()
		}
	})
}

// TODO: selected tab doesn't work here
func (a *App) renderStaticEditors() g.Widget {
	tabs := make([]*g.TabItemWidget, 0)

	idx := 0
	for idx < len(a.editors) {
		editor := a.editors[idx]
		if !editor.IsVisible() {
			editor.Cleanup()

			a.editors = append(a.editors[:idx], a.editors[idx+1:]...)

			continue
		}

		hadFocus := editor.HasFocus()

		// common shortcut
		editor.RegisterKeyboardShortcuts(
			g.WindowShortcut{
				Key:      g.KeyS,
				Modifier: g.ModControl,
				Callback: func() {
					editor.Save()
				},
			},
		)

		editor.RegisterKeyboardShortcuts(
			editor.KeyboardShortcuts()...,
		)

		tabs = append(tabs, g.TabItemf("Editor %d", idx).Layout(
			editor.GetLayout(),
		))

		if !hadFocus && editor.HasFocus() {
			a.focusedEditor = editor
		}

		idx++
	}
	return g.TabBar().TabItems(tabs...)
}
