#!/usr/bin/env python3
"""Validate the checked-in z0-z6 map tile pyramid."""

from __future__ import annotations

import argparse
from pathlib import Path


PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"
MAX_ZOOM = 6


def expected_tiles() -> set[Path]:
    return {
        Path(str(z), str(x), f"{y}.png")
        for z in range(MAX_ZOOM + 1)
        for x in range(2**z)
        for y in range(2**z)
    }


def verify_map(root: Path) -> None:
    expected = expected_tiles()
    actual = {path.relative_to(root) for path in root.rglob("*.png")}
    missing = sorted(expected - actual)
    extra = sorted(actual - expected)

    if missing or extra:
        details = []
        if missing:
            details.append(f"missing={len(missing)} (first: {missing[0]})")
        if extra:
            details.append(f"extra={len(extra)} (first: {extra[0]})")
        raise SystemExit("invalid map tile layout: " + ", ".join(details))

    non_png = []
    lfs_pointers = []
    for relative in sorted(actual):
        path = root / relative
        prefix = path.read_bytes()[:64]
        if prefix.startswith(b"version https://git-lfs.github.com/spec/v1"):
            lfs_pointers.append(relative)
        elif not prefix.startswith(PNG_SIGNATURE):
            non_png.append(relative)

    if lfs_pointers:
        raise SystemExit(
            f"Git LFS objects were not pulled (first pointer: {lfs_pointers[0]})"
        )
    if non_png:
        raise SystemExit(f"invalid PNG data (first file: {non_png[0]})")

    print(f"verified {len(actual)} map tiles in {root}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("root", nargs="?", type=Path, default=Path("map"))
    args = parser.parse_args()
    verify_map(args.root.resolve())


if __name__ == "__main__":
    main()
