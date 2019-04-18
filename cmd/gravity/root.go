package gravity

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Enables Viper remote key/value store
	_ "github.com/spf13/viper/remote"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gravity",
	Short: "Terraform dynamic state-driver",
	Long:  `Terraform dynamic state-driver`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("%s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads config from Consul
func initConfig() {
	viper.SetDefault("clone_dir", ".gravity")
	viper.BindEnv("github_token", "GITHUB_TOKEN")

	viper.SetDefault("consul_addr", "localhost:8500")
	viper.BindEnv("consul_addr", "GRAVITY_CONSUL_ADDR")

	consulAddr := viper.GetString("consul_addr")
	consulConfigKey := "gravity/config"
	viper.AddRemoteProvider("consul", consulAddr, consulConfigKey)
	viper.SetConfigType("hcl")

	log.Debug().Str("consulAddr", consulAddr).Str("consulConfigKey", consulConfigKey).
		Msg("Reading config from Consul")

	err := viper.ReadRemoteConfig()
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
