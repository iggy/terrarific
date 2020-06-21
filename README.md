# terrarific

terraform cloud (nee enterprise) cli written in golang

![Test](https://github.com/iggy/terrarific/workflows/Test/badge.svg?branch=master)

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

## Disclaimer

I work for Pluto TV (part of ViacomCBS), but I wrote this on my own time.
