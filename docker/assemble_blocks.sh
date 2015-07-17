#!/bin/bash

blocks_dir=blocks
docker_dir=docker
template_dir=templates

docker_file=Dockerfile

gitea_config_file=conf.tmp
gitea_config=config
gitea_init_file=$docker_dir/init_gitea.sh

fig_file=fig.yml
fig_config=fig

gitea_init_template=$template_dir/init_gitea.sh.tpl

if [ "$#" == 0 ]; then
    blocks=`ls $blocks_dir`
    if [ -z "$blocks" ]; then
        echo "No Blocks available in $blocks_dir"
    else
        echo "Available Blocks:"
        for block in $blocks; do
            echo "    $block"
        done
    fi
    exit 0
fi

for file in $gitea_config_file $fig_file; do
    if [ -e $file ]; then
        echo "Deleting $file"
        rm $file
    fi
done

for dir in $@; do
    current_dir=$blocks_dir/$dir
    if [ ! -d "$current_dir" ]; then
        echo "$current_dir is not a directory"
        exit 1
    fi

    if [ -e $current_dir/$docker_file ]; then
        echo "Copying $current_dir/$docker_file to $docker_dir/$docker_file"
        cp $current_dir/$docker_file $docker_dir/$docker_file
    fi

    if [ -e $current_dir/$gitea_config ]; then
        echo "Adding $current_dir/$gitea_config to $gitea_config_file"
        cat $current_dir/$gitea_config >> $gitea_config_file
        echo "" >> $gitea_config_file
    fi

    if [ -e $current_dir/$fig_config ]; then
        echo "Adding $current_dir/$fig_config to $fig_file"
        cat $current_dir/fig >> $fig_file
        echo "" >> $fig_file
    fi
done

echo "Creating $gitea_init_file"
sed "/{{ CONFIG }}/{
r $gitea_config_file
d
}" $gitea_init_template > $gitea_init_file

if [ -e $gitea_config_file ]; then
    echo "Removing temporary Gitea config"
    rm $gitea_config_file
fi
