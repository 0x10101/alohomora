package report

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/gosuri/uitable"
	"github.com/steps0x29a/alohomora/term"
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

func fmtHeader(val string) string {
	return term.Bold(term.White(val))
}

func fmtValue(val string) string {
	return term.BrightBlue(val)
}

func reportLine(header, value string) {
	fmt.Printf("%s\t\t%s\n", fmtHeader(header), fmtValue(value))
}

func (report *Report) Print() {
	fmt.Println()

	t := uitable.New()
	t.MaxColWidth = 80
	t.Wrap = true // wrap columns
	t.AddRow("Server started", fmtValue(report.StartTimestamp.Format("02.01.2006 15:04:05.999")))
	t.AddRow("Server stopped", fmtValue(report.EndTimestamp.Format("02.01.2006 15:04:05.999")))
	t.AddRow("Charset", fmtValue(report.Charset))
	t.AddRow("Offset", fmtValue(report.Offset.String()))
	t.AddRow("Password length", fmtValue(fmt.Sprintf("%d", report.Length)))
	t.AddRow("Jobsize", fmtValue(report.Jobsize.String()))
	t.AddRow("Finished runs", fmtValue(report.FinishedRuns.String()))
	t.AddRow("Job type", fmtValue(report.JobType))
	t.AddRow("Target", fmtValue(report.Target))
	t.AddRow("Passwords tried (approx.)", fmtValue(report.PasswordsTried.String()))
	t.AddRow("Max. clients connected", fmtValue(fmt.Sprintf("%d", report.MaxClientsConnected)))
	if report.Success {
		t.AddRow("Password cracked", term.Bold(term.BrightGreen("YES")))
		t.AddRow("Username", term.Bold(term.BrightGreen(report.AccessData.Username)))
		t.AddRow("Password", term.Bold(term.BrightGreen(report.AccessData.Password)))
		t.AddRow("Successful client", fmtValue(fmt.Sprintf("%s (%s)", report.SuccessClientID, report.SuccessClientAddress)))

	} else {
		t.AddRow("Password cracked", term.Bold(term.BrightRed("NO")))
	}

	/*
		reportLine("Server started:", report.StartTimestamp.String())
		reportLine("Server stopped:", report.EndTimestamp.String())
		reportLine("Charset used:", report.Charset)
		reportLine("Offset used:", report.Offset.String())
		reportLine("Password len:", fmt.Sprintf("%d", report.Length))
		reportLine("Jobsize used:", report.Jobsize.String())
		reportLine("Finished runs:", report.FinishedRuns.String())
		reportLine("Type of job:", report.JobType)
		reportLine("Target used:", report.Target)
		reportLine("Overall tries:", report.PasswordsTried.String())
		reportLine("Max clients:", fmt.Sprintf("%d", report.MaxClientsConnected))
		fmt.Printf("%s\t\t", fmtHeader("Password found:"))
		if report.Success {
			fmt.Printf("%s\n", term.Bold(term.BrightGreen("YES")))
			fmt.Printf("%s\t\t%s\n", fmtHeader("Username:"), term.Bold(term.BrightGreen(report.AccessData.Username)))
			fmt.Printf("%s\t\t%s\n", fmtHeader("Password:"), term.Bold(term.BrightGreen(report.AccessData.Password)))
			reportLine("Winner client:", fmt.Sprintf("%s (%s)", report.SuccessClientID, report.SuccessClientAddress))
		} else {
			fmt.Printf("%s\n", term.Bold(term.BrightRed("NO")))
		}*/

	fmt.Println(t)

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
