package entity

import "encoding/xml"

type (
	Position struct {
		Key   string `xml:"Key"`
		Value string `xml:"Value"`
	}
	Entry struct {
		Positions []Position `xml:"String"`
	}
	Group struct {
		Name   string  `xml:"Name"`
		Notes  string  `xml:"Notes"`
		Entrys []Entry `xml:"Entry"`
		Groups []Group `xml:"Group"`
	}
	KeePassFile struct {
		XMLName xml.Name `xml:"KeePassFile"`
		Meta    struct {
			Generator           string `xml:"Generator"`
			DatabaseName        string `xml:"DatabaseName"`
			DatabaseDescription string `xml:"DatabaseDescription"`
		} `xml:"Meta"`
		Root struct {
			Group Group `xml:"Group"`
		} `xml:"Root"`
	}
)
