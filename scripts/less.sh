#!/bin/sh
echo "compiling LESS Files"
lessc ../public/ng/less/gitea.less ../public/ng/css/gitea.css
lessc ../public/ng/less/ui.less ../public/ng/css/ui.css
echo "done"
