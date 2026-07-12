#!/usr/bin/env python3
"""Build and validate sav_cli_oss against the local Palworld 1.0 fixtures."""

import json
import subprocess
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SAVS_DIR = ROOT / "savs"
OUTPUT_DIR = SAVS_DIR / ".test-output"
IMAGE = "pst-sav-cli-oss-test"
FIXTURES = (
    "2026.07.12-13.59.00",
    "2026.07.12-18.35.48",
)


def validate_output(path: Path) -> dict[str, object]:
    with path.open(encoding="utf-8") as output_file:
        data = json.load(output_file)

    players = data.get("players", [])
    guilds = data.get("guilds", [])
    pals = [pal for player in players for pal in player.get("pals", [])]
    items = [
        item
        for player in players
        for container in (player.get("items") or {}).values()
        for item in container
    ]

    assert players, "players must not be empty"
    assert guilds, "guilds must not be empty"
    assert pals, "pals must not be empty"
    assert items, "items must not be empty"
    assert all(
        {"player_uid", "level", "pals", "items"} <= player.keys()
        for player in players
    ), "player output contract is incomplete"
    assert all(
        {"melee", "ranged", "defense", "type"} <= pal.keys() for pal in pals
    ), "pal output contract is incomplete"

    hp_talents = [pal["melee"] for pal in pals]
    hp_talent_nonzero = sum(value != 0 for value in hp_talents)
    assert hp_talent_nonzero, "Palworld 1.0 Talent_HP values were not mapped"

    return {
        "players": len(players),
        "guilds": len(guilds),
        "pals": len(pals),
        "items": len(items),
        "hp_iv_nonzero": hp_talent_nonzero,
        "hp_iv_range": [min(hp_talents), max(hp_talents)],
    }


def main() -> None:
    for fixture in FIXTURES:
        level_sav = SAVS_DIR / fixture / "Level.sav"
        players_dir = SAVS_DIR / fixture / "Players"
        if not level_sav.is_file() or not players_dir.is_dir():
            raise SystemExit(f"Missing fixture: {fixture}")

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    subprocess.run(
        [
            "docker",
            "build",
            "--progress=plain",
            "-f",
            "docker/Dockerfile.oss",
            "-t",
            IMAGE,
            ".",
        ],
        cwd=ROOT,
        check=True,
    )

    for fixture in FIXTURES:
        output_path = OUTPUT_DIR / f"{fixture}.json"
        log_path = OUTPUT_DIR / f"{fixture}.log"
        with log_path.open("w", encoding="utf-8") as log_file:
            subprocess.run(
                [
                    "docker",
                    "run",
                    "--rm",
                    "-v",
                    f"{SAVS_DIR / fixture}:/save:ro",
                    "-v",
                    f"{OUTPUT_DIR}:/out",
                    IMAGE,
                    "/app/sav_cli_oss/sav_cli",
                    "-f",
                    "/save/Level.sav",
                    "-o",
                    f"/out/{fixture}.json",
                ],
                check=True,
                stdout=log_file,
                stderr=subprocess.STDOUT,
            )
        print(f"{fixture}: {validate_output(output_path)}")


if __name__ == "__main__":
    main()
