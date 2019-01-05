package base

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type config struct {
	ConnStr     string `json:"connect_string"`
	Port        string `json:"port"`
	PathTreeXML string `json:"xml_file"`
	Err         *log.Logger
	Warn        *log.Logger
	Info        *log.Logger
}

var (
	_config         *config
	_onceConfig     sync.Once
	//_treeXML        *epTree
	//_onceTreeXML    sync.Once
	//_onceHeaders    sync.Once
	//_headers        *Headers
	//_onceTablesMeta sync.Once
	//_tablesMata     *map[int]TableInfo
	//_onceEmptyText  sync.Once
	//_emptyText      *map[int]map[int]string
)

//GetConfig получение объекта конфига
func GetConfig() *config {
	_onceConfig.Do(func() {
		_config = new(config)
		_config.load()
	})
	return _config
}

func (c *config) load() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	confFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	dc := json.NewDecoder(confFile)
	if err := dc.Decode(&c); err != nil {
		log.Fatal("Read config file: ", err)
	}

	if c.ConnStr == "" {
		log.Fatal("Can`t read connection string: ", c.ConnStr)
	}

	c.Err = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	c.Warn = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime)
	c.Info = log.New(os.Stderr, "[INFO] ", log.Ldate|log.Ltime)
}




//GetEpTree получение дерева таблиц
//func GetEpTree() *epTree {
//	_onceTreeXML.Do(func() {
//		_treeXML = new(epTree)
//		_treeXML.loadTree(GetConfig().PathTreeXML)
//	})
//	return _treeXML
//}
//
////GetHeaders получение всех заголовок всех таблиц singletone
//func GetHeaders() *Headers {
//	_onceHeaders.Do(func() {
//		_headers = new(Headers)
//
//		h, err := NewDatabase().GetHeaders()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		*_headers = *h
//
//	})
//
//	return _headers
//}
//
//func GetTablesMeta() *map[int]TableInfo {
//	_onceTablesMeta.Do(func() {
//
//		_tablesMata = new(map[int]TableInfo)
//
//		tm, err := NewDatabase().GetTablesInfo()
//		if err != nil {
//			return
//		}
//
//		*_tablesMata = tm
//	})
//	return _tablesMata
//}
//
////
//func GetEmptyText() *map[int]map[int]string {
//	_onceEmptyText.Do(func() {
//		_emptyText = new(map[int]map[int]string)
//
//		et, err := NewDatabase().GetTextForEmptyTable()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		*_emptyText = et
//	})
//	return _emptyText
//}
