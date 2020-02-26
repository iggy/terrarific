// Copyright 2020 pluto.tv
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
//

// terraform cloud (nee enterprise) cli written in golang
// For now it only creates workspaces, but the idea is to expand the functionality later
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

type sliceFlags []string

func (i *sliceFlags) String() string {
	return "foo"
}

func (i *sliceFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

// for future use
// func workspaceCreate() {
//   log.Println("workspaceCreate")
// }

// the terraform cloud organization and workspace we want to add/make changes to
var organization = flag.String("org", "", "The terraform cloud organization to use.")
var workspace = flag.String("workspace", "", "The workspace name.")

func main() {
	var envFlags sliceFlags
	flag.Var(&envFlags, "env", "Environment variables to set on the workspace (-env KEY=Value)")
	flag.Parse()

	config := &tfe.Config{
		Token: os.Getenv("TERRAFORM_CLOUD_TOKEN"),
	}

	// create the client
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// Check if workspace exists already
	w, err := client.Workspaces.Read(ctx, *organization, *workspace)
	if err != nil {
		log.Print("Failed reading workspace: ", err)

		// Create a new workspace
		w, err = client.Workspaces.Create(ctx, *organization, tfe.WorkspaceCreateOptions{
			Name:             tfe.String(*workspace),
			AutoApply:        tfe.Bool(false),
			TerraformVersion: tfe.String("latest"),
		})
		if err != nil {
			log.Fatal("failed creating workspace: ", err, w)
		}

		// Since the workspace didn't exist, we can assume the ENV vars didn't either
		for ev := range envFlags {
			envkv := strings.Split(envFlags[ev], "=")
			log.Printf("%q\n", envkv)
			v, err := client.Variables.Create(ctx, w.ID, tfe.VariableCreateOptions{
				Key:   &envkv[0],
				Value: &envkv[1],
			})
			if err != nil {
				log.Fatal("failed setting ")
			}
			log.Print("v = ", v)
		}

	} else {
		log.Print("workspace already exists. Updating attrs if necessary")
		log.Print(w)

		// Update the workspace
		w, err = client.Workspaces.Update(ctx, *organization, w.Name, tfe.WorkspaceUpdateOptions{
			AutoApply:        tfe.Bool(false),
			TerraformVersion: tfe.String("latest"),
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

			log.Printf("%v:%v\n", key, value)
			varList, err := client.Variables.List(ctx, w.ID, tfe.VariableListOptions{})
			if err != nil {
				log.Fatal("Failed getting list of variables", err)
			}
			for li := range varList.Items {
				log.Printf("li: %v | %v\n", li, varList.Items[li])
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
}
