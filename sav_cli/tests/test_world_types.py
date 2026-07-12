import sys
import unittest
from pathlib import Path


sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from world_types import Pal


class PalOutputTests(unittest.TestCase):
    def test_palworld_1_0_hp_talent_is_exposed_as_melee(self):
        pal = Pal(
            {
                "OwnerPlayerUId": {
                    "value": "00000001-0000-0000-0000-000000000000"
                },
                "Talent_HP": {"value": {"value": 73}},
            },
            real_date_time_ticks=0,
            filetime=0,
        ).to_dict()

        self.assertEqual(pal["melee"], 73)

if __name__ == "__main__":
    unittest.main()
