package main

import (
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
	return nil
}

func syncRepo() error {
	return nil
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
