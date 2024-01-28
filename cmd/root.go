/*
Released under YOLO licence. Idgaf what you do.
*/
package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

type config struct {
	Resolvers []string `yaml:"resolvers"`
}

var (
	// cfg is the config for use by all subcommands
	cfg = config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "steve",
	Short: "Azure resource attack surface mapping.",
	Long: `
  ⇱_ [◨_◧] /

  Steve - by @sam_phisher

  Steve can mine for Azure resources from a list of domains given to him,
  craft a list of likely resource names or hashcat-style masks for resource
  discovery, and test a list of resource names for presence across known
  Azure domains.

  Configure your DNS resolvers in:
    Linux: ~/.Steve/config.yml
    Windows: %userprofile%\Steve\config.yml
`,
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
	// determine user's home dir
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal(err)
	}

	// concat home dir and steve's config dir
	var configDirectory string
	if runtime.GOOS == "windows" {
		configDirectory = filepath.Join(homeDirectory, "Steve")
	} else {
		configDirectory = filepath.Join(homeDirectory, ".Steve")
	}

	// create dir if it doesn't exist
	err = os.MkdirAll(configDirectory, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}
	configFile := filepath.Join(configDirectory, "config.yaml")

	// check config file exists.
	if _, err := os.Stat(configFile); err != nil {
		// if not, add azrecon.Resolvers to cfg and marshal it to the
		// file, then continue.
		err = writeDefaultConfig(configFile)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// unmarshal config to cfg and continue.
	err = readConfigFromFile(configFile)
	if err != nil {
		logger.Fatal(err)
	}
}
