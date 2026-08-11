package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/toqueteos/webbrowser"
	"os/exec"
	"strconv"
	"strings"
	"listdbg2/listdb"
	"github.com/rsn604/taps"
)

func DetailPanel() *taps.Panel {
	var styleMatrix = [][]string{
		{"label", "aqua", "default"},
		{"select", "yellow", "default"},
		{"select_focus", "black", "yellow"},
		{"edit", "white", "black"},
		{"edit_focus", "yellow,underline", "black"},
		{"GOMAP", "aqua", "default"},
		{"GOMAP_FOCUS", "black", "red"},
	}

	var doc = `
StartX = 0
StartY = 0
EndX = 9999
EndY = 9999

[[Field]]	
Name = "R"
Data = "<R>"
X = 0
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]
Name = "N"
Data = "<N>"
X = 10
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "P"
Data = "<P>"
X = 20
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "Q"
Data = "<Q>"
X = 30
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]
Name = "L_ID01"
Data = "ID"
X = 0
Y = 1
Style = "label"
FieldType = "label"

[[Field]]
Name = "L_ID02"
#Data = "ID"
X = 0
Y = 2
Style = "label"
FieldType = "label"

[[Field]]
Name = "L_CATEGORY"
Data = "Category"
X = 0
Y = 3
Style = "label"
FieldType = "label"

[[Field]]
Name = "E_CATEGORY"
X = 0
Y = 4
#Style = "edit, edit_focus"
Style = "select, select_focus"
FieldLen = 10
#FieldType = "edit"
FieldType = "select"

[[Field]]
Name = "L_FIELD01"
#Data = "Title"
X = 0
Y = 5
Style = "label"
FieldType = "label"

[[Field]]
Name = "E_FIELD01"
X = 0
Y = 6
Style = "edit, edit_focus"
FieldLen = 9999
FieldType = "edit"

[[Field]]
Name = "L_FIELD02"
X = 0
Y = 7
Style = "label"
FieldType = "label"

[[Field]]
Name = "E_FIELD02"
X = 0
Y = 8
Style = "edit, edit_focus"
FieldLen = 9999
FieldType = "edit"

[[Field]]
Name = "L_NOTE"
Data = "Note"
X = 0
Y = 9
Style = "label"
FieldType = "label"

[[Field]]
Name = "E_NOTE"
X = 0
Y = 10
Style = "edit, edit_focus"
FieldLen = 9998
Rows = 9988
FieldType = "edit"

[[Field]]
Name = "D"
Data = "<D>"
X = 0
Y = 9999
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "U"
Data = "<U>"
X = 10
Y = 9999
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "I"
Data = "<I>"
X = 20
Y = 9999
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"
`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type DetailApp struct {
	panel *taps.Panel
}

// -------------------------------------------------
//
//	Set field to invoke Browser
//
// -------------------------------------------------
func getIntentPath(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		return ""
	}
	return path
}

func execCommand(t string) {
	if getIntentPath("am") == "" {
		webbrowser.Open("http://maps.google.co.jp/maps?q=" + t)
	} else {
		parm := strings.Split("start -a android.intent.action.VIEW -d geo:0,0?q="+t, " ")
		_ = exec.Command("am", parm...).Run()

	}
}

// -------------------------------------------------
//	get Manager
// -------------------------------------------------
func (m *DetailApp) getManager(common *Common) listdb.Manager{
	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}
	return manager
}

func (m *DetailApp) createListItem(listItem *listdb.ListItem) *listdb.ListItem {
	listItem.ID = listItem.ID
	listItem.Category = m.panel.Get("E_CATEGORY")
	listItem.Field01 =  m.panel.Get("E_FIELD01")
	listItem.Field02 =  m.panel.Get("E_FIELD02")
	listItem.Note = strings.Join(m.panel.GetList("E_NOTE"), "\n")
	return listItem
}

// -------------------------------------------------
//
//	Update
//
// -------------------------------------------------
func (m *DetailApp) update(common *Common, listItem *listdb.ListItem) {
	manager := m.getManager(common)
	listItem = m.createListItem(listItem)
	manager.Update(common.tableName, listItem.ID, listItem)
	manager.Close()
}

// -------------------------------------------------
//
//	Delete
//
// -------------------------------------------------
func (m *DetailApp) delete(common *Common, listItem *listdb.ListItem) {
	manager := m.getManager(common)
	manager.Delete(common.tableName, listItem.ID)
	manager.Close()
}

