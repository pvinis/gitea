#!/bin/sh

PID="log/supervisord.pid"
CONF="etc/supervisord.conf"

EXEPATH='/usr/bin/gitea_start'
if [ ! -f $EXEPATH ]; then
    gitea_scripts_path=$(cd `dirname $0`; pwd)
    echo $gitea_scripts_path
    sudo ln -s $gitea_scripts_path'/start.sh' /usr/bin/gitea_start
fi

LOGDIR="log"
if [ ! -d $LOGDIR ]; then
    mkdir $LOGDIR
fi

stop() {
    if [ -f $PID ]; then
        kill `cat -- $PID`
        rm -f -- $PID
        echo "stopped"
    fi
}

start() {
    echo "starting"
    if [ ! -f $PID ]; then
        supervisord -c $CONF
        echo "started"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
esac
