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

	tfe "github.com/hashicorp/go-tfe"
)

// for future use
// func workspaceCreate() {
//   log.Println("workspaceCreate")
// }

// the terraform cloud organization and workspace we want to add/make changes to
var organization = flag.String("org", "", "The terraform cloud organization to use.")
var workspace = flag.String("workspace", "", "The workspace name.")

func main() {
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
			Name: tfe.String(*workspace),
			AutoApply:        tfe.Bool(false),
			TerraformVersion: tfe.String("latest"),
		})
		if err != nil {
			log.Fatal("failed creating workspace: ", err)
		}
		log.Print(w)

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
		log.Print(w)
	}
}
