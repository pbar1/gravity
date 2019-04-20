package gravity

import (
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/einsteinplatform/gravity/pkg/git"
	"github.com/einsteinplatform/gravity/pkg/terraform"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// StartServer clones the provided Git repositories and begins
// supervising the Terraform projects contained within them
func StartServer(repoURLs []string, cloneDirBase, githubToken string) {
	var wg sync.WaitGroup

	for _, repoURL := range repoURLs {
		go superviseRepo(repoURL, cloneDirBase, githubToken, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func superviseRepo(repoURL, cloneDirBase, githubToken string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Debug().Str("repoURL", repoURL).Msg("Parsing repo URL")
	u, err := url.Parse(repoURL)
	if err != nil {
		log.Fatal().Str("repoURL", repoURL).Msg("Parsing repo URL failed")
	}

	cloneDirRepo := path.Join(u.Host, strings.TrimSuffix(u.Path, ".git"))
	cloneDirFull := path.Join(cloneDirBase, cloneDirRepo)

	if _, err := os.Stat(cloneDirFull); os.IsNotExist(err) {
		log.Debug().Str("repo", repoURL).Msg("Cloning")
		err = git.Clone(repoURL, cloneDirFull, githubToken)
		if err != nil {
			log.Fatal().Str("repo", repoURL).Msg("Cloning failed")
		}
	} else {
		log.Debug().Str("repo", cloneDirRepo).Msg("Repo exits; pulling")
		err = git.Pull(cloneDirFull, githubToken)
		if err != nil {
			log.Fatal().Str("repo", cloneDirRepo).Msg("Pulling failed")
		}
	}

	log.Debug().Str("repo", cloneDirRepo).Msg("Searching for Terraform backend directories")
	backendDirs, err := terraform.FindStatefulDirs(cloneDirFull)
	if err != nil {
		log.Fatal().Str("repo", cloneDirRepo).Msg("Backend dir search failed")
	}

	log.Debug().Str("backendDirs", strings.Join(backendDirs, "")).Msg("Found backend directories")
	for _, backendDir := range backendDirs {
		go superviseBackend(backendDir, wg)
		wg.Add(1)
	}
}

func superviseBackend(backendDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Debug().Str("backendDir", backendDir).Msg("Running Terraform init")
	_, err := terraform.Init(backendDir)
	if err != nil {
		log.Fatal().Str("backendDir", backendDir).Msg("Init failed")
	}

	for {
		log.Debug().Str("backendDir", backendDir).Msg("Running Terraform plan")
		planOut, err := terraform.Plan(backendDir)
		if err != nil {
			log.Error().Str("backendDir", backendDir).Msg("Plan failed")
			continue
		}

		if !strings.Contains(*planOut, "No changes. Infrastructure is up-to-date.") {
			log.Info().Str("backendDir", backendDir).Msg("Drift detected! Running Terraform apply")
			_, err := terraform.Apply(backendDir)
			if err != nil {
				log.Fatal().Str("backendDir", backendDir).Msg("Apply failed")
			}
			log.Info().Str("backendDir", backendDir).Msg("Apply succeeded")
		} else {
			log.Debug().Str("backendDir", backendDir).Msg("No changes")
		}

		time.Sleep(time.Duration(viper.GetInt("default_period")) * time.Second)
	}
}

// TODO: if signal ctrl-c received, graceful shutdown -> break
