package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"news-weeder/internal/config"
)

var serviceConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "news-weeder",
	Short: "News Weeder server",
	Long: `There is news-weeder project which has been release for filtering semantic similar news by 
			comparing embeddings of source article content text data.`,

	Run: func(cmd *cobra.Command, _ []string) {
		fromEnv, _ := cmd.Flags().GetBool("from-env")

		var parseErr error
		if fromEnv {
			disabledDotenv, _ := cmd.Flags().GetBool("with-dotenv")
			serviceConfig, parseErr = config.LoadEnv(disabledDotenv)
		} else {
			filePath, _ := cmd.Flags().GetString("config")
			serviceConfig, parseErr = config.FromFile(filePath)
		}

		if parseErr != nil {
			log.Fatal(parseErr)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() *config.Config {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return serviceConfig
}

func init() {
	flags := rootCmd.Flags()
	flags.StringP("config", "c", "./configs/development.toml", "Parse options from config file.")
	flags.BoolP("from-env", "e", false, "Parse options from env.")
	flags.BoolP("with-dotenv", "j", false, "Parse options from existing .env file.")
}
