package tui
import (
	"github.com/rsn604/taps"
	//"fmt"
	//"log"
	"os"
)

func App() {
	/*	
	file, err := os.OpenFile("TestList.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	*/
	common := NewCommon()
	common.reset()
	if len(os.Args) == 3 {
		common.databaseName = os.Args[1]
		common.connectString = os.Args[2]
	} else {
		common.databaseName = "BOLT"
		common.connectString = "./db/ListDB.boltdb"
	}
	common.cols, _ = taps.GetWindowSize()
	m := &ListApp{}
	m.Run(common)
}
