#!/usr/bin/env python3
"""Build Go binaries and assemble one platform's release files."""

from __future__ import annotations

import argparse
import os
from pathlib import Path
import shutil
import subprocess
import tarfile
import tempfile
import zipfile

from verify_map import verify_map


ROOT = Path(__file__).resolve().parents[1]
ARCH_NAMES = {"amd64": "x86_64", "arm64": "aarch64"}


def run(command: list[str], *, env: dict[str, str]) -> None:
    print("+", subprocess.list2cmdline(command), flush=True)
    subprocess.run(command, cwd=ROOT, env=env, check=True)


def build_go(
    source: str,
    destination: Path,
    version: str | None,
    goos: str,
    goarch: str,
) -> None:
    destination.parent.mkdir(parents=True, exist_ok=True)
    env = os.environ.copy()
    env.update({"GOOS": goos, "GOARCH": goarch, "CGO_ENABLED": "0"})
    ldflags = "-s -w"
    if version is not None:
        ldflags += f" -X main.version={version}"
    run(
        ["go", "build", "-trimpath", "-ldflags", ldflags, "-o", str(destination), source],
        env=env,
    )


def archive_directory(source: Path, destination: Path, windows: bool) -> None:
    if windows:
        with zipfile.ZipFile(destination, "w", compression=zipfile.ZIP_DEFLATED) as archive:
            for path in sorted(source.iterdir()):
                archive.write(path, path.name)
    else:
        with tarfile.open(destination, "w:gz") as archive:
            for path in sorted(source.iterdir()):
                archive.add(path, arcname=path.name)


def package(version: str, goos: str, goarch: str, sav_cli: Path, output: Path) -> list[Path]:
    if goarch not in ARCH_NAMES:
        raise ValueError(f"unsupported architecture: {goarch}")
    for required in (ROOT / "assets", ROOT / "index.html", ROOT / "pal-conf.html"):
        if not required.exists():
            raise FileNotFoundError(f"frontend build artifact is missing: {required}")
    verify_map(ROOT / "map")

    windows = goos == "windows"
    executable_suffix = ".exe" if windows else ""
    platform = f"{goos}_{ARCH_NAMES[goarch]}"
    output.mkdir(parents=True, exist_ok=True)

    with tempfile.TemporaryDirectory(prefix=f"pst-{platform}-") as temp:
        stage = Path(temp)
        pst = stage / f"pst{executable_suffix}"
        build_go("main.go", pst, version, goos, goarch)

        sav_name = f"sav_cli{executable_suffix}"
        shutil.copy2(sav_cli, stage / sav_name)
        if not windows:
            (stage / sav_name).chmod((stage / sav_name).stat().st_mode | 0o111)
            pst.chmod(pst.stat().st_mode | 0o111)
        shutil.copy2(ROOT / "LICENSE", stage / "LICENSE")
        shutil.copy2(ROOT / "NOTICE", stage / "NOTICE")
        if windows:
            shutil.copy2(ROOT / "script" / "start.bat", stage / "start.bat")

        extension = "zip" if windows else "tar.gz"
        package_path = output / f"pst_{version}_{platform}.{extension}"
        archive_directory(stage, package_path, windows)

    agent_path = output / f"pst-agent_{version}_{platform}{executable_suffix}"
    build_go("cmd/pst-agent/main.go", agent_path, None, goos, goarch)
    if not windows:
        agent_path.chmod(agent_path.stat().st_mode | 0o111)
    return [package_path, agent_path]


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--version", required=True)
    parser.add_argument("--goos", required=True, choices=("linux", "windows", "darwin"))
    parser.add_argument("--goarch", required=True, choices=tuple(ARCH_NAMES))
    parser.add_argument("--sav-cli", required=True, type=Path)
    parser.add_argument("--output", type=Path, default=ROOT / "dist" / "release")
    args = parser.parse_args()
    for artifact in package(
        args.version, args.goos, args.goarch, args.sav_cli.resolve(), args.output.resolve()
    ):
        print(f"created {artifact}")


if __name__ == "__main__":
    main()
