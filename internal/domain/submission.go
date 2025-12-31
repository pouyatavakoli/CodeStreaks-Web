package domain

type Submission struct {
	CreationTimeSeconds int64  `json:"creationTimeSeconds" db:"creationTimeSeconds"`
	Verdict             string `json:"verdict" db:"verdict"`
}
