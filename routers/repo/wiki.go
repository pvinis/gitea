package repo

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gitea/gitea/models"
	"github.com/go-gitea/gitea/modules/auth"
	"github.com/go-gitea/gitea/modules/base"
	"github.com/go-gitea/gitea/modules/git"
	"github.com/go-gitea/gitea/modules/log"
	"github.com/go-gitea/gitea/modules/middleware"
)

const (
	WIKI_EMPTY     base.TplName = "repo/wiki/empty"
	WIKI_ADD       base.TplName = "repo/wiki/add"
	WIKI_VIEW      base.TplName = "repo/wiki/view"
	WIKI_PAGELIST  base.TplName = "repo/wiki/pagelist"
	WIKI_GIT       base.TplName = "repo/wiki/git"
)

func Wiki(ctx *middleware.Context) {
	ctx.Data["Title"] = ctx.Repo.Repository.Name

	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}

	canWrite, err := models.HasAccess(ctx.User, wr, models.ACCESS_MODE_WRITE)
	if err != nil {
		ctx.Handle(500, "wiki.Wiki", err)
		return
	}

	has, err := models.HasWikiPage(wr, "home")
	if err != nil {
		ctx.Handle(500, "wiki.Wiki", err)
		return
	}

	if !has {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}

	p, err := models.GetWikiPage(wr, "home")
	if err != nil {
		ctx.Handle(500, "wiki.Wiki", err)
		return
	}

	ctx.Data["Page"] = p
	ctx.Data["CanWrite"] = canWrite
	ctx.Data["FileContent"] = string(base.RenderMarkdown([]byte(p.Content), ctx.Repo.WikiLink))
	ctx.HTML(200, WIKI_VIEW)
}

func ViewWikiPage(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}

	p, err := models.GetWikiPage(wr, ctx.Params(":slug"))
	if err != nil {
		ctx.Handle(404, "wiki.ViewWikiPage", err)
	}

	canWrite, err := models.HasAccess(ctx.User, wr, models.ACCESS_MODE_WRITE)
	if err != nil {
		ctx.Handle(500, "wiki.Wiki", err)
	}

	ctx.Data["Page"] = p
	ctx.Data["CanWrite"] = canWrite
	ctx.Data["FileContent"] = string(base.RenderMarkdown([]byte(p.Content), ctx.Repo.WikiLink))
	ctx.HTML(200, WIKI_VIEW)
}

func CreateWikiPage(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.Data["PageTitle"] = "Home"
	} else {
		has, err := models.HasWikiPage(wr, "Home")
		if err != nil {
			ctx.Handle(500, "wiki.CreateWikiPage", err)
			return
		}
		if !has {
			ctx.Data["PageTitle"] = "Home"
		}
	}

	ctx.Data["IsCreate"] = true
	canWrite, err := models.HasAccess(ctx.User, ctx.Repo.Repository, models.ACCESS_MODE_WRITE)
	if err != nil {
		ctx.Handle(500, "wiki.CreateWikiPage", err)
	}

	if !canWrite {
		ctx.Handle(401, "wiki.CreateWikiPage", errors.New(ctx.Tr("wiki.rights")))
	}

	ctx.HTML(200, WIKI_ADD)
}

func EditWikiPage(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.Data["PageTitle"] = "Home"
	}

	ctx.Data["IsCreate"] = false
	p, err := models.GetWikiPage(wr, ctx.Params(":slug"))
	if err != nil {
		ctx.Handle(404, "wiki.ViewWikiPage", err)
	}

	canWrite, err := models.HasAccess(ctx.User, ctx.Repo.Repository, models.ACCESS_MODE_WRITE)
	if err != nil {
		ctx.Handle(500, "wiki.EditWikiPage", err)
	}

	if !canWrite {
		ctx.Handle(401, "wiki.EditWikiPage", errors.New(ctx.Tr("wiki.rights")))
	}

	ctx.Data["Page"] = p
	ctx.Data["PageTitle"] = p.Title
	ctx.HTML(200, WIKI_ADD)
}

