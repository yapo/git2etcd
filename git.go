package main

import (
	"errors"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitobj "gopkg.in/src-d/go-git.v4/plumbing/object"
	gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func openOrCloneRepo() error {
	var err error
	gitRepo, err = git.PlainOpen(viper.GetString("repo.path"))
	if err != nil || gitRepo == nil {
		log.WithError(err).Warn("Couldn't find repo locally, trying to clone it")
		cloneOptions := &git.CloneOptions{}
		cloneOptions.URL = viper.GetString("repo.url")
		cloneOptions.Auth, err = getGitAuth()
		if err != nil {
			return err
		}
		if viper.GetString("repo.branch") == "" {
			// Default value is not correctly assigned to repo.branch when using json config, forcing it here
			viper.Set("repo.branch", "master")
		}
		cloneOptions.SingleBranch = true
		cloneOptions.ReferenceName = plumbing.ReferenceName("refs/heads/" + viper.GetString("repo.branch"))
		log.WithFields(log.Fields{
			"url":    viper.GetString("repo.url"),
			"branch": viper.GetString("repo.branch"),
			"path":   viper.GetString("repo.path"),
		}).Info("Cloning repo")
		gitRepo, err = git.PlainClone(viper.GetString("repo.path"), false, cloneOptions)
		if err != nil {
			return err
		}
	}
	if err := syncRepo(gitRepo); err != nil {
		log.WithError(err).Warn("Couldn't sync repo")
	}
	return nil
}

func syncRepo(repo *git.Repository) error {
	wt, err := repo.Worktree()
	if err != nil {
		return errors.New("Couldn't get WorkTree: " + err.Error())
	}
	po := &git.PullOptions{}
	po.Auth, err = getGitAuth()
	if err != nil {
		return err
	}
	if err := wt.Pull(po); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.New("Couldn't pull: " + err.Error())
	}
	head, err := repo.Head()
	if err != nil {
		return errors.New("Couldn't checkout head: " + err.Error())
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return errors.New("Couldn't get commit: " + err.Error())
	}
	tree, err := commit.Tree()
	if err != nil {
		return errors.New("Couldn't get commit tree: " + err.Error())
	}
	err = tree.Files().ForEach(func(f *gitobj.File) error {
		if !etcdExists(f.Name) {
			if err := etcdCreate(f.Name); err != nil {
				log.WithError(err).WithField("name", f.Name).Warn("Couldn't create key")
			}
		} else {
			if err := etcdSet(f.Name); err != nil {
				log.WithError(err).WithField("name", f.Name).Warn("Couldn't set key")
			}
		}
		return nil
	})
	if err != nil {
		return errors.New("Couldn't walk in files: " + err.Error())
	}
	log.Info("Repo synced")
	return nil
}

func getGitAuth() (gittransport.AuthMethod, error) {
	if viper.GetString("auth.type") == "ssh" {
		var signer ssh.Signer
		sshFile, err := os.Open(viper.GetString("auth.ssh.key"))
		if err != nil {
			return nil, errors.New("Couldn't open SSH key: " + err.Error())
		}
		sshB, err := ioutil.ReadAll(sshFile)
		if err != nil {
			return nil, errors.New("Couldn't read SSH key: " + err.Error())
		}
		if viper.GetString("auth.ssh.passphrase") != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(sshB, []byte(viper.GetString("auth.ssh.passphrase")))
		} else {
			signer, err = ssh.ParsePrivateKey(sshB)
		}
		if err != nil {
			return nil, errors.New("Couldn't parse SSH key: " + err.Error())
		}
		sshAuth := &gitssh.PublicKeys{User: "git", Signer: signer}
		return sshAuth, nil
	}
	httpAuth := &githttp.BasicAuth{
		Username: viper.GetString("auth.http.username"),
		Password: viper.GetString("auth.http.password"),
	}
	return httpAuth, nil
}
