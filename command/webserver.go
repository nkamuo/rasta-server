package command

import (
	"fmt"

	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/web"
	"github.com/spf13/cobra"
)

func StartWebServer(config web.WebServerConfig) (err error) {

	sysConfig, err := initializers.LoadConfig()
	if err != nil {
		return err
	}

	if config.AssetDir == "" {
		config.AssetDir = sysConfig.ASSET_DIR
	} else {
		sysConfig.ASSET_DIR = config.AssetDir
	}
	if config.PublicPrefix == "" {
		config.PublicPrefix = sysConfig.PUBLIC_PREFIX
	} else {
		sysConfig.PUBLIC_PREFIX = config.PublicPrefix
	}
	if config.Addr == "" {
		config.Addr = sysConfig.SERVER_ADDRESS
	} else {
		sysConfig.SERVER_ADDRESS = config.Addr
	}
	if config.Port == 0 {
		if sysConfig.SERVER_PORT != nil {
			config.Port = *sysConfig.SERVER_PORT
		}
	} else {
		sysConfig.SERVER_PORT = &config.Port
	}

	r, err := web.BuildWebServer(config)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	err = r.Run(addr)
	return
}

func buildWebServerCommand(command *cobra.Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// username, _ := rootCmd.Flags().GetString("username")
		// password, _ := rootCmd.Flags().GetString("password")

		address, _ := command.Flags().GetString("address")
		port, _ := command.Flags().GetUint("port")
		publicPrefix, _ := command.Flags().GetString("public-prefix")
		htdocs, _ := command.Flags().GetString("htdocs")
		index, _ := command.Flags().GetString("index")
		directoryListing, _ := command.Flags().GetBool("directory-listing")

		// addr := fmt.Sprintf("%s:%s", config.Addr, config.Port)

		// sysConfig.SERVER_ADDRESS = addr

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

	config, err := initializers.LoadConfig()
	if err != nil {
		fmt.Println("CONFIG ERROR:", err)
		return
	}

	webserverCmd.Run = buildWebServerCommand(webserverCmd)
	SERVER_ADDRESS := config.SERVER_ADDRESS
	SERVER_PORT := config.SERVER_PORT
	if SERVER_PORT == nil {
		P := uint(8090)
		SERVER_PORT = &P
	}

	rootCmd.AddCommand(webserverCmd)
	webserverCmd.PersistentFlags().StringP("address", "a", SERVER_ADDRESS, "the ip address to bind the server to")
	webserverCmd.PersistentFlags().UintP("port", "p", *SERVER_PORT, "the port to bind the request to")

	webserverCmd.PersistentFlags().StringP("htdocs", "", "", "The assets folder for HTML, CSS, JS files")
	webserverCmd.PersistentFlags().StringP("public-prefix", "", "", "The prefix that will be served from the htdocs folder")
	webserverCmd.PersistentFlags().StringP("index", "i", "", "The path to the index file")
	webserverCmd.PersistentFlags().BoolP("directory-listing", "l", true, "Allow Directory Listing")
}
