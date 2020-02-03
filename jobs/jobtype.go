package jobs

// JobType is an enum that defines the type of crackjob (WPA2, MD5, ...)
type JobType int

const (
	// WPA2 identifies a WPA2 crackjob
	WPA2 JobType = 1
)
