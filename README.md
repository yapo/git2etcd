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
`auth.http.username` | Username (if HTTP auth type) | `""`
`auth.http.password` | Password (if HTTP auth type) | `""`

#### JSON file
You can use a JSON config file that you would put either in current folder or in a folder you can precise with the `-config` flag. 

~~~json
{
  "host": {
    "listen": "",
    "port": "4242",
    "hook": "hook"
  },
  "repo": {
    "url": "https://github.com/0rax/fishline.git",
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
      "key": "~/.ssh/id_rsa",
      "public": "~/.ssh/id_rsa.pub"
    }
  }
}
~~~

> I don't speak JSON !

Well, you can use TOML, YAML, HCL ... 
#### Env vars

Who needs a file when you can use environment variables ? `host.port` can be `G2E_HOST_POST` and so on.

Contributing
-------

We'd love to get your feedback with [issues](https://github.com/blippar/git2etcd/issues/new) or even [pull requests](https://github.com/blippar/git2etcd/pulls).

Authors
-------

- Clem Dal Palu ([@dal-papa](http://www.github.com/dal-papa))
- JP Roemer ([0rax](http://www.github.com/0rax))