func CreateWikiPagePost(ctx *middleware.Context, form auth.CreateWikiPageForm) {
	send := func(status int, data interface{}, err error) {
		if err != nil {
			log.Error(4, "wiki.CreateWikiPagePost(?): %s", err)

			ctx.JSON(status, map[string]interface{}{
				"ok":     false,
				"status": status,
				"error":  err.Error(),
			})
		} else {
			ctx.JSON(status, map[string]interface{}{
				"ok":     true,
				"status": status,
				"data":   data,
			})
		}
	}

	var err error

	_, err = ctx.Repo.Repository.GetCollaborators()
	if err != nil {
		send(500, nil, err)
		return
	}

	if ctx.HasError() {
		send(400, nil, errors.New(ctx.Flash.ErrorMsg))
		return
	}

	canWrite, err := models.HasAccess(ctx.User, ctx.Repo.Repository, models.ACCESS_MODE_WRITE)
	if err != nil {
		send(500, nil, err)
		return
	}

	if !canWrite {
		send(401, nil, errors.New(ctx.Tr("wiki.rights")))
		return
	}

	p := &models.WikiPage {
		Title:   form.Title,
		Content: form.Content,
		Repo:    ctx.Repo.Repository.WikiRepo,
	}

	err = models.NewWikiPage(p, ctx.Repo.Repository, ctx.User)
	if err != nil {
		send(500, nil, err)
		return
	}

	send(200, fmt.Sprintf("%s/wiki/%s", ctx.Repo.RepoLink, p.Alias), nil)
}

func WikiPageList(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}

	wikiRepoPath, err := wr.WikiRepoPath()
	if err != nil {
		ctx.Handle(500, "wiki.WikiPageList", err)
	}

	filelist, err := filepath.Glob(wikiRepoPath + "/*.md")
	if err != nil {
		ctx.Handle(500, "wiki.WikiPageList", err)
	}

	pagelist := make([]*models.WikiPage, 0)
	repoPath, err := wr.RepoPath()
	if err != nil {
		ctx.Handle(500, "wiki.WikiPageList", err)
	}

	wikiGitRepo, err := git.OpenRepository(repoPath)
	if err != nil {
		ctx.Handle(500, "wiki.WikiPageList", err)
	}
	for _, p := range filelist {
		/**
		 * A little piece of ugly code
		 **/
		commits, err := wikiGitRepo.CommitsByFileAndRange("master", filepath.Base(p), 0)
		if err != nil {
			ctx.Handle(500, "wiki.WikiPageList", err)
		}

		page := strings.Split(filepath.Base(p), ".")
		ptitle := strings.Replace(page[0], "-", " ", -1)
		pagelist = append(pagelist, &models.WikiPage{
			Alias:  page[0],
			Title:  strings.Title(ptitle),
			Repo:   wr,
			Commit: commits.Front().Value.(*git.Commit),
		})
	}

	ctx.Data["Pagelist"] = pagelist
	ctx.HTML(200, WIKI_PAGELIST)
}

func WikiGit(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}

	ctx.Data["WikiRepo"] = wr
	ctx.Data["PageIsGit"] = true
	ctx.HTML(200, WIKI_GIT)
}

func WikiPageRemove(ctx *middleware.Context) {
	send := func(status int, err error) {
		if err != nil {
			log.Error(4, "wiki.WikiPageRemove(?): %s", err)

			ctx.JSON(status, map[string]interface{}{
				"ok":     false,
				"status": status,
				"error":  err.Error(),
			})
		} else {
			ctx.JSON(status, map[string]interface{}{
				"ok":     true,
				"status": status,
			})
		}
	}

	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		send(404, errors.New("There is no wiki repo to delete page"))
		return
	}

	fmt.Println(ctx.User)
	canWrite, err := models.HasAccess(ctx.User, ctx.Repo.Repository, models.ACCESS_MODE_WRITE)
	if err != nil || !canWrite {
		send(401, errors.New(ctx.Tr("wiki.rights")))
		return
	}

	p := &models.WikiPage {
		Alias: ctx.Params(":slug"),
		Title: strings.Title(ctx.Params(":slug")),
		Repo:  ctx.Repo.Repository.WikiRepo,
	}

	if err := p.Delete(ctx.User); err != nil {
		send(500, err)
	}

	send(200, nil)
}