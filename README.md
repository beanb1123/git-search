# git-search
Tool for search a pattern regexp into a git repository for all remote branches.

# Overview
Powered by [Go-Git](https://github.com/go-git/go-git/v5),
this tool clones a repo in memory and searches for a string *(or regex)* for all remote branches.

For install, use `go install`
````shell
go install github.com/agustin-del-pino/git-search
````

Once install, just execute it.

````shell
gss "YOUR-STRING-OR-REGEX" -r "REPO-URL" -u "USERNAME" -p "PASSWORD_OR_TOKEN"
````

Then a `result.json` will be generated that contains the collected information. It looks like this.

````json
{
  "repo_name": {
    "branch_name": {
      "file_name": [
        "content_matched"
      ]
    }
  }
}
````

*For each repository where the search will make a new goroutine starts.*  

# Reference
*The `-h` flag provides the **Help Information** of the command.*

## Arguments
There only one positional argument `SEARCH` which is a 
`STRING|REGEX` to be used for the *grep process*.

## Flags
| Flag           | Shorthand | Description                                                              |
|----------------|-----------|--------------------------------------------------------------------------|
| `--help`       | `-h`      | the command help                                                         |
| `--user`       | `-u`      | set the username for auth (mandatory)                                    |
| `--password`   | `-p`      | set the password fot auth (mandatory)                                    |
| `--repo`       | `-r`      | add a repo name (accumulative) (mandatory if `-f` is not set)            |
| `--file-repos` | `-f`      | set the source file to get the repo names (mandatory if `-r` is not set) |
| `--out-file`   | `-o`      | set the output filepath (default: `./result.json`)                       |
