Use the Github API to fetch configured SSH keys for a team.

Build by (making sure your GOPATH environment variable is properly configured) - 
```
go get
go build
```

Run
```
./fetch-keys --help
```

The tool uses an OAuth Token to connect to the Github API, make sure you've added a personal access token via your [Github account settings](https://github.com/settings/tokens) and have either properly edited the [authtoken.go](authtoken.go#L3) source file (!!! with caution), or pass the appropriate token via the command line.
