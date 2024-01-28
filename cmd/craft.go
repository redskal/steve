/*
Released under YOLO licence. Idgaf what you do.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/redskal/steve/pkg/azrecon"
	"github.com/spf13/cobra"
)

// craftCmd represents the craft command
var craftCmd = &cobra.Command{
	Use:   "craft",
	Short: "Craft all permutations or hashcat-style masks of a given list of resource names.",
	Long: `
Craft allows you to generate all permutations or hashcat-style masks of
a given list of Azure resource names.

Use this for further brute-force discovery with your favourite tools, or
with hashcat piped into dnsx.

Example:
	$ cat resource-names.txt | steve craft | steve forage

	or (if you've got time and several 4090s)

	$ for line in $(cat resource-names.txt | steve craft -m);
	  do
	    hashcat -m 0 -a 3 --stdout $line | steve forage;
	  done;
`,
	Run: func(cmd *cobra.Command, args []string) {
		// grab our argument values
		generateMasks, err := cmd.Flags().GetBool("mask")
		if err != nil {
			logger.Fatal(err)
		}

		// no stdin? no worky.
		if !hasStdin() {
			logger.Fatalln("input to be supplied through stdin")
		}

		var inputResults []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputResults = append(inputResults, scanner.Text())
		}

		// generate output
		var output []string
		if generateMasks {
			output = azrecon.GetMaskCombinations(inputResults)
		} else {
			output = azrecon.GetNameCombinations(inputResults)
		}
		if len(output) == 0 {
			logger.Fatalln("no results from azrecon package")
		}
		for _, line := range output {
			fmt.Println(line)
		}
	},
}

func init() {
	rootCmd.AddCommand(craftCmd)

	craftCmd.Flags().BoolP("mask", "m", false, "Produce hashcat-style masks instead of string permutations.")
}
