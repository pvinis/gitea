package models

import (
	"fmt"
	"errors"
	"io/ioutil"
	"path/filepath"
	"os"
	"time"

	"github.com/Unknwon/com"

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

	err = p.create(u)
	if err != nil {
		return err
	}

	return nil
}

func (p *WikiPage) create(u *User) error {
	repoPath, err := p.Repo.RepoPath()

	tmpDir := filepath.Join(os.TempDir(), com.ToStr(time.Now().Nanosecond()))
	os.MkdirAll(tmpDir, os.ModePerm)

	_, stderr, err := process.Exec(
		fmt.Sprintf("WikiPage create(git clone): %s", repoPath),
		"git", "clone", repoPath, tmpDir)
	if err != nil {
		return errors.New("git clone: " + stderr)
	}

	filename := sanitize.Path(p.Title)
	if err := ioutil.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("%s.md", filename)),
		[]byte(p.Content), 0644); err != nil {
		return err
	}
	p.Alias = filename

	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit(git add): %s", tmpDir),
		"git", "add", "--all"); err != nil {
		return errors.New("git add: " + stderr)
	}

	sig := u.NewGitSig()
	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit(git commit): %s", tmpDir),
		"git", "commit", fmt.Sprintf("--author='%s <%s>'", sig.Name, sig.Email),
		"-m", fmt.Sprintf("Create page %s", p.Title)); err != nil {
		return errors.New("git commit: " + stderr)
	}

	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit(git push): %s", tmpDir),
		"git", "push", "origin", "master"); err != nil {
		return errors.New("git push: " + stderr)
	}

	return nil
}