// -------------------------------------------------
//
//	Insert
//
// -------------------------------------------------
func (m *DetailApp) insert(common *Common, listItem *listdb.ListItem) {
	manager := m.getManager(common)
	listItem = m.createListItem(listItem)
	listItem.ID, _ =  manager.Insert(common.tableName, listItem)
	manager.Update(common.tableName, listItem.ID, listItem)
	manager.Close()
}

// -------------------------------------------------
//
//	Paging
//
// -------------------------------------------------
var isLastItem bool

func (m *DetailApp) firstPage(common *Common) bool {
	return common.selectedItem == 1
}

func (m *DetailApp) lastPage() bool {
	return isLastItem
}

func (m *DetailApp) nextPage(common *Common) {
	if !m.lastPage() {
		common.selectedItem += 1
	}
}

func (m *DetailApp) priorPage(common *Common) {
	if !m.firstPage(common) {
		common.selectedItem -= 1
	}
	if common.selectedItem < 1 {
		common.selectedItem = 1
	}
}

func (m *DetailApp) isDataExist(common *Common) bool {
	return common.tableName != ""
}

// ---------------------------------------
// Main
// ---------------------------------------
func (m *DetailApp) doFormat(common *Common) listdb.ListItem{
	manager := m.getManager(common)
	isLastItem = false
	listdb := manager.SearchDB(common.tableName, common.category, common.search, common.selectedItem, 2)
	listdata := listdb.GetListData()
	if len(listdata) < 2 {
		isLastItem = true
	}
	manager.Close()

	// -----------------------------
	m.panel.Store(strconv.Itoa(listdata[0].ID), "L_ID02")
	m.panel.Store(listdata[0].Category, "E_CATEGORY")


	if strings.ToLower(listdb.FieldName01) == "address" || strings.ToLower(listdb.FieldName01) == "map" {
		f := m.panel.GetDataField("L_FIELD01")
		f.FieldType = "SELECT"
		m.panel.Store(listdb.FieldName01+" -> MAP", "L_FIELD01")
	}else{
		m.panel.Store(listdb.FieldName01, "L_FIELD01")
		m.panel.ResetFieldStyle("L_FIELD01", "GOMAP, GOMAP_FOCUS")
	}
	
	if strings.ToLower(listdb.FieldName02) == "address" || strings.ToLower(listdb.FieldName02) == "map" {
		f  := m.panel.GetDataField("L_FIELD02")
		f.FieldType = "SELECT"
		m.panel.Store(listdb.FieldName02+" -> MAP", "L_FIELD02")
		m.panel.ResetFieldStyle("L_FIELD02", "GOMAP, GOMAP_FOCUS")
	}else{
		m.panel.Store(listdb.FieldName02, "L_FIELD02")
	}

	m.panel.Store(listdata[0].Field01, "E_FIELD01")
	m.panel.Store(listdata[0].Field02, "E_FIELD02")
	m.panel.StoreList(strings.Split(listdata[0].Note, "\n"), "E_NOTE")

	if m.lastPage() {
		m.panel.Store("", "N")
	} else {
		m.panel.Store("<N>", "N")
	}
	if m.firstPage(common) {
		m.panel.Store("", "P")
	} else {
		m.panel.Store("<P>", "P")
	}

	m.panel.Say()
	return listdata[0]
}

func (m *DetailApp) Run(common *Common) string {
	m.panel = DetailPanel()

	category := &CategoryApp{}
	confirm := &Confirm{}

	for {
		listItem := m.doFormat(common)
		k, n := m.panel.Read()

		if k == tcell.KeyEscape {
			break
		}

		//--------------------------------		
		if n == "R" {
			break
		}
		if n == "N" {
			m.nextPage(common)
		}
		if n == "P" {
			m.priorPage(common)
		}

		if n == "Q" {
			return "Q"
		}

		//--------------------------------		
 		if n == "D" {
			rs := confirm.Run("Delete record ?")
			if rs == "OK" {
				m.delete(common, &listItem)
				return ""
			}
		}
		
		if n == "U" {
			rs := confirm.Run("Update record ?")
			if rs == "OK" {
				m.update(common, &listItem)
			}
		}
		
		if n == "I" {
			rs := confirm.Run("Insert record ?")
			if rs == "OK" {
				m.insert(common, &listItem)
				return ""
			}
		}

		//--------------------------------
		if n == "E_CATEGORY" {
			rc := category.Run(common)
			if rc != "Clear" && rc != "" {
				m.panel.Store(rc, "E_CATEGORY")
			}
		}

		if n == "L_FIELD01" {
			execCommand(m.panel.Get("E_FIELD01"))
		}
		
		if n == "L_FIELD02" {
			execCommand(m.panel.Get("E_FIELD02"))
		}
	}
	return ""

}
