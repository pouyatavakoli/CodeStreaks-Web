package domain

import "time"

type Submission struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	UserID                 uint      `gorm:"index;not null" json:"user_id"`
	CodeforcesSubmissionID int64     `gorm:"uniqueIndex;not null" json:"codeforces_submission_id"`
	//ProblemName            string    `json:"problem_name"`
	//ContestID              int       `json:"contest_id"`
	//ProblemIndex           string    `json:"problem_index"`
	Verdict                string    `json:"verdict"`
	//ProgrammingLanguage    string    `json:"programming_language"`
	SubmittedAt            time.Time `gorm:"index;not null" json:"submitted_at"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// CodeforcesSubmission represents the API response structure
type CodeforcesSubmission struct {
	ID                  int               `json:"id"`
	//ContestID           int               `json:"contestId"`
	CreationTimeSeconds int64             `json:"creationTimeSeconds"`
	//ProblemsetName      string            `json:"problemsetName"`
	//Problem             CodeforcesProblem `json:"problem"`
	//Author              CodeforcesAuthor  `json:"author"`
	//ProgrammingLanguage string            `json:"programmingLanguage"`
	Verdict             string            `json:"verdict"`
}

// type CodeforcesProblem struct {
// 	ContestID int      `json:"contestId"`
// 	Index     string   `json:"index"`
// 	Name      string   `json:"name"`
// 	Type      string   `json:"type"`
// 	Rating    int      `json:"rating"`
// 	Tags      []string `json:"tags"`
// }

// type CodeforcesAuthor struct {
// 	ContestID        int      `json:"contestId"`
// 	Members          []Member `json:"members"`
// 	ParticipantType  string   `json:"participantType"`
// 	Ghost            bool     `json:"ghost"`
// 	StartTimeSeconds int64    `json:"startTimeSeconds"`
// }

type Member struct {
	Handle string `json:"handle"`
}

// CodeforcesUserInfo represents user info from API
type CodeforcesUserInfo struct {
	Handle    string `json:"handle"`
	Rating    int    `json:"rating"`
	MaxRating int    `json:"maxRating"`
	Rank      string `json:"rank"`
	Country   string `json:"country"`
}

type CodeforcesAPIResponse struct {
	Status string                 `json:"status"`
	Result []CodeforcesSubmission `json:"result"`
}

type CodeforcesUserAPIResponse struct {
	Status string               `json:"status"`
	Result []CodeforcesUserInfo `json:"result"`
}
