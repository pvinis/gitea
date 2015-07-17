// +build !pam

// Copyright 2014-2015 The Gogs Authors. All rights reserved.
// Copyright 2015 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pam

import (
	"errors"
)

func PAMAuth(serviceName, userName, passwd string) error {
	return errors.New("PAM not supported")
}
