/*
Copyright © 2020 Iggy Jackson <iggy@theiggy.com>

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
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// workspacesDescribeCmd represents the workspacesDescribe command
var workspacesDescribeCmd = &cobra.Command{
	Use:   "describe <organization> <workspace>",
	Args:  cobra.ExactArgs(2),
	Short: "Print info about a workspace",
	Long: `Print info about a workspace.
`,
	Run: func(cmd *cobra.Command, args []string) {
		organization := args[0]
		workspace := args[1]

		fmt.Printf("%s (ws: %s) info:\n\n", args[1], args[0])
		w, err := client.Workspaces.Read(ctx, organization, workspace)
		if err != nil {
			log.Fatalln("failed to get workspace info", err)
		}

		// fmt.Println(w)

		tw := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(tw, "ID:\t", w.ID)
		fmt.Fprintln(tw, "AutoApply:\t", w.AutoApply)
		fmt.Fprintln(tw, "CanQueueDestroyPlan:\t", w.CanQueueDestroyPlan)
		fmt.Fprintln(tw, "CreatedAt:\t", w.CreatedAt)
		fmt.Fprintln(tw, "Environment:\t", w.Environment)
		fmt.Fprintln(tw, "FileTriggersEnabled:\t", w.FileTriggersEnabled)
		fmt.Fprintln(tw, "Locked:\t", w.Locked)
		fmt.Fprintln(tw, "MigrationEnvironment:\t", w.MigrationEnvironment)
		fmt.Fprintln(tw, "Name:\t", w.Name)
		fmt.Fprintln(tw, "Operations:\t", w.Operations)
		fmt.Fprintln(tw, "QueueAllRuns:\t", w.QueueAllRuns)
		fmt.Fprintln(tw, "TerraformVersion:\t", w.TerraformVersion)
		fmt.Fprintln(tw, "TriggerPrefixes:\t", w.TriggerPrefixes)
		fmt.Fprintln(tw, "WorkingDirectory:\t", w.WorkingDirectory)
		tw.Flush()
	},
}

func init() {
	workspacesCmd.AddCommand(workspacesDescribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workspacesDescribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workspacesDescribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
