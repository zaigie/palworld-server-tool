#!/bin/bash
# https://github.com/thijsvanloef/palworld-server-docker/blob/main/scripts/init.sh

if [[ ! "${PUID}" -eq 0 ]] && [[ ! "${PGID}" -eq 0 ]]; then
    printf "\e[0;32m*****EXECUTING USERMOD*****\e[0m\n"
    usermod -o -u "${PUID}" steam
    groupmod -o -g "${PGID}" steam
else
    printf "\033[31mRunning as root is not supported, please fix your PUID and PGID!\n"
    exit 1
fi

mkdir -p /palworld/backups
chown -R steam:steam /palworld /home/steam/

term_handler() {
    if [ "${RCON_ENABLED,,}" = true ]; then
        rcon-cli save
        rcon-cli "shutdown 1"
    else # Does not save
        kill -SIGTERM "$(pidof PalServer-Linux-Test)"
    fi
    tail --pid="$killpid" -f 2>/dev/null
}

trap 'term_handler' SIGTERM

su steam -c ./start.sh &
# Process ID of su

LEVEL_SAV_PATH=$(find /palworld/Pal/Saved/SaveGames/0 -name "Level.sav")

if [ -f "$LEVEL_SAV_PATH" ]; then
    echo "Starting pst-agent..."
    nohup pst-agent -f "$LEVEL_SAV_PATH" > /var/log/pst-agent.log 2>&1 &
else
    echo "Level.sav not found, pst-agent will not start."
fi

killpid="$!"
wait "$killpid"

mapfile -t backup_pids < <(pgrep backup)
if [ "${#backup_pids[@]}" -ne 0 ]; then
    echo "Waiting for backup to finish"
    for pid in "${backup_pids[@]}"; do
        tail --pid="$pid" -f 2>/dev/null
    done
fi