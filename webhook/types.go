package webhook

// DockerHubPayload is the data structure which holds
// the form in which request bodies are expected to be in
// when Docker hub invokes the webhook from a build trigger.
type DockerHubPayload struct {
	PushData    *DockerHubPushData   `json:"push_data"`
	CallbackURL string               `json:"callback_url"`
	Repository  *DockerHubRepository `json:"repository"`
}

// DockerHubPushData holds information specific to the push event
// from a Docker hub trigger.
type DockerHubPushData struct {
	PushedAt int64    `json:"pushed_at"`
	Images   []string `json:"images"`
	Tag      string   `json:"tag"`
	Pusher   string   `json:"pusher"`
}

// DockerHubRepository provides that data structure for all information about the docker
// repository the build occured for.
type DockerHubRepository struct {
	CommentCount    string  `json:"comment_count"`
	DateCreated     float64 `json:"date_created"`
	Description     string  `json:"description"`
	Dockerfile      string  `json:"dockerfile"`
	FullDescription string  `json:"full_description"`
	IsOfficial      bool    `json:"is_official"`
	IsPrivate       bool    `json:"is_private"`
	IsTrusted       bool    `json:"is_trusted"`
	Name            string  `json:"name"`
	Namespace       string  `json:"namespace"`
	Owner           string  `json:"owner"`
	RepoName        string  `json:"repo_name"`
	RepoURL         string  `json:"repo_url"`
	StarCount       int64   `json:"star_count"`
	Status          string  `json:"status"`
}
