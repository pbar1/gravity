package gravity

import (
	"os"
	"strings"
	"time"

	"github.com/einsteinplatform/gravity/pkg/clone"
	"github.com/einsteinplatform/gravity/pkg/gravity"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Gravity server",
	Long: `Starts the Gravity server, which will monitor Git repos it subscribes
to for Terraform code. It will attempt to run a plan and notify if there
are any changes from the desired state.`,
	Run: func(cmd *cobra.Command, args []string) {
		repos := viper.GetStringSlice("repos")
		githubToken := viper.GetString("github_token")
		cloneDir := viper.GetString("clone_dir")

		// TODO: clone into subdirs, don't remove dirs
		for _, repo := range repos {
			log.Debug().Str("repo", repo).Msg("Cloning")
			err := clone.Clone(repo, cloneDir, githubToken)
			if err != nil {
				log.Fatal().Str("repo", repo).Msg("Cloning failed")
			}
			defer os.RemoveAll(cloneDir)

			log.Debug().Str("repo", repo).Msg("Running Terraform init")
			_, err = gravity.Init(cloneDir)
			if err != nil {
				log.Fatal().Str("repo", repo).Msg("Init failed")
			}

			for {
				log.Debug().Str("repo", repo).Msg("Running Terraform plan")
				planOut, err := gravity.Plan(cloneDir)
				if err != nil {
					log.Fatal().Str("repo", repo).Msg("Plan failed")
				}

				if !strings.Contains(*planOut, "No changes. Infrastructure is up-to-date.") {
					log.Info().Str("repo", repo).Msg("Plan has detected drift!! Attempting to rectify")
					_, err := gravity.Apply(cloneDir)
					if err != nil {
						log.Fatal().Str("repo", repo).Msg("Apply failed")
					}
					log.Info().Str("repo", repo).Msg("Apply succeeded")
				} else {
					log.Debug().Str("repo", repo).Msg("No changes")
				}

				time.Sleep(10 * time.Second)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
