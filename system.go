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
	TableID   string        `xml:"table_id,attr"`
	TableName string        `xml:"table_name,attr"`
	TreeItem  []*nodeEpTree `xml:"TreeItem"`
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

	if d.Decode(&e) != nil {
		GetConfig().Err.Fatalf("Can`t read tree file from %v err %v", path, err)
	}
}

func GetEpTree() *epTree {
	_onceTreeXML.Do(func() {
		_treeXML = new(epTree)
		_treeXML.loadTree(GetConfig().PathTreeXML)
	})
	return _treeXML
}

func GetHeaders() *Headers {
	_onceHeaders.Do(func() {
		_headers = new(Headers)
		_headers.Fetch()
	})
	return _headers
}

func GetTablesMeta() *TablesMeta {
	_onceTablesMeta.Do(func() {
		_tablesMata = new(TablesMeta)
		_tablesMata.Fetch()
	})
	return _tablesMata
}

