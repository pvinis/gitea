package repo

import (
    "github.com/go-gitea/gitea/modules/base"
    "github.com/go-gitea/gitea/modules/middleware"
)

const (
    WIKI base.TplName = "repo/wiki/home"
)

func Wiki(ctx *middleware.Context) {
    ctx.Data["Title"] = ctx.Repo.Repository.Name
//
//    wr, err := ctx.Repo.Repository.WikiRepo()
//    if err != nil {
//        ctx.Handle(500, "GetBranches", err)
//        return
//    }
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