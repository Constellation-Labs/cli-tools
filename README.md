DAG Command Line Utility Tools
---
![Release](https://img.shields.io/github/v/release/constellation-labs/cli-tools?style=for-the-badge)
---

```
DAG Command Line Utility Tools

Usage:
  cli-tools [flags]
  cli-tools [command]

Available Commands:
  completion   Generate the autocompletion script for the specified shell
  fix-balances Fix wrongly indexed zero balances
  help         Help about any command

Flags:
  -h, --help      help for cli-tools
      --verbose   verbose output
  -v, --version   version for cli-tools

Use "cli-tools [command] --help" for more information about a command.
```

## Fix balances
```
Fix wrongly indexed zero balances

Usage:
  cli-tools fix-balances [flags]

Flags:
      --dry-run             When true, works in read-only mode
      --from int            Starting ordinal
  -h, --help                help for fix-balances
      --opensearch string   Opensearch url
      --to int              Ending ordinal
      --workers int         No. of parallel workers (default 10)

Global Flags:
      --verbose   verbose output
```
