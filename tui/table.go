package tui

import (
	"github.com/gdamore/tcell/v2"
	"sort"
	"listdbg2/listdb"
	"github.com/rsn604/taps"
)

// ---------------------------------------
// Panel
// ---------------------------------------
func TablePanel() *taps.Panel {
	var styleMatrix = [][]string{
		{"list", "white", "default"},
		{"list_focus", "white", "aqua"},
	}

	var doc = `
StartX = 20
StartY = 2
EndX = 50
EndY = 19
Rect = true

[[Field]]	
Name = "TABLE"
X = 2
Y = 1
Rows = 16
Style = "list, list_focus"
FieldType = "select"
`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type TableApp struct {
	panel *taps.Panel
}

// -------------------------------------------------
//
//	Table
//
// -------------------------------------------------
func getTable(common *Common) ([]string, int) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	dbNames, _ := manager.GetDbNames()
	sort.Strings(dbNames)
	manager.Close()

	current := 0
	for i, c := range dbNames {
		if c == common.tableName {
			current = i
			break
		}
	}
	
	return dbNames, current
}

// ---------------------------------------
// Main
// ---------------------------------------
func (t *TableApp) Run(common *Common) string {
	if t.panel == nil{
		t.panel = TablePanel()
	}
	dbNames, current := getTable(common)
	t.panel.StoreList(dbNames, "TABLE")
		if current == 0{
		t.panel.SelectFocus = 0
	}
	t.panel.SetListStart("TABLE", current-t.panel.SelectFocus)
	t.panel.Say()
	k, n := t.panel.Read()
	if k == tcell.KeyEnter {
		return t.panel.Get(n)
	}
	return ""
}
