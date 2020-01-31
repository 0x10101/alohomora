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
	StartTimestamp       time.Time  `xml:"started"`
	EndTimestamp         time.Time  `xml:"stopped"`
	Charset              string     `xml:"charset"`
	Offset               *big.Int   `xml:"offset"`
	Length               uint       `xml:"passlen"`
	Jobsize              *big.Int   `xml:"jobsize"`
	FinishedRuns         *big.Int   `xml:"runs"`
	Success              bool       `xml:"success"`
	SuccessClientID      string     `xml:"client"`
	SuccessClientAddress net.Addr   `xml:"clientaddr"`
	AccessData           AccessData `xml:"access"`
	JobType              string     `xml:"type"`
	Target               string     `xml:"target"` // might be type in the future
	MaxClientsConnected  uint       `xml:"maxclients"`
	PasswordsTried       *big.Int   `xml:"tries"`
}

// The AccessData type wraps a generic username and password
// combination for reporting purposes.
type AccessData struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
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
