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
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

// organizationsListCmd represents the list command
var organizationsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List organizations that your API token can access",
	Long: `Output a table of the terraform cloud/enterprise organizations your API key has access 
to. The fields displayed are Name, Email, Creation Date, Enterprise Plan, and Two Factor Auth`,
	Run: func(cmd *cobra.Command, args []string) {
		pgNum := 0
		pgSize := 100

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		fmt.Fprintln(w, "Name\tEmail\tCreated\tPlan\tTwoFactor")

		for {
			opts := tfe.OrganizationListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: pgNum,
					PageSize:   pgSize,
				},
			}
			list, err := client.Organizations.List(ctx, opts)
			if err != nil {
				log.Fatalln("Failed to list organizations", err)
			}
			for _, org := range list.Items {
				fmt.Fprintf(
					w,
					"%s\t%s\t%s\t%t\n",
					org.Name,
					org.Email,
					org.CreatedAt.Format("2006-01-02"),
					org.TwoFactorConformant,
				)
			}

			if list.Pagination.CurrentPage == list.Pagination.TotalPages {
				break
			}
			pgNum++
		}
		w.Flush()
	},
}

func init() {
	organizationsCmd.AddCommand(organizationsListCmd)
}
