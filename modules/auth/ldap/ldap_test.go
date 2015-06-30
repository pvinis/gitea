// Copyright github.com/juju2013. All rights reserved.
// Copyright 2014-2015 The Gogs Authors. All rights reserved.
// Copyright 2015 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ldap

// import (
// 	"fmt"
// 	"testing"
// )

// var ldapServer = "ldap.itd.umich.edu"
// var ldapPort = 389
// var baseDN = "dc=umich,dc=edu"
// var filter = []string{
// 	"(cn=cis-fac)",
// 	"(&(objectclass=rfc822mailgroup)(cn=*Computer*))",
// 	"(&(objectclass=rfc822mailgroup)(cn=*Mathematics*))"}
// var attributes = []string{
// 	"cn",
// 	"description"}
// var msadsaformat = ""

// func TestLDAP(t *testing.T) {
// 	AddSource("test", ldapServer, ldapPort, baseDN, attributes, filter, msadsaformat)
// 	user, err := LoginUserLdap("xiaolunwen", "")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	fmt.Println(user)
// }
