#!/bin/sh

if [ ! -d "$GITEA_CUSTOM_CONF_PATH" ]; then
    mkdir -p $GITEA_CUSTOM_CONF_PATH

echo "
{{ CONFIG }}
" >> $GITEA_CUSTOM_CONF

fi

exec "$@"
