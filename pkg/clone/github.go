package clone

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// TODO: generalize to GitLab, etc

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

// Pull pulls the provided Git directory, assuming it is a GitHub repo
func Pull(dir, githubToken string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: "dummy", // anything except an empty string
			Password: githubToken,
		},
		SingleBranch: true,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}
