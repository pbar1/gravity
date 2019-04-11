package gravity

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fooCmd represents the foo command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays a list of keys that can be read from Consul",
	Long:  `Displays a list of keys that can be read from Consul`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.Get("port"))     // 8080
		fmt.Println(viper.Get("hostname")) // myhostname.com
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fooCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fooCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
