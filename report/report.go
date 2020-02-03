package report

import (
	"encoding/xml"
	"math/big"
	"net"
	"time"
)

// The Report type wraps everything that the server can
// report to the user.
type Report struct {
	XMLName              xml.Name
	StartTimestamp       time.Time  `xml:"started" json:"started"`
	EndTimestamp         time.Time  `xml:"stopped" json:"stopped"`
	Charset              string     `xml:"charset" json:"charset"`
	Offset               *big.Int   `xml:"offset"  json:"offset"`
	Length               uint       `xml:"passlen" json:"passlen"`
	Jobsize              *big.Int   `xml:"jobsize" json:"jobsize"`
	FinishedRuns         *big.Int   `xml:"runs"    json:"run"`
	Success              bool       `xml:"success" json:"success"`
	SuccessClientID      string     `xml:"client"  json:"client"`
	SuccessClientAddress net.Addr   `xml:"clientaddr" json:"clientaddr"`
	AccessData           AccessData `xml:"access" json:"access"`
	JobType              string     `xml:"type" json:"type"`
	Target               string     `xml:"target" json:"target"` // might be type in the future
	MaxClientsConnected  uint       `xml:"maxclients" json:"maxclients"`
	PasswordsTried       *big.Int   `xml:"tries" json:"tries"`
}

// The AccessData type wraps a generic username and password
// combination for reporting purposes.
type AccessData struct {
	Username string `xml:"username" json:"username"`
	Password string `xml:"password" json:"password"`
}

// New initializes a new empty report instance
func New() *Report {
	report := &Report{
		XMLName:              xml.Name{Local: "report"},
		StartTimestamp:       time.Now(),
		EndTimestamp:         time.Time{},
		Charset:              "",
		Offset:               big.NewInt(0),
		Length:               0,
		Jobsize:              big.NewInt(0),
		FinishedRuns:         big.NewInt(0),
		Success:              false,
		SuccessClientID:      "",
		SuccessClientAddress: nil,
		AccessData:           AccessData{},
		JobType:              "",
		Target:               "",
		MaxClientsConnected:  0,
		PasswordsTried:       big.NewInt(0),
	}

	report.XMLName = xml.Name{Local: "report"}
	return report
}
