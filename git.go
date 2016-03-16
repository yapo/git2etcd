package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/libgit2/git2go.v23"
)

func openOrCloneRepo() error {
	var err error
	gitRepo, err = git.OpenRepository(viper.GetString("repo.path"))
	if err != nil || gitRepo == nil {
		log.WithError(err).Warn("Couldn't find repo locally, trying to clone it")
		cloneOptions := &git.CloneOptions{}
		cloneOptions.FetchOptions = &git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				CredentialsCallback:      credentialsCallback,
				CertificateCheckCallback: certificateCheckCallback,
			},
		}
		if viper.IsSet("repo.branch") {
			cloneOptions.CheckoutBranch = viper.GetString("repo.branch")
		}
		gitRepo, err = git.Clone(viper.GetString("repo.url"), viper.GetString("repo.path"), cloneOptions)
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
	if err := repo.CheckoutHead(nil); err != nil {
		return errors.New("Couldn't checkout head: " + err.Error())
	}
	head, err := repo.Head()
	if err != nil {
		return errors.New("Couldn't checkout head: " + err.Error())
	}
	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return errors.New("Couldn't lookup commit: " + err.Error())
	}
	tree, err := commit.Tree()
	if err != nil {
		return errors.New("Couldn't get commit tree: " + err.Error())
	}
	tree.Walk(walkCallback)
	log.Info("Repo synced")
	return nil
}

func walkCallback(name string, treeEntry *git.TreeEntry) int {
	if treeEntry.Type == git.ObjectTree {
		return 0
	}
	name += treeEntry.Name
	if !etcdExists(name) {
		if err := etcdCreate(name); err != nil {
			log.WithError(err).Warn("Couldn't create key")
		}
	} else {
		if err := etcdSet(name); err != nil {
			log.WithError(err).Warn("Couldn't set key")
		}
	}
	return 0
}

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	if viper.GetString("auth.type") == "ssh" {
		ret, cred := git.NewCredSshKey("git", viper.GetString("auth.ssh.public"), viper.GetString("auth.ssh.key"), viper.GetString("auth.ssh.passphrase"))
		return git.ErrorCode(ret), &cred
	}
	ret, cred := git.NewCredUserpassPlaintext(viper.GetString("auth.http.username"), viper.GetString("auth.http.password"))
	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}
