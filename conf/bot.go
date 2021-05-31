package conf

type GitHubBot struct {
	GitHubActor string `json:"gitHubActor,omitempty"`
	GitHubEmail string `json:"gitHubEmail,omitempty"`
	GitHubToken string `json:"gitHubToken,omitempty"`
}
