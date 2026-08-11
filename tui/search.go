package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rsn604/taps"
)

// ---------------------------------------
// Panel
// ---------------------------------------
func SearchPanel() *taps.Panel {
	var styleMatrix = [][]string{
		{"label", "yellow", "default"},
		{"select", "yellow", "default"},
		{"select_focus", "black", "yellow"},
		{"edit", "white,underline", "black"},
		{"edit_focus", "yellow", "black"},
	}

	var doc = `
StartX = 10
StartY = 5
EndX = 47
EndY = 11
Rect = true

[[Field]]
Name = "L01"
Data = "Search : "
X = 2
Y = 2
Style = "label"
FieldType = "label"

[[Field]]
Name = "S01"
X = 12
Y = 2
Style = "edit, edit_focus"
FieldLen = 22
FieldType = "edit"
ExitKey = ["Enter"]

[[Field]]
Name = "OK"
Data = "OK"
X = 12
Y = 4
Style = "select, select_focus"
FieldType = "select"

[[Field]]
Name = "Cancel"
Data = "Cancel"
X = 19
Y = 4
Style = "select, select_focus"
FieldType = "select"
`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type SearchApp struct {
	panel *taps.Panel
}

// ---------------------------------------
// Main
// ---------------------------------------
func (t *SearchApp) Run() string {
	if t.panel == nil{
		t.panel = SearchPanel()
	}
	t.panel.Say()
	k, n := t.panel.Read()
	if k == tcell.KeyEnter {
		if n == "OK" || n == "S01" {
			return t.panel.Get("S01")
		}
		if n == "Cancel" {
			t.panel.Store("", "S01")
		}
	}
	return ""
}
