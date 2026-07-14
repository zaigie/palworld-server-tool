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

from logger import configure_logging, log
from structurer import convert_sav, structure_player, structure_guild

MAX_ERROR_BODY_LENGTH = 512


def _compact_error_text(text):
    compact = " ".join(str(text).split()) or "<empty response body>"
    if len(compact) > MAX_ERROR_BODY_LENGTH:
        return f"{compact[:MAX_ERROR_BODY_LENGTH]}... <truncated>"
    return compact


def _http_error_details(response):
    status = _http_status(response)

    try:
        payload = json.loads(response.text)
    except (TypeError, json.JSONDecodeError):
        payload = None

    if isinstance(payload, dict) and payload.get("error"):
        return f"{status}; error: {_compact_error_text(payload['error'])}"
    return f"{status}; response body: {_compact_error_text(response.text)}"


def _http_status(response):
    reason = str(response.reason).strip() if response.reason else ""
    return f"HTTP {response.status_code}{f' {reason}' if reason else ''}"


def _elapsed(started):
    return f"{time.perf_counter() - started:.2f}s"


def _sync_resource(requests, resource, url, token, records):
    started = time.perf_counter()
    try:
        response = requests.put(
            url,
            headers={"Authorization": f"Bearer {token}"},
            json=records,
            timeout=10,
        )
    except requests.RequestException as error:
        log(
            f"Failed to sync {resource}={len(records)}: "
            f"{type(error).__name__}: {error} ({_elapsed(started)})",
            "ERROR",
        )
        return False

    if response.status_code != 200:
        log(
            f"Failed to sync {resource}={len(records)}: "
            f"{_http_error_details(response)} ({_elapsed(started)})",
            "ERROR",
        )
        return False

    log(
        f"Synced {resource}={len(records)}: {_http_status(response)} "
        f"({_elapsed(started)})"
    )
    return True


def main():
    start = time.perf_counter()
    failed_requests = 0
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", "-f", help="File to convert", type=str, default="Level.sav")
    parser.add_argument("--clear", "-c", help="Clear input file", action="store_true")
    parser.add_argument("--output", "-o", help="Output file", type=str, default="structure.json")
    parser.add_argument("--request", "-r", help="Request base URL", type=str, default="")
    parser.add_argument("--token", "-t", help="Request token", type=str, default="")
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Show detailed parser and decompression logs",
    )
    args = parser.parse_args()
    configure_logging(verbose=args.verbose)

    output = args.output
    if args.request == "" and not output.endswith(".json"):
        output = output + ".json"

    if not os.path.exists(args.file):
        log(f"File not exists: {args.file}", "ERROR")
        return 1

    log(f"Decoding save: {os.path.basename(args.file)}")
    phase_start = time.perf_counter()
    convert_sav(args.file)
    log(f"Decoded save ({_elapsed(phase_start)})")
    filetime = os.stat(args.file).st_mtime

    dir_path = os.path.join(os.path.dirname(args.file), "Players")

    log("Structuring save data")
    phase_start = time.perf_counter()
    players, player_save_warnings = structure_player(dir_path, filetime=filetime)
    guilds = structure_guild(filetime)

    # Fill save_last_online from the player's guild membership record.
    for player in players:
        for guild in guilds:
            for guild_player in guild["players"]:
                if player["player_uid"] == guild_player["player_uid"]:
                    player["save_last_online"] = guild_player["last_online"]
                    break

    pal_count = sum(len(player.get("pals", [])) for player in players)
    base_camp_count = sum(len(guild.get("base_camp", [])) for guild in guilds)
    warning_summary = (
        f", player_save_warnings={player_save_warnings}"
        if player_save_warnings
        else ""
    )
    log(
        "Structured save: "
        f"players={len(players)}, pals={pal_count}, guilds={len(guilds)}, "
        f"base_camps={base_camp_count}{warning_summary} ({_elapsed(phase_start)})"
    )

    if args.request == "":
        with open(output, "w", encoding="utf-8") as f:
            json.dump(
                {"players": players, "guilds": guilds}, f, indent=4, ensure_ascii=False
            )
        log(f"Wrote structured save: {output}")
    else:
        import requests

        player_url = urljoin(args.request, "player")
        guild_url = urljoin(args.request, "guild")
        log(f"Syncing save data: {args.request}")
        if not _sync_resource(
            requests,
            "players",
            player_url,
            args.token,
            players,
        ):
            failed_requests += 1

        if not _sync_resource(
            requests,
            "guilds",
            guild_url,
            args.token,
            guilds,
        ):
            failed_requests += 1

    if failed_requests:
        warning_summary = (
            f", warnings={player_save_warnings}" if player_save_warnings else ""
        )
        log(
            f"Save sync failed: requests_failed={failed_requests}{warning_summary} "
            f"({_elapsed(start)})",
            "ERROR",
        )
        return 1

    try:
        if args.clear:
            os.remove(args.file)
            if os.path.exists(dir_path):
                shutil.rmtree(dir_path)
    except FileNotFoundError:
        pass

    operation = "Save sync" if args.request else "Save export"
    warning_summary = (
        f" with warnings={player_save_warnings}" if player_save_warnings else ""
    )
    log(f"{operation} completed{warning_summary} ({_elapsed(start)})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
