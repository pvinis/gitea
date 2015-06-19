package models

import (
	"fmt"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-gitea/gitea/modules/git"
	"github.com/go-gitea/gitea/modules/log"
	"github.com/go-gitea/gitea/modules/process"

	"github.com/kennygrant/sanitize"
)

type WikiPage struct {
	Title    string
	Alias    string
	Content  string
	Repo     *Repository
}

func NewWikiPage(p *WikiPage, r *Repository, u *User) error {

	var err error

	if r.WikiRepo != nil {
		p.Repo = r.WikiRepo
	} else {
		p.Repo, err = CreateRepository(r.Owner, fmt.Sprintf("%s.wiki", r.Name), "", "", "", false, false, false, true, r.Id)
		if err != nil {
			return err
		}
	}

	err = p.update(u)
	if err != nil {
		return err
	}

	return nil
}

func (p *WikiPage) update(u *User) error {
	repoPath, err := p.Repo.RepoPath()
	repoRealPath, err := p.Repo.WikiRepoPath()

	f, err := os.Open(repoRealPath)
	if err != nil {
		userDir := filepath.Join(repoRealPath, "../")
		os.MkdirAll(userDir, os.ModePerm)

		_, stderr, err := process.Exec(
			fmt.Sprintf("WikiPage create(git clone): %s", repoPath),
			"git", "clone", repoPath, repoRealPath)
		if err != nil {
			return errors.New("git clone: " + stderr)
		}
	}
	f.Close()

	filename := sanitize.Path(p.Title)

	fp := filepath.Join(repoRealPath, fmt.Sprintf("%s.md", filename))
	f, err = os.Open(fp)
	commitMsg := "Edit"
	if err != nil {
		commitMsg = "Create"
	}
	f.Close()

	if err := ioutil.WriteFile(fp,
		[]byte(p.Content), 0644); err != nil {
		return err
	}

	p.Alias = filename

	if err = p.commitChanges(repoRealPath, commitMsg, u.NewGitSig()); err != nil {
		return err
	}

	return nil
}

func (p *WikiPage) Delete(u *User) error {
	wikiRepoPath, err := p.Repo.WikiRepoPath()
	if err != nil {
		return err
	}

	if err = os.Remove(fmt.Sprintf("%s/%s.md", wikiRepoPath, p.Alias)); err != nil {
		return err
	}

	if err = p.commitChanges(wikiRepoPath, "Delete", u.NewGitSig()); err != nil {
		return err
	}

	return nil
}

func (p *WikiPage) commitChanges(repoRealPath, commitMsg string, sig *git.Signature) error {
	var (
		stderr string
		err    error
	)

	if _, stderr, err = process.ExecDir(-1,
		repoRealPath, fmt.Sprintf("WikiPage.create(git add): %s", repoRealPath),
		"git", "add", "--all"); err != nil {
		return errors.New("git add: " + stderr)
	}

	if _, stderr, err = process.ExecDir(-1,
		repoRealPath, fmt.Sprintf("WikiPage.create(git commit): %s", repoRealPath),
		"git", "commit", fmt.Sprintf("--author='%s <%s>'", sig.Name, sig.Email),
		"-m", fmt.Sprintf("%s page %s", commitMsg, p.Title)); err != nil {
		return errors.New("git commit: " + stderr)
	}

	if _, stderr, err = process.ExecDir(-1,
		repoRealPath, fmt.Sprintf("WikiPage.create(git push): %s", repoRealPath),
		"git", "push", "origin", "master"); err != nil {
		return errors.New("git push: " + stderr)
	}

	return err
}

func GetWikiPage(r *Repository, a string) (*WikiPage, error) {
	p := &WikiPage {
		Alias: a,
		Repo:  r,
	}

	wikiRepoPath, err := r.WikiRepoPath()
	if err != nil {
		return nil, err
	}

	c, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.md", wikiRepoPath, p.Alias))
	if err != nil {
		return nil, err
	}

	p.Content = string(c)
	p.Title = strings.Title(strings.Replace(p.Alias, "-", "", -1))

	return p, nil
}

func WikiUpdate() {
	if err := x.Where("is_wiki = 1").Iterate(new(Repository), func(idx int, bean interface{}) error {
		r := bean.(*Repository)

		wikiRepoPath, err := r.WikiRepoPath()
		if err != nil {
			return err
		}

		if _, stderr, err := process.ExecDir(10*time.Minute,
			wikiRepoPath, fmt.Sprintf("WikiUpdate: %s", wikiRepoPath),
			"git", "pull", "origin", "master"); err != nil {
			desc := fmt.Sprintf("Fail to update wiki repository(%s): %s", wikiRepoPath, stderr)
			log.Error(4, desc)
			return nil
		}

		return nil
	}); err != nil {
		log.Error(4, "WikiUpdate: %v", err)
	}
}