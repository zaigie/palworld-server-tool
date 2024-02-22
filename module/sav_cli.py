import os
import json
import time
import requests
import argparse
from urllib.parse import urljoin

from structurer import convert_sav, structure_player, structure_guild
from logger import log


if __name__ == "__main__":
    start = time.time()
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--file", "-f", help="File to convert", type=str, default="Level.sav"
    )
    parser.add_argument("--clear", "-c", help="Clear input file", action="store_true")
    parser.add_argument(
        "--output", "-o", help="Output file", type=str, default="structure.json"
    )
    parser.add_argument("--request", "-r", help="Request", type=str, default="")
    parser.add_argument("--token", "-t", help="Request token", type=str, default="")
    args = parser.parse_args()

    if args.request == "":
        output = args.output
        if not args.output.endswith(".json"):
            output = args.output + ".json"

    converted = convert_sav(args.file)
    filetime = os.stat(args.file).st_mtime
    players = structure_player(converted)
    guilds = structure_guild(converted, filetime)

    # Add last_online to players
    for player in players:
        for guild in guilds:
            guild_players = guild["players"]
            for guild_player in guild_players:
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
    else:
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
            log(f"Put Players data error: {player_res.text}")

        log(f"Put guilds to {guild_url} with Guilds: {len(guilds)}")
        guild_res = requests.put(
            guild_url,
            headers={"Authorization": f"Bearer {args.token}"},
            json=guilds,
            timeout=10,
        )
        if guild_res.status_code != 200:
            log(f"Put Guilds data error: {guild_res.text}")

    try:
        if args.clear:
            os.remove(args.file)
    except FileNotFoundError:
        pass

    log(f"Done in {round(time.time() - start,3)}s")
