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

// forageCmd represents the hunt command
var forageCmd = &cobra.Command{
	Use:   "forage",
	Short: "Checks a list of resource names for existence within Azure.",
	Long: `
Feed Forage a list of resource names on stdin, and it will enumerate known
Azure domains for the existence of those resources.

Example:
	$ cat resource-names.txt | steve forage

	If you want to enumerate permutations, add craft in:
	$ cat resource-names.txt | steve craft | steve forage
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

		// no stdin? no worky.
		if !hasStdin() {
			logger.Fatalln("input to be supplied through stdin")
		}

		gather := make(chan []azrecon.Resource)
		tracker := make(chan empty)

		// threads for reading stdin and resolving domains
		for i := 0; i < threadCount; i++ {
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					results, err := azrecon.CheckResourceExists(
						scanner.Text(),
						cfg.Resolvers[rand.Intn(len(cfg.Resolvers))],
					)
					if err != nil {
						continue
					}
					gather <- results
				}
				var e empty
				tracker <- e
			}()
		}

		// gathering thread
		go func() {
			for resources := range gather {
				for _, resource := range resources {
					if unique(resource.Domain) {
						if verbose {
							fmt.Printf("%s [type: %s]\n", resource.Domain, resource.Type)
						} else {
							fmt.Println(resource.Domain)
						}
					}
				}
			}
			var e empty
			tracker <- e
		}()

		// wait for threads
		for i := 0; i < threadCount; i++ {
			<-tracker
		}
		close(gather)
		<-tracker
	},
}

func init() {
	rootCmd.AddCommand(forageCmd)

	forageCmd.Flags().IntP("threads", "t", 50, "Number of threads for enumerating domains.")
	forageCmd.Flags().BoolP("verbose", "v", false, "Include information in output.")
}
