package gravity

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/einsteinplatform/gravity/internal/app/gravity"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Gravity server",
	Long: `Starts the Gravity server, which will monitor Git repos it subscribes
to for Terraform code. It will attempt to run a plan and notify if there
are any changes from the desired state.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoURLs := viper.GetStringSlice("repos")
		githubToken := viper.GetString("github_token")
		cloneDirBase := viper.GetString("clone_dir")

		gravity.StartServer(repoURLs, cloneDirBase, githubToken)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
