// Copyright 2014-2015 The Gogs Authors. All rights reserved.
// Copyright 2015 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

type ObjectType string

const (
	COMMIT ObjectType = "commit"
	TREE   ObjectType = "tree"
	BLOB   ObjectType = "blob"
	TAG    ObjectType = "tag"
)
