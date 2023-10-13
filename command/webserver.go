package command

import (
	"fmt"

	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/web"
	"github.com/spf13/cobra"
)

func StartWebServer(config web.WebServerConfig) (err error) {
	r, err := web.BuildWebServer(config)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%s", config.Addr, config.Port)
	err = r.Run(addr)
	return
}

func buildWebServerCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// username, _ := rootCmd.Flags().GetString("username")
		// password, _ := rootCmd.Flags().GetString("password")

		address, _ := command.Flags().GetString("address")
		port, _ := command.Flags().GetString("port")
		publicPrefix, _ := command.Flags().GetString("public-prefix")
		htdocs, _ := command.Flags().GetString("htdocs")
		index, _ := command.Flags().GetString("index")
		directoryListing, _ := command.Flags().GetBool("directory-listing")

		config := web.WebServerConfig{
			Addr:                  address,
			Port:                  port,
			AssetDir:              htdocs,
			IndexFile:             index,
			PublicPrefix:          publicPrefix,
			AllowDirectoryListing: directoryListing,
		}

		StartWebServer(config)

	}
}

// webserverCmd to start the rasta webserver
var webserverCmd = &cobra.Command{
	Use:   "serve",
	Short: "get repo details",
	Long:  `Get Repo information using the Cobra Command`,
}

func init() {
	// Add flags in root command if required

	config, err := initializers.LoadConfig(".")
	if err != nil {
		fmt.Println("CONFIG ERROR:", err)
		return
	}

	webserverCmd.Run = buildWebServerCommand(webserverCmd)
	SERVER_ADDRESS := config.SERVER_ADDRESS
	SERVER_PORT := config.SERVER_PORT
	if SERVER_PORT == "" {
		SERVER_PORT = "8090"
	}

	rootCmd.AddCommand(webserverCmd)
	webserverCmd.PersistentFlags().StringP("address", "a", SERVER_ADDRESS, "the ip address to bind the server to")
	webserverCmd.PersistentFlags().StringP("port", "p", SERVER_PORT, "the port to bind the request to")

	webserverCmd.PersistentFlags().StringP("htdocs", "", "", "The assets folder for HTML, CSS, JS files")
	webserverCmd.PersistentFlags().StringP("public-prefix", "", "", "The prefix that will be served from the htdocs folder")
	webserverCmd.PersistentFlags().StringP("index", "i", "", "The path to the index file")
	webserverCmd.PersistentFlags().BoolP("directory-listing", "l", true, "Allow Directory Listing")
}
