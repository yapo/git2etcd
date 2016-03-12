git2etcd
======
[![GoDoc](https://godoc.org/github.com/blippar/git2etcd?status.svg)](http://godoc.org/github.com/blippar/git2etcd) [![Build Status](https://ci.userctl.xyz/api/badges/blippar/git2etcd/status.svg)](https://ci.userctl.xyz/blippar/git2etcd)


Simple binary to sync a Git repository with an etcd config. Built and tested with Go 1.4+

Installing
----------

First, install libgit2, then :

```
go get github.com/blippar/git2etcd
```

Configuring
-------

Key | Description | Default
----|-------------|--------
`host.listen` | Host to listen to | `""`
`host.port` | Port to listen to | `"4242"`
`host.hook` | Name of the Webhook endpoint | `"hook"`
`repo.url` | URL of the repo to sync | `"https://github.com/0rax/fishline.git"`
`repo.path` | Path where to clone the repo | `"/opt/git2etcd/repo"`
`etcd.hosts` | List of etcd hosts | `["http://127.0.0.1:2379"]`
`auth.type`  | Type of authentication for Git | `"ssh"`
`auth.ssh.key` | Path to the SSH private key (if ssh auth type) | `"~/.ssh/id_rsa"`
`auth.ssh.public` | Path to the SSH public key (if ssh auth type) | `"~/.ssh/id_rsa.pub"`

Contributing
-------

We'd love to get your feedback with [issues](https://github.com/blippar/git2etcd/issues/new) or even [pull requests](https://github.com/blippar/git2etcd/pulls).

Authors
-------

- Clem Dal Palu ([@dal-papa](http://www.github.com/dal-papa))
- JP Roemer ([0rax](http://www.github.com/0rax))
