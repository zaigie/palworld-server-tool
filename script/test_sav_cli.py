#!/usr/bin/env python3
"""Build and validate the current image against local Palworld 1.0 saves."""

import argparse
import json
import subprocess
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SAVS_DIR = ROOT / "savs"
OUTPUT_DIR = SAVS_DIR / ".test-output"
IMAGE_PREFIX = "pst-sav-cli-test"
LOCK_FILE = ROOT / "docker" / "sav-cli-requirements.lock"
PLATFORMS = ("linux/arm64", "linux/amd64")
FIXTURES = (
    "2026.07.12-13.59.00",
    "2026.07.12-18.35.48",
)
BUILD_ONLY_PACKAGES = ("build-base", "git", "python3-dev")
SOURCE_PACKAGES = {"palooz": "0.2.0", "palsav-flex": "0.2.0"}
SAV_CLI_PATH = "/app/sav_cli"


def locked_python_packages() -> dict[str, str]:
    packages = dict(SOURCE_PACKAGES)
    for line in LOCK_FILE.read_text(encoding="utf-8").splitlines():
        requirement = line.strip()
        if not requirement or requirement.startswith("#"):
            continue
        name, version = requirement.split("==", maxsplit=1)
        packages[name] = version
    return packages


def image_name(platform: str) -> str:
    return f"{IMAGE_PREFIX}-{platform.rsplit('/', maxsplit=1)[-1]}"


def validate_image(image: str, platform: str) -> dict[str, object]:
    subprocess.run(
        [
            "docker",
            "run",
            "--platform",
            platform,
            "--rm",
            "--entrypoint",
            "/bin/sh",
            image,
            "-c",
            f'test -x "{SAV_CLI_PATH}" && test -z "$SAVE__DECODE_PATH"',
        ],
        check=True,
    )

    installed_build_packages = [
        package
        for package in BUILD_ONLY_PACKAGES
        if subprocess.run(
            [
                "docker",
                "run",
                "--platform",
                platform,
                "--rm",
                "--entrypoint",
                "apk",
                image,
                "info",
                "-e",
                package,
            ],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        ).returncode
        == 0
    ]
    assert not installed_build_packages, (
        f"build-only packages remain in runtime image: {installed_build_packages}"
    )

    image_size = int(
        subprocess.check_output(
            ["docker", "image", "inspect", image, "--format", "{{.Size}}"],
            text=True,
        ).strip()
    )
    architecture = subprocess.check_output(
        ["docker", "image", "inspect", image, "--format", "{{.Architecture}}"],
        text=True,
    ).strip()
    expected_architecture = platform.rsplit("/", maxsplit=1)[-1]
    assert architecture == expected_architecture, (
        f"expected {expected_architecture} image, got {architecture}"
    )

    expected_packages = locked_python_packages()
    probe = (
        "import importlib.metadata as metadata, json, platform; "
        "normalize=lambda name: name.lower().replace('_', '-'); "
        "print(json.dumps({'python': platform.python_version(), "
        "'packages': {normalize(dist.metadata['Name']): dist.version "
        "for dist in metadata.distributions()}}, "
        "sort_keys=True))"
    )
    runtime = json.loads(
        subprocess.check_output(
            [
                "docker",
                "run",
                "--platform",
                platform,
                "--rm",
                "--entrypoint",
                "/opt/sav-cli-venv/bin/python",
                image,
                "-c",
                probe,
            ],
            text=True,
        )
    )
    assert runtime["python"].startswith("3.12."), (
        f"unexpected Python version: {runtime['python']}"
    )
    assert runtime["packages"] == expected_packages, (
        f"runtime packages differ from lock: {runtime['packages']}"
    )
    subprocess.run(
        [
            "docker",
            "run",
            "--platform",
            platform,
            "--rm",
            "--entrypoint",
            "/opt/sav-cli-venv/bin/python",
            image,
            "-m",
            "pip",
            "check",
        ],
        check=True,
        stdout=subprocess.DEVNULL,
    )

    return {
        "architecture": architecture,
        "size_bytes": image_size,
        "python": runtime["python"],
        "python_packages": len(runtime["packages"]),
        "build_only_packages": [],
    }


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
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--no-cache",
        action="store_true",
        help="Pull pinned inputs and rebuild every Docker layer",
    )
    args = parser.parse_args()

    for fixture in FIXTURES:
        level_sav = SAVS_DIR / fixture / "Level.sav"
        players_dir = SAVS_DIR / fixture / "Players"
        if not level_sav.is_file() or not players_dir.is_dir():
            raise SystemExit(f"Missing fixture: {fixture}")

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    for platform in PLATFORMS:
        image = image_name(platform)
        architecture = platform.rsplit("/", maxsplit=1)[-1]
        build_command = [
            "docker",
            "build",
            "--progress=plain",
            "--platform",
            platform,
        ]
        if args.no_cache:
            build_command.extend(("--no-cache", "--pull"))
        build_command.extend(("-t", image, "."))
        subprocess.run(build_command, cwd=ROOT, check=True)
        print(f"{platform} image: {validate_image(image, platform)}")

        for fixture in FIXTURES:
            output_name = f"{architecture}-{fixture}.json"
            output_path = OUTPUT_DIR / output_name
            log_path = OUTPUT_DIR / f"{architecture}-{fixture}.log"
            with log_path.open("w", encoding="utf-8") as log_file:
                subprocess.run(
                    [
                        "docker",
                        "run",
                        "--platform",
                        platform,
                        "--rm",
                        "-v",
                        f"{SAVS_DIR / fixture}:/save:ro",
                        "-v",
                        f"{OUTPUT_DIR}:/out",
                        image,
                        SAV_CLI_PATH,
                        "-f",
                        "/save/Level.sav",
                        "-o",
                        f"/out/{output_name}",
                    ],
                    check=True,
                    stdout=log_file,
                    stderr=subprocess.STDOUT,
                )
            print(f"{platform} {fixture}: {validate_output(output_path)}")


if __name__ == "__main__":
    main()
