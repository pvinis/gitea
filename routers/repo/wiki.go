package repo

import (
	"errors"
//	"fmt"
	"github.com/go-gitea/gitea/models"
	"github.com/go-gitea/gitea/modules/auth"
	"github.com/go-gitea/gitea/modules/base"
	"github.com/go-gitea/gitea/modules/log"
	"github.com/go-gitea/gitea/modules/middleware"
	"fmt"
)

const (
	WIKI       base.TplName = "repo/wiki/home"
	WIKI_EMPTY base.TplName = "repo/wiki/empty"
	WIKI_ADD   base.TplName = "repo/wiki/add"
)

func Wiki(ctx *middleware.Context) {
	ctx.Data["Title"] = ctx.Repo.Repository.Name

	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.HTML(200, WIKI_EMPTY)
		return
	}
	//    refName := wr.DefaultBranch
	//    if !ctx.Repo.GitRepo.IsBranchExist(refName) {
	//        brs, err := ctx.Repo.GitRepo.GetBranches()
	//        if err != nil {
	//            ctx.Handle(500, "GetBranches", err)
	//            return
	//        }
	//        refName = brs[0]
	//    }
	//    ctx.Repo.Commit, err = ctx.Repo.GitRepo.GetCommitOfBranch(refName)
	//
	//
	//    ctx.Data["Pages"] = wr
	ctx.HTML(200, WIKI)
}

func CreateWikiPage(ctx *middleware.Context) {
	wr := ctx.Repo.Repository.WikiRepo
	if wr == nil {
		ctx.Data["Title"] = "Home"
	}

	_, err := ctx.Repo.Repository.GetCollaborators()
	if err != nil {
		ctx.Handle(400, "wiki.CreateWikiPage", err)
		return
	}

	if !ctx.Repo.IsOwner() {
		ctx.Handle(401, "wiki.CreateWikiPage", errors.New(ctx.Tr("wiki.rights")))
		return
	}

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

	fmt.Println(form)
	if !ctx.Repo.IsOwner() {
		send(401, nil, errors.New(ctx.Tr("wiki.rights")))
		return
	}

	p := &models.WikiPage {
		Name:    form.Title,
		Content: form.Content,
		Repo:    ctx.Repo.Repository.WikiRepo,
	}

	err = models.NewWikiPage(p, ctx.Repo.Repository, ctx.User)
	if err != nil {
		send(500, nil, err)
		return
	}

	send(200, p.Alias, nil)
	// Create repo if needed
	// Create page
	// Redirect to page
}