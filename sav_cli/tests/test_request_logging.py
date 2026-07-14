import importlib.util
import sys
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock


SAV_CLI_DIR = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SAV_CLI_DIR))
SPEC = importlib.util.spec_from_file_location("sav_cli_main", SAV_CLI_DIR / "sav_cli.py")
sav_cli = importlib.util.module_from_spec(SPEC)
with mock.patch.dict(
    sys.modules,
    {
        "structurer": SimpleNamespace(
            convert_sav=mock.Mock(),
            structure_player=mock.Mock(),
            structure_guild=mock.Mock(),
        ),
        "logger": SimpleNamespace(log=mock.Mock()),
    },
):
    SPEC.loader.exec_module(sav_cli)


class RequestLoggingTests(unittest.TestCase):
    def test_empty_error_response_still_logs_http_diagnostics(self):
        response = SimpleNamespace(
            status_code=502,
            reason="Bad Gateway",
            text="",
        )
        requests = SimpleNamespace(put=mock.Mock(return_value=response))

        with tempfile.NamedTemporaryFile() as save_file:
            with (
                mock.patch.object(
                    sys,
                    "argv",
                    [
                        "sav_cli.py",
                        "--file",
                        save_file.name,
                        "--request",
                        "http://127.0.0.1:8080/api/",
                    ],
                ),
                mock.patch.dict(sys.modules, {"requests": requests}),
                mock.patch.object(sav_cli, "convert_sav"),
                mock.patch.object(sav_cli, "structure_player", return_value=[]),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                sav_cli.main()

        errors = [
            call.args[0]
            for call in log.call_args_list
            if len(call.args) > 1 and call.args[1] == "ERROR"
        ]
        self.assertEqual(len(errors), 2)
        for message in errors:
            self.assertIn("HTTP 502 Bad Gateway", message)
            self.assertIn("empty response body", message)


if __name__ == "__main__":
    unittest.main()
