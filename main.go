package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	_sep = string(os.PathSeparator)
)

var (
	flagConfigPath = flag.String("conf_dir", "", "Path to look for a config file. (directory)")
	gitRepo        *git.Repository
	etcdClient     etcd.KeysAPI
)

func main() {
	flag.Parse()
	setConfig(*flagConfigPath)

	// etcd Client connection
	if err := etcdConnect(); err != nil {
		log.WithError(err).Fatal("Couldn't connect to etcd")
	}

	// Git repository opening/cloning
	if err := openOrCloneRepo(); err != nil {
		log.WithError(err).Fatal("Couldn't find repo or clone it")
	}

	go func(repo *git.Repository) {
		syncCycle := time.Duration(viper.GetInt("repo.synccycle")) * time.Second
		if syncCycle > 0 {
			for {
				select {
				case <-time.After(syncCycle):
					if err := syncRepo(repo); err != nil {
						log.WithError(err).Warn("Couldn't sync automatically")
					}
				}
			}
		} else {
			log.Info("No sync cycle")
		}
	}(gitRepo)

	// HTTP serving
	http.HandleFunc("/"+viper.GetString("host.hook"), hookHandler)
	http.HandleFunc("/sync", syncHandler)
	http.HandleFunc("/status", statusHandler)
	log.Fatal(http.ListenAndServe(viper.GetString("host.listen")+":"+viper.GetString("host.port"), nil))
}

func setConfig(path string) {

	// Default values
	viper.SetDefault("host.listen", "")
	viper.SetDefault("host.port", "4242")
	viper.SetDefault("host.hook", "hook")

	viper.SetDefault("repo.path", "data/")
	viper.SetDefault("repo.url", "https://github.com/djaque/git2etcd.git")
	viper.SetDefault("repo.branch", "master")
	viper.SetDefault("repo.synccycle", 3600)

	viper.SetDefault("etcd.hosts", []string{"http://127.0.0.1:2379"})

	// Getting config from file
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/git2etcd/")
	viper.AddConfigPath("$HOME/.git2etcd")
	viper.AddConfigPath(".")
	if len(path) > 0 {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Warn("Couldn't read config file. Will use defaults.")
	}

	// Setting environment config
	viper.SetEnvPrefix("g2e")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	_, err := etcdClient.Get(context.Background(), "/", nil)
	if err != nil && err == etcd.ErrClusterUnavailable {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	if err := syncRepo(gitRepo); err != nil {
		http.Error(w, "Couldn't sync repo: "+err.Error(), http.StatusInternalServerError)
	}
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-GitHub-Event") == "ping" {
		log.Info("Ping received")
	} else if r.Header.Get("X-GitHub-Event") == "push" {
		treatPushEvent(w, r)
	} else {
		// handle other events than push.
	}
}

func treatPushEvent(w http.ResponseWriter, r *http.Request) {
	var event github.PushEvent
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Couldn't read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &event); err != nil {
		log.WithError(err).Error("Couldn't Unmarshal json payload")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Push received from ", *event.Repo.FullName)
	added := make(map[string]bool)
	modified := make(map[string]bool)
	removed := make(map[string]bool)
	for _, commit := range event.Commits {
		for _, ca := range commit.Added {
			added[ca] = true
		}
		for _, cm := range commit.Modified {
			modified[cm] = true
		}
		for _, cr := range commit.Removed {
			removed[cr] = true
		}
	}
	wt, err := gitRepo.Worktree()
	if err != nil {
		log.WithError(err).Error("Couldn't get repo's WorkTree")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = wt.Pull(nil); err != nil {
		log.WithError(err).Error("Couldn't pull repo's HEAD")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Repository head is now ", *event.After)
	commit, err := gitRepo.CommitObject(plumbing.NewHash(*event.After))
	if err != nil {
		log.WithError(err).Error("Couldn't get HEAD's commit")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tree, err := commit.Tree()
	if err != nil {
		log.WithError(err).Error("Couldn't get HEAD's tree")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	treatAdded(added, tree)
	treatModified(modified, tree)
	treatRemoved(removed, tree)
}

func treatAdded(added map[string]bool, tree *object.Tree) {
	for file := range added {
		_, err := tree.File(file)
		if err != nil {
			log.WithError(err).Warn("Couldn't get file: ", file)
			continue
		}
		if err := etcdCreate(file); err != nil {
			log.WithError(err).Warn("Couldn't create key")
		}
	}
}

func treatModified(modified map[string]bool, tree *object.Tree) {
	for file := range modified {
		_, err := tree.File(file)
		if err != nil {
			log.WithError(err).Warn("Couldn't get file: ", file)
			continue
		}
		if err := etcdSet(file); err != nil {
			log.WithError(err).Warn("Couldn't set key")
		}
	}
}

func treatRemoved(removed map[string]bool, tree *object.Tree) {
	for file := range removed {
		_, err := tree.File(file)
		if err != nil {
			log.WithError(err).Warn("Couldn't get file: ", file)
			continue
		}
		if err := etcdDelete(file); err != nil {
			log.WithError(err).Warn("Couldn't delete key")
		}
	}
}
