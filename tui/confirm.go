package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rsn604/taps"
)

// ---------------------------------------
// Panel
// ---------------------------------------
func ConfirmPanel() *taps.Panel {

	var styleMatrix = [][]string{
		{"label", "yellow", "default"},
		{"select", "yellow", "default"},
		{"select_focus", "black", "yellow"},
	}

	var doc = `
StartX = 10
StartY = 5
EndX = 37
EndY = 11
Rect = true

[[Field]]
Name = "L_CONFIRM"
#Data = "Confirm : "
X = 2
Y = 2
FieldLen = 20
Style = "label"
FieldType = "label"


[[Field]]
Name = "OK"
Data = "OK"
X = 6
Y = 4
Style = "select, select_focus"
FieldType = "select"

[[Field]]
Name = "Cancel"
Data = "Cancel"
X = 15
Y = 4
Style = "select, select_focus"
FieldType = "select"
`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type Confirm struct {
	panel *taps.Panel
}

// ---------------------------------------
// Main
// ---------------------------------------
func (m *Confirm) Run(confirmMsg string) string {
	if m.panel == nil {
		m.panel = ConfirmPanel()
	}
	m.panel.Store(confirmMsg, "L_CONFIRM")
	m.panel.Say()
	k, n := m.panel.Read()
	if k == tcell.KeyEnter {
		if n == "OK" {
			return "OK"
		}
		if n == "Cancel" {
			return ""
		}
	}
	return ""
}
