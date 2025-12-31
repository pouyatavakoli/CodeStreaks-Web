package domain

import "time"

type Submission struct {
	CreationTimeSeconds int64  `json:"creationTimeSeconds" db:"creationTimeSeconds"`
	Verdict             string `json:"verdict" db:"verdict"`
}

func (s *Submission) GetSubmissionDate() time.Time {
	utcTime := time.Unix(s.CreationTimeSeconds, 0)

	IranTime := 3*time.Hour + 30*time.Minute
	return utcTime.Add(IranTime)
}
