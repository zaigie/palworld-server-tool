# SPDX-License-Identifier: Apache-2.0
# Derived from zaigie/palworld-server-tool sav_cli @ fb45624 (Apache-2.0).
# Runtime deps (palsav-flex/palooz/ooz) are GPL-3.0-or-later, so a Docker image
# built from the root Dockerfile includes these runtime components.
"""``sav_cli`` for Palworld 1.0.

Parses a Level.sav (and per-player saves under ``Players/``) and either writes a
``{"players": [...], "guilds": [...]}`` JSON document (default), or PUTs the
players/guilds to a palworld-server-tool backend when ``--request`` is given.

Usage:
    python sav_cli.py -f /path/to/Level.sav -o structure.json
    python sav_cli.py -f /path/to/Level.sav --request http://host/api/ --token TOKEN
"""

import os
import sys
import json
import shutil
import time
import argparse
from urllib.parse import urljoin

from structurer import convert_sav, structure_player, structure_guild
from logger import log


def _http_error_details(response):
    reason = str(response.reason).strip() if response.reason else ""
    status = f"HTTP {response.status_code}"
    if reason:
        status = f"{status} {reason}"

    body = response.text.strip() or "<empty response body>"
    return f"{status}; response body: {body}"


def main():
    start = time.time()
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", "-f", help="File to convert", type=str, default="Level.sav")
    parser.add_argument("--clear", "-c", help="Clear input file", action="store_true")
    parser.add_argument("--output", "-o", help="Output file", type=str, default="structure.json")
    parser.add_argument("--request", "-r", help="Request base URL", type=str, default="")
    parser.add_argument("--token", "-t", help="Request token", type=str, default="")
    args = parser.parse_args()

    output = args.output
    if args.request == "" and not output.endswith(".json"):
        output = output + ".json"

    if not os.path.exists(args.file):
        log(f"File not exists: {args.file}", "ERROR")
        sys.exit(1)

    convert_sav(args.file)
    filetime = os.stat(args.file).st_mtime

    dir_path = os.path.join(os.path.dirname(args.file), "Players")

    players = structure_player(dir_path, filetime=filetime)
    guilds = structure_guild(filetime)

    # Fill save_last_online from the player's guild membership record.
    for player in players:
        for guild in guilds:
            for guild_player in guild["players"]:
                if player["player_uid"] == guild_player["player_uid"]:
                    player["save_last_online"] = guild_player["last_online"]
                    break

    if args.request == "":
        with open(output, "w", encoding="utf-8") as f:
            json.dump(
                {"players": players, "guilds": guilds}, f, indent=4, ensure_ascii=False
            )
        log(f"Players: {len(players)}")
        log(f"Guilds: {len(guilds)}")
        log(f"Wrote {output}")
    else:
        import requests

        player_url = urljoin(args.request, "player")
        guild_url = urljoin(args.request, "guild")
        log(f"Put players to {player_url} with Players: {len(players)}")
        player_res = requests.put(
            player_url,
            headers={"Authorization": f"Bearer {args.token}"},
            json=players,
            timeout=10,
        )
        if player_res.status_code != 200:
            log(
                f"Put Players data error: {_http_error_details(player_res)}",
                "ERROR",
            )

        log(f"Put guilds to {guild_url} with Guilds: {len(guilds)}")
        guild_res = requests.put(
            guild_url,
            headers={"Authorization": f"Bearer {args.token}"},
            json=guilds,
            timeout=10,
        )
        if guild_res.status_code != 200:
            log(
                f"Put Guilds data error: {_http_error_details(guild_res)}",
                "ERROR",
            )

    try:
        if args.clear:
            os.remove(args.file)
            if os.path.exists(dir_path):
                shutil.rmtree(dir_path)
    except FileNotFoundError:
        pass

    log(f"Done in {round(time.time() - start, 3)}s")


if __name__ == "__main__":
    main()
