package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	etcd "github.com/coreos/etcd/client"
	"github.com/spf13/viper"
)

func etcdConnect() error {
	hosts := []string{}
	if viper.IsSet("etcd.host") {
		hosts = []string{viper.GetString("etcd.host")}
	} else {
		hosts = viper.GetStringSlice("etcd.hosts")
	}
	cfg := etcd.Config{
		Endpoints:               hosts,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		return err
	}
	etcdClient = etcd.NewKeysAPI(cli)
	_, err = etcdClient.Get(context.Background(), "/foo", nil)
	if err != nil && err.Error() == etcd.ErrClusterUnavailable.Error() {
		return err
	}
	return nil
}

func etcdCreate(file string) error {
	fd, err := os.Open(viper.GetString("repo.path") + _sep + file)
	if err != nil {
		return errors.New("Couldn't open file " + file + " : " + err.Error())
	}
	val, err := ioutil.ReadAll(fd)
	if err != nil {
		return errors.New("Couldn't read file " + file + " : " + err.Error())
	}
	_, err = etcdClient.Create(context.Background(), file, strings.TrimSpace(string(val)))
	if err != nil {
		return errors.New("Couldn't create file " + file + " : " + err.Error())
	}
	return nil
}

func etcdSet(file string) error {
	fd, err := os.Open(viper.GetString("repo.path") + _sep + file)
	if err != nil {
		return errors.New("Couldn't open file " + file + " : " + err.Error())
	}
	val, err := ioutil.ReadAll(fd)
	if err != nil {
		return errors.New("Couldn't read file " + file + " : " + err.Error())
	}
	_, err = etcdClient.Set(context.Background(), file, strings.TrimSpace(string(val)), nil)
	if err != nil {
		return errors.New("Couldn't set file " + file + " : " + err.Error())
	}
	return nil
}

func etcdDelete(file string) error {
	_, err := os.Stat(viper.GetString("repo.path") + _sep + file)
	if err != nil && os.IsNotExist(err) {
		_, err = etcdClient.Delete(context.Background(), file, nil)
		if err != nil {
			return errors.New("Couldn't set file " + file + " : " + err.Error())
		}
	} else {
		return errors.New("File still exsits, not deleting " + file + " : " + err.Error())
	}
	return nil
}

func etcdExists(file string) bool {
	_, err := etcdClient.Get(context.Background(), file, nil)
	if err != nil && strings.Contains(err.Error(), "Key not found") {
		return false
	}
	return true
}
