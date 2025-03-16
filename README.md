# azstorage

A command line tool for daily operations on Azure Storage Accounts.

## Usage

### azstorage remove

```
Removes the blobs listed in the specified file

Usage:
  azstorage remove [flags]

Flags:
      --account string     the name of the storage account
      --container string   the name of the blob container
  -h, --help               help for remove
      --list-file string   path to the file containing directories to delete
      --processors int     the number of concurrent processors deleting found blobs (default 16)
      --walkers int        the number of concurrent directory walkers (default 4)
```

The options `account`, `container`, and `list-file` are required to specify.

The storage account key must be specified with environment variable `AZURE_STORAGE_ACCOUNT_KEY`.

The list file specified by `--list-file` option contains directories to delete on each line.

Example of list file:
```
dir1
dir2/foo
dir3/bar/baz
```

The blank lines and lines starting with `#` are skipped.

## How to build

```
make clean
make
```
