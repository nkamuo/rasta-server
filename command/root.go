package command

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rasta",
	Short: `rasta is a CLI client of task control.`,
	Long:  `rasta is a CLI client of task control. you can use our task control portal also for managing your tasks`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("cwd", "d", "", "The Current working directory")
	rootCmd.PersistentFlags().StringP("env", "e", "", "the .env path relative to the Current working directory")
	// 	// Add flags in root command if required

	// 	// rootCmd.AddCommand(addCmd)
	// 	rootCmd.PersistentFlags().StringP("username", "u", "", "the username of git")
	// 	rootCmd.PersistentFlags().StringP("password", "p", "", "the access token of the git")
}
