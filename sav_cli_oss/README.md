# sav_cli_oss — open-source Palworld 1.0 save parser

An open-source replacement for palworld-server-tool's closed
`sav_cli`, capable of parsing **Palworld 1.0** saves (tested against
`v1.0.0.100427`). It reads a `Level.sav` plus the per-player saves under
`Players/` and emits the same `{"players": [...], "guilds": [...]}` JSON that
the palworld-server-tool backend consumes — or PUTs it straight to the backend.

## Why this exists

The original `sav_cli` was closed after commit `fb45624`. It was built on the
`palworld_save_tools` library, whose bundled decoders and the tool's own
`skip_decode` optimizations no longer parse a 1.0 `Level.sav` (they crash with
`EOF not reached`).

This version keeps the original structuring logic but rewires it onto
[`palsav`](../../PalworldSaveTools) (the `palsav-flex` package from
PalworldSaveTools), which ships full 1.0 mappings and Oodle (`PlM1`)
decompression, and it does a **full decode** instead of the old skip machinery.

## What changed for 1.0

Verified against a real `v1.0.0.100427` save:

| Field | Old (pre-1.0) | 1.0 |
|-------|---------------|-----|
| Player/Pal HP | `HP` | **`Hp`** (FixedPoint64 at `.value.Value.value`) |
| Item slots | `RawData.value.permission.{type_a,item_static_id,type_b}` | **`RawData.value.{slot_index, count, item.static_id}`** |
| `Level`, talents, ranks | ByteProperty | ByteProperty (number nested at `.value.value`) |
| `MaxHP`, `ShieldMaxHP`, `MaxSP` | present | **not persisted** in 1.0 player saves (default 0) |

Guild (`GroupSaveDataMap`) and base-camp (`BaseCampSaveData`) shapes from
palsav are compatible with the original accessors.

## Setup

Requires Python 3.11+ and a C++ compiler (MSVC on Windows, gcc/clang on
Linux/macOS) to build the native `palooz` Oodle module.

```bash
# 1. Create a venv
python -m venv .venv
# Windows:
.venv\Scripts\activate
# Linux/macOS:
# source .venv/bin/activate

# 2. Build & install the native Oodle module (palooz) from PalworldSaveTools.
#    On Windows, the upstream setup.py uses GCC-only flags; use MSVC-compatible
#    ones (/O2 /fp:fast /GR-). A patched copy is not required if you build on
#    Linux/macOS. See "Building palooz on Windows" below.
pip install <PalworldSaveTools>/src/palsav/palooz

# 3. Install the palsav parser (editable) + JSON dep
pip install orjson
pip install --no-build-isolation -e <PalworldSaveTools>/src/palsav

# 4. (optional) for --request mode
pip install requests
```

### Building palooz on Windows

`palooz`'s upstream `setup.py` passes GCC/Clang flags (`-O3 -flto …`) that MSVC
`cl.exe` rejects. Edit `setup.py` so Windows uses MSVC flags:

```python
if sys.platform == 'win32':
    extra_compile_args = ['/O2', '/fp:fast', '/GR-']
else:
    extra_compile_args = ['-O3', '-flto', '-fno-exceptions', '-fno-rtti', '-ffast-math', '-fno-strict-aliasing']
```

Then `pip install <path>/palooz`. setuptools auto-detects the installed Visual
Studio Build Tools (verified with MSVC 14.44 / VS 2022+).

## Reproducible Docker build

`docker/Dockerfile.oss` pins the multi-architecture base image by digest and
installs Python packages from `docker/sav-cli-requirements.lock`. The palsav
and palooz sources are pinned separately by `PST_TOOLS_REF`; both local source
packages are installed without build isolation or dependency resolution so
their build cannot silently pull newer packages.

Maintainers can rebuild both supported architectures from scratch and parse
the two local fixtures under `savs/` with:

```bash
python3 script/test_sav_cli_oss.py --no-cache
```

The test uses disposable containers, verifies the image architecture, Python
minor version, exact installed package set, dependency consistency, and output
contract. Update the lock file deliberately and rerun this command whenever a
dependency or base image is changed.

## Run

```bash
# Write JSON locally (default)
python sav_cli.py -f /path/to/Level.sav -o structure.json

# PUT to a palworld-server-tool backend
python sav_cli.py -f /path/to/Level.sav --request http://host/api/ --token TOKEN
```

`Players/` is expected next to `Level.sav`; per-player item containers are read
from those saves.

## Output shape

```jsonc
{
  "players": [
    {
      "player_uid": "1234567890",     // decimal of first 8 hex of PlayerUId
      "nickname": "ExamplePlayer",
      "level": 20, "exp": ..., "hp": ..., "max_hp": 0,
      "shield_hp": ..., "shield_max_hp": 0, "max_status_point": 0,
      "status_point": { "最大HP": 0, ... },
      "full_stomach": 74.2,
      "save_last_online": "2026-01-01T00:00:00Z",
      "pals": [ { "level": 11, "type": "Kitsunebi", "gender": "Male",
                  "ranged": 89, "defense": 23, "skills": ["SalePrice_Up_1"], ... } ],
      "items": { "CommonContainerId": [ {"SlotIndex":0,"ItemId":"money","StackCount":12032}, ... ], ... }
    }
  ],
  "guilds": [
    {
      "name": "ExampleGuild", "base_camp_level": 10,
      "admin_player_uid": "1234567890",
      "players": [ {"player_uid":"...","nickname":"...","last_online":"..."} ],
      "base_camp": [ {"id":"...","area":3500.0,"location_x":...,"location_y":...} ]
    }
  ]
}
```

Fields the backend fills itself (`user_id`, `steam_id`, `ip`, `ping`,
`location`, `building_count`, `account_name`) are intentionally omitted.

## Files

- `sav_cli.py` — CLI entrypoint (argparse, JSON-out or PUT).
- `structurer.py` — decode + structure players / pals / guilds / base camps.
- `world_types.py` — flatten decoded property trees into the output shape.
- `logger.py` — minimal `log(text, level)`.

## License

This code is **Apache-2.0** (derived from zaigie/palworld-server-tool `sav_cli`
@ `fb45624`). At runtime it depends on `palsav-flex`, `palooz`, and `ooz`, which
are **GPL-3.0-or-later**; therefore a Docker image built from
`docker/Dockerfile.oss` is a **GPL-3.0-or-later** combined work. The `ooz`
decompressor is an independent reimplementation (© 2016 Powzix), not RAD/Epic's
proprietary Oodle SDK.
