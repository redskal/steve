/*
Released under YOLO licence. Idgaf what you do.
*/
package cmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/redskal/steve/pkg/azrecon"
	"github.com/spf13/cobra"
)

// mineCmd represents the mine command
var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Given a list of domains, mine will check for Azure resources.",
	Long: `
Mine will read a list of domains from stdin and check each of them
for CNAMEs pointing to known Azure domains, indicating an Azure resource.
Optionally, provide -c to check for potential subdomain takeovers.

Example:
	# identify Azure resources from CNAMEs
	$ subfinder -d microsoft.com | steve mine

	# identify Azure resources with potential subdomain takeovers
	$ subfinder -d microsoft.com | steve mine -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		// get flags
		threadCount, err := cmd.Flags().GetInt("threads")
		if err != nil {
			logger.Fatal(err)
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			logger.Fatal(err)
		}
		checkTakeover, err := cmd.Flags().GetBool("check-takeover")
		if err != nil {
			logger.Fatal(err)
		}

		// no stdin? no worky.
		if !hasStdin() {
			logger.Fatalln("input to be supplied through stdin")
		}

		gather := make(chan azrecon.Domain)
		tracker := make(chan empty)

		// threads for pulling from stdin and checking CNAMEs
		for i := 0; i < threadCount; i++ {
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					domainResult, err := azrecon.CheckAzureCnames(
						scanner.Text(),
						cfg.Resolvers[rand.Intn(len(cfg.Resolvers))],
						checkTakeover,
					)
					if err != nil {
						continue
					}
					gather <- domainResult
				}
				var e empty
				tracker <- e
			}()
		}

		// gathering thread
		go func() {
			for result := range gather {
				for _, cname := range result.Cnames {
					if unique(cname.Cname) {
						var takeoverString string
						if checkTakeover && verbose {
							takeoverString = fmt.Sprintf("[takeover: %v]", cname.Takeover)
						} else if checkTakeover && cname.Takeover {
							takeoverString = " (Potential subdomain takeover)"
						}

						if verbose {
							fmt.Printf("%s [src: %s] [type: %s] %s\n", cname.Cname, result.Domain, cname.Type, takeoverString)
						} else {
							fmt.Printf("%s%s\n", cname.Cname, takeoverString)
						}
					}
				}
			}
			var e empty
			tracker <- e
		}()

		// clean up on isle threads
		for i := 0; i < threadCount; i++ {
			<-tracker
		}
		close(gather)
		<-tracker
	},
}

func init() {
	rootCmd.AddCommand(mineCmd)

	mineCmd.Flags().BoolP("check-takeover", "c", false, "Check for subdomain takeover.")
	mineCmd.Flags().IntP("threads", "t", 50, "Number of threads for querying CNAME records.")
	mineCmd.Flags().BoolP("verbose", "v", false, "Include information in output.")
}
