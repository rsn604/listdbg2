package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"strings"
	"listdbg2/listdb"
	"github.com/rsn604/taps"
	//"log"
)

func ListPanel() *taps.Panel {
	var styleMatrix = [][]string{
		{"label", "yellow", "default"},
		{"select", "yellow", "default"},
		{"select_focus", "black", "yellow"},
		{"list", "white", "default"},
		{"list_focus", "white", "aqua"},
	}

	var doc = `
StartX = 0
StartY = 0
EndX = 9999
EndY = 9999

[[Field]]	
Name = "T"
Data = "<T>"
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
Name = "C"
Data = "<C>"
X = 30
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "S"
Data = "<S>"
X = 40
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

[[Field]]	
Name = "Q"
Data = "<Q>"
X = 50
Y = 0
FieldLen = 3
Style = "select, select_focus"
FieldType = "select"

# -------------------------------------------------
[[Field]]	
Name = "LIST"
X = 0
Y = 1
Rows = 9998
Style = "list, list_focus"
FieldType = "select"

# -------------------------------------------------
[[Field]]
Name = "L01"
X = 0
Y = 9999
Style = "label"
FieldType = "label"
`
	return taps.NewPanel(doc, styleMatrix, "")
}

// -------------------------------------------------
type ListApp struct {
	panel *taps.Panel
}

// -------------------------------------------------
//
//	Paging
//
// -------------------------------------------------
var isLast bool

func (m *ListApp) firstPage(common *Common) bool {
	return common.from == 1
}

func (m *ListApp) lastPage() bool {
	return isLast
}

func (m *ListApp) nextPage(common *Common) {
	if !(m.lastPage()) {
		common.from += common.rows
	}
}

func (m *ListApp) priorPage(common *Common) {
	if !(m.firstPage(common)) {
		common.from -= common.rows
	}
	if common.from < 1 {
		common.from = 1
	}
}

func (m *ListApp) isDataExist(common *Common) bool {
	return common.tableName != ""
}

func getFieldsInfo(listdata []listdb.ListItem, common *Common) (int, []int) {
	maxLength := 0
	var fieldsLength []int
	var flen int
	for _, s := range listdata {
		flen = runewidth.StringWidth(strings.TrimSpace(s.Field01))
		if flen > maxLength {
			maxLength = flen
		}
		fieldsLength = append(fieldsLength, flen)
	}
	if maxLength > common.cols {
		maxLength = common.cols
	}
	return maxLength, fieldsLength
}

func (m *ListApp) createListData(common *Common) ([]string, int) {
	var listData []string

	manager := listdb.GetManager(common.databaseName)
	err := manager.Connect(common.databaseName, common.connectString)
	if err != nil {
		panic(err)
	}

	listdb := manager.SearchDB(common.tableName, common.category, common.search, common.from, common.rows+1)
	recordCount := manager.GetRecordCount(common.tableName, common.category, common.search)
	manager.Close()

	data := listdb.GetListData()
	if len(data) <= common.rows {
		isLast = true
	} else {
		isLast = false
	}

	maxLength, fieldsLength := getFieldsInfo(data, common)
	var fld string
	var flen int
	for i, s := range data {
		if i >= common.rows {
			break
		}
		if maxLength < fieldsLength[i] {
			listData = append(listData, strings.TrimSpace(s.Field01)+s.Field02)
		} else {
			fld = strings.TrimSpace(s.Field01) + strings.Repeat(" ", (maxLength-fieldsLength[i]+1)) + strings.TrimSpace(s.Field02)
			flen = runewidth.StringWidth(fld)
			if flen < common.cols {
				fld = fld + strings.Repeat(" ", common.cols-flen)
			}
			listData = append(listData, fld)
		}
	}
	
	//log.Printf("listData[0]:%s listData[1]:%s recordCount:%d\n", listData[0], listData[1], recordCount)

	return listData, recordCount
}

// ---------------------------------------
// Main
// ---------------------------------------
func (m *ListApp) doFormat(common *Common) {
	if m.isDataExist(common) {
		listData, recordCount := m.createListData(common)
		m.panel.StoreList(listData, "LIST")
		m.panel.Store(fmt.Sprintf("%d/%d ", common.from, recordCount)+common.tableName, "L01")

		m.panel.GetDataField("C").Enabled()
		m.panel.GetDataField("S").Enabled()
	}else{
		m.panel.GetDataField("C").Disabled()
		m.panel.GetDataField("S").Disabled()
	}
	if m.lastPage() || !(m.isDataExist(common)){
		m.panel.GetDataField("N").Disabled()
	} else {
		m.panel.GetDataField("N").Enabled()
	}
	if m.firstPage(common) {
		m.panel.GetDataField("P").Disabled()
	} else {
		m.panel.GetDataField("P").Enabled()
	}
	m.panel.Say()
}

func (m *ListApp) Run(common *Common) {
	m.panel = ListPanel()
	_, common.rows = m.panel.GetListCount("LIST")
	table := &TableApp{}
	search := &SearchApp{}
	category := &CategoryApp{}
	detail := &DetailApp{}

	for {
		m.doFormat(common)
		k, n := m.panel.Read()
		if k == tcell.KeyEscape {
			break
		}

		if n == "T" {
			rt := table.Run(common)
			if rt != ""{
				common.reset()
				common.tableName = rt
			}
		}
		if n == "N" {
			m.nextPage(common)
		}
		if n == "P" {
			m.priorPage(common)
		}
		if n == "C" {
			rc := category.Run(common)
			if rc != ""{
				common.resetPaging()
				common.category = rc
				if common.category == "Clear" {
					common.category = ""
				}
			}
		}
		if n == "S" {
			rs := search.Run()
			common.resetPaging()
			common.search = rs
		}
		if n == "Q" {
			break
		}

		if k == tcell.KeyEnter && len(n) > 4 && n[:4] == "LIST" {
			num := m.panel.GetListFocus(n)
			if num >=0{
				common.selectedItem = common.from + num
				rs := detail.Run(common)
				if rs == "Q" {
					break
				}
			}
		}

	}
}

