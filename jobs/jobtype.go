package jobs

// JobType is an enum that defines the type of crackjob (WPA2, MD5, ...)
type JobType int

const (
	// UNDEFINED identifies crackjobs that have not been defined
	UNDEFINED JobType = 0

	// WPA2 identifies a WPA2 crackjob
	WPA2 JobType = 1

	// MD5 identifies an MD5 crackjob
	MD5 JobType = 2
)

func (t JobType) String() string {
	return [...]string{"<undefined>", "WPA2", "MD5"}[t]
}
