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
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

// workspacesEnsureCmd represents the workspacesEnsure command
var workspacesEnsureCmd = &cobra.Command{
	Use:   "ensure <organization> <workspace>",
	Args:  cobra.ExactArgs(2),
	Short: "A shortcut to create/update a workspace to match the args",
	Long: `If a workspace doesn't exist, create it and update it to match the args passed. If 
it does exist, just update it to match the args.`,
	Run: func(cmd *cobra.Command, args []string) {
		organization := args[0]
		workspace := args[1]
		envFlags, err := cmd.Flags().GetStringArray("env")
		if err != nil {
			log.Fatalln("can't get env flag(s)", err)
		}
		versFlag, err := cmd.Flags().GetString("version")
		if err != nil {
			log.Fatalln("failed to get version flag", err)
		}
		autoApplyFlag, err := cmd.Flags().GetBool("autoapply")
		if err != nil {
			log.Fatalln("failed to get autoapply flag", err)
		}
		remoteFlag, err := cmd.Flags().GetBool("remote")
		if err != nil {
			log.Fatalln("failed to get remote flag", err)
		}
		workDirFlag, err := cmd.Flags().GetString("workingdir")
		if err != nil {
			log.Fatalln("failed to get workdir flag", err)
		}

		// Check if workspace exists already
		w, err := client.Workspaces.Read(ctx, organization, workspace)
		if err != nil {
			log.Printf("Workspace %s doesn't exist. Creating (%q)", workspace, err)

			w, err = client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{
				AutoApply:           tfe.Bool(autoApplyFlag),
				FileTriggersEnabled: tfe.Bool(false),
				Name:                tfe.String(workspace),
				Operations:          tfe.Bool(remoteFlag),
				TerraformVersion:    tfe.String(versFlag),
				TriggerPrefixes:     []string{},
				WorkingDirectory:    tfe.String(workDirFlag),
			})
			if err != nil {
				log.Fatal("failed creating workspace: ", err, w)
			}

			// Since the workspace didn't exist, we can assume the ENV vars didn't either
			for ev := range envFlags {
				envkv := strings.Split(envFlags[ev], "=")
				log.Printf("Creating variable on workspace %s: %q\n", workspace, envkv)
				v, nerr := client.Variables.Create(ctx, w.ID, tfe.VariableCreateOptions{
					Key:         &envkv[0],
					Value:       &envkv[1],
					Description: tfe.String("Set by terrarific"),
					Category:    tfe.Category(tfe.CategoryEnv),
				})
				if nerr != nil {
					log.Fatal("failed setting ", err, v)
				}
			}

		} else {
			log.Printf("workspace %s already exists. Updating attrs if necessary", w.Name)

			// Update the workspace
			w, err = client.Workspaces.Update(ctx, organization, workspace, tfe.WorkspaceUpdateOptions{
				AutoApply:        tfe.Bool(false),
				TerraformVersion: tfe.String(versFlag),
			})
			if err != nil {
				log.Fatal("failed updating workspace attrs (AutoApply|TerraformVersion): ", err)
			}

			// check if the vars exist and create/update appropriately
			for ev := range envFlags {
				envkv := strings.Split(envFlags[ev], "=")
				key := envkv[0]
				value := envkv[1]
				foundEnvInVarList := false

				varList, err := client.Variables.List(ctx, w.ID, &tfe.VariableListOptions{})
				if err != nil {
					log.Fatal("Failed getting list of variables", err)
				}
				for li := range varList.Items {
					if varList.Items[li].Key == key {
						foundEnvInVarList = true
						if varList.Items[li].Value != value || !varList.Items[li].Sensitive {
							// key exists, but value/Sensitive doesn't match, update it
							vu, err := client.Variables.Update(ctx, w.ID, varList.Items[li].ID, tfe.VariableUpdateOptions{
								Key:         &key,
								Value:       &value,
								Description: tfe.String("Set by terrarific"),
								Sensitive:   tfe.Bool(true),
							})
							if err != nil {
								log.Fatal("Failed updating existing env var ", err, vu)
							}
						}
					}
				}
				if !foundEnvInVarList {
					// we never found the env var in the list from tf cloud, create it
					v, err := client.Variables.Create(ctx, w.ID, tfe.VariableCreateOptions{
						Key:         &key,
						Value:       &value,
						Description: tfe.String("Set by terrarific"),
						Category:    tfe.Category(tfe.CategoryEnv),
						Sensitive:   tfe.Bool(true),
					})
					if err != nil {
						log.Fatal("failed setting ", err, v)
					}
				}
			}
		}
	},
}

func init() {
	workspacesCmd.AddCommand(workspacesEnsureCmd)

	workspacesEnsureCmd.Flags().Bool("autoapply", false, "Whether to automatically apply after a successful plan")
	workspacesEnsureCmd.Flags().StringArrayP("env", "e", []string{}, "Environment variables to set on the workspace. Can be passed multiple times. (-e KEY=Value)")
	workspacesEnsureCmd.Flags().StringP("version", "v", "latest", "Version of terraform to run in the terraform cloud runners (ex. latest, 0.12.29, 0.11.14")
	workspacesEnsureCmd.Flags().Bool("remote", true, "Whether to run commands in terraform cloud runners (remote) or locally on a workstation")
	workspacesEnsureCmd.Flags().String("workingdir", "", "Directory to run terraform commands in")

	// workspacesEnsureCmd.Flags().StringArrayP("tfvar", "t", []string{}, "Terraform variables to set on the workspace. Can be passed multiple times. (-t KEY=Value)")
}
