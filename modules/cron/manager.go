// Copyright 2014-2015 The Gogs Authors. All rights reserved.
// Copyright 2015 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cron

import (
	"fmt"

	"github.com/go-gitea/gitea/models"
	"github.com/go-gitea/gitea/modules/setting"
)

var c = New()

func NewCronContext() {
	c.AddFunc("Update mirrors", "@every 1h", models.MirrorUpdate)
	c.AddFunc("Update wiki repos", "@every 10m", models.WikiUpdate)
	c.AddFunc("Deliver hooks", fmt.Sprintf("@every %ds", setting.Webhook.TaskInterval), models.DeliverHooks)
	if setting.Git.Fsck.Enable {
		c.AddFunc("Repository health check", fmt.Sprintf("@every %dh", setting.Git.Fsck.Interval), models.GitFsck)
	}
	c.Start()
}

func ListEntries() []*Entry {
	return c.Entries()
}
