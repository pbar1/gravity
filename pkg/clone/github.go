package clone

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Clone clones the provided GitHub repository into the provided directory
func Clone(url, dir, githubToken string) error {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "dummy", // anything except an empty string
			Password: githubToken,
		},
		SingleBranch: true,
	})
	return err
}
