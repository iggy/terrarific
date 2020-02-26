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
