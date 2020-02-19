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

	// SHA1 identifies a SHA1 crackjob
	SHA1 JobType = 3

	// SHA256 identifies a SHA256 crackjob
	SHA256 JobType = 4

	// SHA512 identifies a SHA512 crackjob
	SHA512 JobType = 5
)

func (t JobType) String() string {
	return [...]string{"<undefined>", "WPA2", "MD5", "SHA1", "SHA256", "SHA512"}[t]
}
