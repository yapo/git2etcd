# git2etcd
[![Go Report Card](https://goreportcard.com/badge/github.com/yapo/git2etcd)](https://goreportcard.com/report/github.com/yapo/git2etcd) [![GoDoc](https://godoc.org/github.com/yapo/git2etcd?status.svg)](http://godoc.org/github.com/yapo/git2etcd) [![GitHub release](https://img.shields.io/github/release/yapo/git2etcd.svg)](https://github.com/yapo/git2etcd/releases/latest) [![LICENSE](https://img.shields.io/github/license/yapo/git2etcd.svg)](https://github.com/yapo/git2etcd/blob/master/LICENSE)


Simple binary to sync a Git repository with an etcd config. Built and tested with Go 1.4+

## Installing

### Docker

```
docker pull yapo/git2etcd
```

### Manually

```
go get github.com/yapo/git2etcd
```

## Configuring

Key                  | Description                    | Default
---------------------|--------------------------------|--------
`host.listen`        | Host to listen to              | `""`
`host.port`          | Port to listen to              | `"4242"`
`host.hook`          | Name of the Webhook endpoint   | `"hook"`
`repo.url`           | URL of the repo to sync        | `"https://github.com/yapo/git2etcd.git"`
`repo.branch`        | Branch of the repo to sync     | `"master"`
`repo.path`          | Path where to clone the repo   | `"data/"`
`repo.synccycle`     | Number of seconds between 2 automatic syncs (if 0, never syncs) | `3600`
`etcd.hosts`         | List of etcd hosts             | `["http://127.0.0.1:2379"]`
`auth.type`          | Type of authentication for Git | `n/a`
`auth.ssh.key`       | Path to the SSH private key (if `ssh` auth type) | `n/a`
`auth.ssh.public`    | Path to the SSH public key (if `ssh` auth type)  | `n/a`
`auth.http.username` | Username (if `http` auth type)   | `n/a`
`auth.http.password` | Password (if `http` auth type)   | `n/a`

#### JSON file
You can use a JSON config file that you would put either in current folder or in a folder you can precise with the `-conf_dir` flag.

```json
{
  "host": {
    "listen": "",
    "port": "4242",
    "hook": "hook"
  },
  "repo": {
    "url": "git@github.com:yapo/git2etcd.git",
    "branch": "master",
    "path": "/opt/git2etcd/repo"
  },
  "etcd": {
    "hosts": [
      "http://127.0.0.1:2379"
    ]
  },
  "auth": {
    "type": "ssh",
    "ssh": {
      "key": "/home/user/.ssh/id_rsa"
    }
  }
}
```

> I don't speak JSON !

Well, you can use TOML, YAML, HCL ...
#### Env vars

Who needs a file when you can use environment variables ? `host.port` can be `G2E_HOST_POST` and so on.

## Contributing

We'd love to get your feedback with [issues](https://github.com/yapo/git2etcd/issues/new) or even [pull requests](https://github.com/yapo/git2etcd/pulls).

## Authors

- Clem Dal Palu ([@dal-papa](http://www.github.com/dal-papa))
- JP Roemer ([0rax](http://www.github.com/0rax))
