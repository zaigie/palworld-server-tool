import importlib.util
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock


SAV_CLI_DIR = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location(
    "sav_cli_structurer", SAV_CLI_DIR / "structurer.py"
)
structurer = importlib.util.module_from_spec(SPEC)
with mock.patch.dict(
    sys.modules,
    {
        "palsav.core": SimpleNamespace(decompress_sav_to_gvas=mock.Mock()),
        "palsav.gvas": SimpleNamespace(GvasFile=mock.Mock()),
        "palsav.paltypes": SimpleNamespace(
            PALWORLD_TYPE_HINTS={},
            PALWORLD_CUSTOM_PROPERTIES={},
        ),
        "world_types": SimpleNamespace(
            Player=mock.Mock(),
            Pal=mock.Mock(),
            Guild=mock.Mock(),
            BaseCamp=mock.Mock(),
        ),
        "logger": SimpleNamespace(log=mock.Mock()),
    },
):
    SPEC.loader.exec_module(structurer)


class StructurerLoggingTests(unittest.TestCase):
    def test_corrupted_player_save_is_counted_as_a_warning(self):
        with (
            mock.patch.object(structurer.os.path, "exists", return_value=True),
            mock.patch.object(
                structurer,
                "_read_gvas",
                side_effect=ValueError("invalid save header"),
            ),
            mock.patch.object(structurer, "log") as log,
        ):
            items, has_warning = structurer.getPlayerItems(
                "player-1", "/save/Players", {}
            )

        self.assertTrue(has_warning)
        self.assertTrue(all(value == [] for value in items.values()))
        log.assert_called_once_with(
            "Skipped corrupted player save: PLAYER1.sav: "
            "ValueError: invalid save header",
            "WARNING",
        )


if __name__ == "__main__":
    unittest.main()
