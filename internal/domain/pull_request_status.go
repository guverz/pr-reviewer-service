package domain

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

func (s PullRequestStatus) IsValid() bool {
	switch s {
	case PullRequestStatusOpen, PullRequestStatusMerged:
		return true
	default:
		return false
	}
}



