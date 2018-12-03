package main

import (
	"encoding/json"
	"encoding/xml"
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
	_treeXML        *epTree
	_onceTreeXML    sync.Once
	_onceHeaders    sync.Once
	_headers        *Headers
	_onceTablesMeta sync.Once
	_tablesMata     *TablesMeta
	_onceEmptyText  sync.Once
	_emptyText      *EmptyText
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

type nodeEpTree struct {
	Name      string        `xml:"name,attr"`
	TableID   string        `xml:"table_id,attr" json:",omitempty"`
	TableName string        `xml:"table_name,attr"  json:",omitempty"`
	TreeItem  []*nodeEpTree `xml:"TreeItem"  json:",omitempty"`
}

type epTree struct {
	TreeItem []*nodeEpTree `xml:"TreeItem"`
}

func (e *epTree) loadTree(path string) {
	file, err := os.Open(path)
	if err != nil {
		GetConfig().Err.Fatalf("Can`t read tree file from %v err %v", path, err)
	}

	d := xml.NewDecoder(file)

	if err := d.Decode(&e); err != nil {
		GetConfig().Err.Fatalf("Can`t read tree file from %v err %v", path, err)
	}
}

//GetEpTree получение дерева таблиц
func GetEpTree() *epTree {
	_onceTreeXML.Do(func() {
		_treeXML = new(epTree)
		_treeXML.loadTree(GetConfig().PathTreeXML)
	})
	return _treeXML
}

//GetHeaders получение всех заголовок всех таблиц singletone
func GetHeaders() *Headers {
	_onceHeaders.Do(func() {
		_headers = new(Headers)
		if err := _headers.Fetch(); err != nil {
			log.Fatal(err)
		}
	})

	return _headers
}

func GetTablesMeta() *TablesMeta {
	_onceTablesMeta.Do(func() {
		_tablesMata = new(TablesMeta)
		if err := _tablesMata.Fetch(); err != nil {
			log.Println(err)
		}
	})
	return _tablesMata
}

//
func GetEmptyText() *EmptyText{
	_onceEmptyText.Do(func() {
		_emptyText = new(EmptyText)
		if err := _emptyText.Fetch(); err != nil {
			log.Fatal(err)
		}
	})
	return _emptyText
}