package gravity

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fooCmd represents the foo command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Retrieves values from Consul's key/value store",
	Long:  `Retrieves values from Consul's key/value store`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fooCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fooCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
