package gravity

import (
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/einsteinplatform/gravity/pkg/git"
	"github.com/einsteinplatform/gravity/pkg/terraform"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// StartServer clones the provided Git repositories and begins
// supervising the Terraform projects contained within them
func StartServer(repoURLs []string, cloneDirBase, githubToken string) {
	sigterm, quit := make(chan os.Signal, 2), make(chan struct{})
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		log.Debug().Msg("Signal interrupt received; cleaning up")
		close(quit)
	}()

	var wg sync.WaitGroup

	for _, repoURL := range repoURLs {
		go fanoutRepo(repoURL, cloneDirBase, githubToken, &wg, quit)
		wg.Add(1)
	}

	wg.Wait()
}

// fanoutRepo clones a git repo, finds its Terraform project directories,
// and starts goroutines for each one
func fanoutRepo(repoURL, cloneDirBase, githubToken string, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()

	log.Debug().Str("repoURL", repoURL).Msg("Parsing repo URL")
	u, err := url.Parse(repoURL)
	if err != nil {
		log.Error().Str("repoURL", repoURL).Msg("Parsing repo URL failed")
		return
	}

	cloneDirRepo := path.Join(u.Host, strings.TrimSuffix(u.Path, ".git"))
	cloneDirFull := path.Join(cloneDirBase, cloneDirRepo)

	if _, err := os.Stat(cloneDirFull); os.IsNotExist(err) {
		log.Debug().Str("repo", repoURL).Msg("Cloning")
		err = git.Clone(repoURL, cloneDirFull, githubToken)
		if err != nil {
			log.Error().Str("repo", repoURL).Msg("Cloning failed")
			return
		}
	} else {
		log.Debug().Str("repo", cloneDirRepo).Msg("Repo exits; pulling")
		err = git.Pull(cloneDirFull, githubToken)
		if err != nil {
			log.Error().Str("repo", cloneDirRepo).Msg("Pulling failed")
			return
		}
	}

	log.Debug().Str("repo", cloneDirRepo).Msg("Searching for Terraform backend directories")
	backendDirs, err := terraform.FindStatefulDirs(cloneDirFull)
	if err != nil {
		log.Error().Str("repo", cloneDirRepo).Msg("Backend dir search failed")
		return
	}

	log.Debug().Str("backendDirs", "["+strings.Join(backendDirs, " ")+"]").Msg("Found backend directories")
	for _, backendDir := range backendDirs {
		go superviseBackend(backendDir, wg, quit)
		wg.Add(1)
	}
}

// TODO fanoutBackend() -> superviseWorkspaces()

func superviseBackend(backendDir string, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()

	log.Debug().Str("backendDir", backendDir).Msg("Running Terraform init")
	_, err := terraform.Init(backendDir)
	if err != nil {
		log.Error().Str("backendDir", backendDir).Msg("Init failed")
		return
	}

planLoop:
	for {
		log.Debug().Str("backendDir", backendDir).Msg("Running Terraform plan")
		planOut, err := terraform.Plan(backendDir)
		if err != nil {
			log.Error().Str("backendDir", backendDir).Msg("Plan failed")
			return
		}

		if !strings.Contains(*planOut, "No changes. Infrastructure is up-to-date.") {
			log.Info().Str("backendDir", backendDir).Msg("Drift detected! Running Terraform apply")
			_, err := terraform.Apply(backendDir)
			if err != nil {
				log.Error().Str("backendDir", backendDir).Msg("Apply failed")
				return
			}
			log.Info().Str("backendDir", backendDir).Msg("Apply succeeded")
		} else {
			log.Debug().Str("backendDir", backendDir).Msg("No changes")
		}

		// TODO: need to randomize when the goroutines run, else API throttling will occur
		wake := makeSleepChan(viper.GetInt("default_period"))

		select {
		case <-quit:
			break planLoop
		case <-wake:
			continue planLoop
		}
	}

	log.Debug().Str("backendDir", backendDir).Msg("Ending backend supervision")
}

// makeSleepChan creates a channel that will close after the provided
// number of seconds
func makeSleepChan(seconds int) chan struct{} {
	wake := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(seconds) * time.Second)
		close(wake)
	}()
	return wake
}
