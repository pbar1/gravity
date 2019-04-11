package gravity

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Enables Viper remote key/value store
	_ "github.com/spf13/viper/remote"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gravity",
	Short: "Einstein Platform developer multitool",
	Long:  `Einstein Platform developer multitool`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads config from Consul
func initConfig() {
	viper.SetDefault("consuladdr", "localhost:8500")
	viper.BindEnv("consuladdr", "GRAVITY_CONSUL_ADDR")

	consulAddr := viper.GetString("consuladdr")
	viper.AddRemoteProvider("consul", consulAddr, "gravity/config")
	viper.SetConfigType("hcl") // Need to explicitly set this to json

	err := viper.ReadRemoteConfig()
	if err != nil {
		log.Fatal(err)
	}
}
