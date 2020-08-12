# terrarific

terraform cloud (nee enterprise) cli written in golang

![Tests](https://github.com/iggy/terrarific/workflows/Tests/badge.svg)

## TODO

* Add functionality to cover the rest of the terraform cloud API
* when go-tfe supports source-name/source-url use it [https://www.terraform.io/docs/cloud/api/workspaces.html#request-body]
        * See if tfc can add source-name/source-url to other resources (variables, etc)


## Usage

```text
List, create, update, delete, etc different things (organizations, workspaces,
variables, etc) in terraform cloud

Usage:
  terrarific [command]

Available Commands:
  completion    Generate shell completion code
  help          Help about any command
  organizations Work with organizations
  workspaces    Work with workspaces

Flags:
      --config string   config file (default is $HOME/.terrarific.yaml)
  -h, --help            help for terrarific
  -v, --version         version for terrarific

Use "terrarific [command] --help" for more information about a command.
```

workspaces
```text
Parent command for manipulating workspaces. This doesn't do anything by itself.
Everything is done via subcommands.

Usage:
  terrarific workspaces [command]

Available Commands:
  create      Create a new workspace
  describe    Print info about a workspace
  ensure      A shortcut to create/update a workspace to match the args
  list        List workspaces in an organization

Flags:
  -h, --help   help for workspaces

Global Flags:
      --config string   config file (default is $HOME/.terrarific.yaml)

Additional help topics:
  terrarific workspaces variables Work with workspace variables

Use "terrarific workspaces [command] --help" for more information about a command.
```

organizations
```text
Parent command for manipulating organizations. This doesn't do anything by itself.
Everything is done via subcommands.

Usage:
  terrarific organizations [command]

Available Commands:
  list        List organizations that your API token can access

Flags:
  -h, --help   help for organizations

Global Flags:
      --config string   config file (default is $HOME/.terrarific.yaml)

Use "terrarific organizations [command] --help" for more information about a command.
```

## Releasing

To create a new release, all you have to do is push a tag (no need to create a separate GH
release). CI will do that for you.

```shell
git tag -a v0.0.4
git push origin --tags
```

## Disclaimer

I work for Pluto TV (part of ViacomCBS), but I wrote this on my own time.
