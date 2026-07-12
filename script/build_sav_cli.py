#!/usr/bin/env python3
"""Build a native, single-file sav_cli executable for the host platform."""

from __future__ import annotations

import argparse
import os
from pathlib import Path
import subprocess
import sys
import tempfile


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PST_TOOLS_REF = "8cb429ae3b1460a6a6a0c31c9964ca8cedb65cc5"
PST_TOOLS_REPOSITORY = "https://github.com/deafdudecomputers/PalworldSaveTools.git"


def run(command: list[str], *, cwd: Path = ROOT) -> None:
    print("+", subprocess.list2cmdline(command), flush=True)
    subprocess.run(command, cwd=cwd, check=True)


def patch_palooz_for_windows(setup_path: Path) -> None:
    source = setup_path.read_text(encoding="utf-8")
    old = """import os
from setuptools import setup, Extension
"""
    new = """import os
import sys
from setuptools import setup, Extension
"""
    flags_old = "extra_compile_args = ['-O3', '-flto', '-fno-exceptions', '-fno-rtti', '-ffast-math', '-fno-strict-aliasing']"
    flags_new = """if sys.platform == 'win32':
    extra_compile_args = ['/O2', '/fp:fast', '/GR-']
else:
    extra_compile_args = ['-O3', '-flto', '-fno-exceptions', '-fno-rtti', '-ffast-math', '-fno-strict-aliasing']"""
    if old not in source or flags_old not in source:
        raise RuntimeError("the pinned palooz setup.py no longer matches the expected patch")
    setup_path.write_text(
        source.replace(old, new, 1).replace(flags_old, flags_new, 1),
        encoding="utf-8",
    )


def build(output_dir: Path, pst_tools_ref: str) -> Path:
    output_dir = output_dir.resolve()
    output_dir.mkdir(parents=True, exist_ok=True)

    run(
        [
            sys.executable,
            "-m",
            "pip",
            "install",
            "--disable-pip-version-check",
            "-r",
            str(ROOT / "docker" / "sav-cli-requirements.lock"),
            "-r",
            str(ROOT / "docker" / "sav-cli-build-requirements.lock"),
        ]
    )

    with tempfile.TemporaryDirectory(prefix="sav-cli-build-") as temp:
        temp_dir = Path(temp)
        tools_dir = temp_dir / "PalworldSaveTools"
        run(["git", "clone", "--filter=blob:none", "--no-checkout", PST_TOOLS_REPOSITORY, str(tools_dir)])
        run(["git", "checkout", "--detach", pst_tools_ref], cwd=tools_dir)

        palsav_dir = tools_dir / "src" / "palsav"
        palooz_dir = palsav_dir / "palooz"
        if sys.platform == "win32":
            patch_palooz_for_windows(palooz_dir / "setup.py")

        run(
            [
                sys.executable,
                "-m",
                "pip",
                "install",
                "--disable-pip-version-check",
                "--no-build-isolation",
                "--no-deps",
                str(palooz_dir),
            ]
        )
        run(
            [
                sys.executable,
                "-m",
                "pip",
                "install",
                "--disable-pip-version-check",
                "--no-build-isolation",
                "--no-deps",
                str(palsav_dir),
            ]
        )

        work_dir = temp_dir / "pyinstaller"
        run(
            [
                sys.executable,
                "-m",
                "PyInstaller",
                "--clean",
                "--noconfirm",
                "--onefile",
                "--name",
                "sav_cli",
                "--paths",
                str(ROOT / "sav_cli"),
                "--collect-all",
                "palsav",
                "--hidden-import",
                "palooz",
                "--distpath",
                str(output_dir),
                "--workpath",
                str(work_dir),
                "--specpath",
                str(temp_dir),
                str(ROOT / "sav_cli" / "sav_cli.py"),
            ]
        )

    executable = output_dir / ("sav_cli.exe" if os.name == "nt" else "sav_cli")
    if not executable.is_file():
        raise RuntimeError(f"PyInstaller did not create {executable}")
    executable.chmod(executable.stat().st_mode | 0o111)
    run([str(executable), "--help"])
    return executable


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, default=ROOT / "dist" / "sav-cli")
    parser.add_argument(
        "--pst-tools-ref",
        default=os.environ.get("PST_TOOLS_REF", DEFAULT_PST_TOOLS_REF),
    )
    args = parser.parse_args()
    executable = build(args.output, args.pst_tools_ref)
    print(f"built {executable}")


if __name__ == "__main__":
    main()
