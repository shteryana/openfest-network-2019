Use the Github API to fetch configured SSH keys for a team.

Build by
```
go get
go build
```
Run
```
./fetch-keys --help
```

TODO: This polutes the current working directory with `{user}.key` files - add
a command line option for custom output target dir
