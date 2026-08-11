package tui

import (
	"github.com/gdamore/tcell/v2"
	"listdbg2/listdb"
	"github.com/rsn604/taps"
)

// ---------------------------------------
// Panel
// ---------------------------------------
func CategoryPanel() *taps.Panel {
	var styleMatrix = [][]string{
		{"list", "white", "default"},
		{"list_focus", "white", "aqua"},
	}

	var doc = `
StartX = 20
StartY = 4
EndX = 50
EndY = 17
Rect = true

[[Field]]	
Name = "CATEGORY"
X = 2
Y = 1
Rows = 12
Style = "list, list_focus"
FieldType = "select"

`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type CategoryApp struct {
	panel *taps.Panel
}

// -------------------------------------------------
//
//	Table
//
// -------------------------------------------------
func getCategory(common *Common) ([]string, int) {
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	categoryList, _ := manager.GetCategoryList(common.tableName)
	manager.Close()
	categoryList = append([]string{"Clear"}, categoryList...)
	current := 0
	for i, c := range categoryList {
		if c == common.category {
			current = i
			break
		}
	}
	return categoryList, current
}

// ---------------------------------------
// Main
// ---------------------------------------
func (t *CategoryApp) Run(common *Common) string {
	if t.panel == nil{
		t.panel = CategoryPanel()
	}
	categoryList, current := getCategory(common)
	t.panel.StoreList(categoryList, "CATEGORY")
	if current == 0{
		t.panel.SelectFocus = 0
	}
	t.panel.SetListStart("CATEGORY", current-t.panel.SelectFocus)
	t.panel.Say()
	k, n := t.panel.Read()
	if k == tcell.KeyEnter {
		return t.panel.Get(n)
	}
	return ""
}
