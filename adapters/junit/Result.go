package junit

import "encoding/xml"

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr"`
	Tests      string     `xml:"tests,attr"`
	Skipped    string     `xml:"skipped,attr"`
	Failures   string     `xml:"failures,attr"`
	Errors     string     `xml:"errors,attr"`
	TimeStamp  string     `xml:"timestamp,attr"`
	Hostname   string     `xml:"hostname,attr"`
	Time       string     `xml:"time,attr"`
	Properties Properties `xml:"properties"`
	TestCases  []TestCase `xml:"testcase"`
	SystemOut  string     `xml:"system-out"`
	SystemErr  string     `xml:"system-err"`
}

type Properties struct{}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Classname string   `xml:"classname,attr"`
	Time      string   `xml:"time,attr"`
	Name      string   `xml:"name,attr"`
	Failure   *string  `xml:"failure,omitempty"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
}

type SystemOut struct{}

type SystemErr struct{}
