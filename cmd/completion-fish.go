/*
Copyright © 2020 Iggy Jackson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// completionFishCmd represents the fish command
var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish shell completion code",
	Long:  `Output (to stdout) shell completion for the fish shell`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenFishCompletion(os.Stdout, true)
		if err != nil {
			log.Fatalln("failed to generate fish completion", err)
		}
	},
}

func init() {
	completionCmd.AddCommand(completionFishCmd)
}
