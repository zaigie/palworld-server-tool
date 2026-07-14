# SPDX-License-Identifier: Apache-2.0
# Derived from zaigie/palworld-server-tool sav_cli @ fb45624 (Apache-2.0).
# Runtime deps (palsav-flex/palooz/ooz) are GPL-3.0-or-later, so a Docker image
# built from the root Dockerfile includes these runtime components.
"""Consistent console logging for sav_cli and its parser dependencies."""

import logging


def configure_logging(verbose=False):
    logging.basicConfig(
        level=logging.INFO,
        format="[SAV-CLI] %(asctime)s | %(levelname)s | %(message)s",
        datefmt="%Y/%m/%d - %H:%M:%S",
        force=True,
    )
    logging.getLogger("palsav").setLevel(logging.INFO if verbose else logging.WARNING)


configure_logging()


_LEVELS = {
    "DEBUG": logging.debug,
    "INFO": logging.info,
    "WARNING": logging.warning,
    "ERROR": logging.error,
    "CRITICAL": logging.critical,
}


def log(text, level="INFO"):
    _LEVELS.get(level.upper(), logging.info)(text)
