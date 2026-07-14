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
        "logger": SimpleNamespace(configure_logging=mock.Mock(), log=mock.Mock()),
    },
):
    SPEC.loader.exec_module(sav_cli)


class RequestLoggingTests(unittest.TestCase):
    def test_partial_player_save_parse_is_reported_as_a_warning(self):
        response = SimpleNamespace(status_code=200, reason="OK", text="")
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
                mock.patch.object(
                    sav_cli,
                    "structure_player",
                    return_value=([], 1),
                ),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                exit_code = sav_cli.main()

        self.assertEqual(exit_code, 0)
        messages = "\n".join(call.args[0] for call in log.call_args_list)
        self.assertRegex(
            messages,
            r"Structured save: .*player_save_warnings=1 \(\d+\.\d{2}s\)",
        )
        self.assertRegex(
            messages,
            r"Save sync completed with warnings=1 \(\d+\.\d{2}s\)",
        )

    def test_json_error_response_logs_only_the_server_reason(self):
        response = SimpleNamespace(
            status_code=401,
            reason="Unauthorized",
            text='{"error":"token expired"}',
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
                mock.patch.object(sav_cli, "structure_player", return_value=([], 0)),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                sav_cli.main()

        error = next(
            call.args[0]
            for call in log.call_args_list
            if call.args[0].startswith("Failed to sync players=")
        )
        self.assertIn("HTTP 401 Unauthorized; error: token expired", error)
        self.assertNotIn('{"error"', error)

    def test_verbose_flag_enables_dependency_process_logs(self):
        with tempfile.NamedTemporaryFile() as save_file:
            with (
                mock.patch.object(
                    sys,
                    "argv",
                    ["sav_cli.py", "--file", save_file.name, "--verbose"],
                ),
                mock.patch.object(sav_cli, "convert_sav"),
                mock.patch.object(sav_cli, "structure_player", return_value=([], 0)),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "configure_logging") as configure_logging,
                mock.patch.object(sav_cli, "log"),
            ):
                exit_code = sav_cli.main()

        self.assertEqual(exit_code, 0)
        configure_logging.assert_called_once_with(verbose=True)

    def test_http_error_body_is_single_line_and_truncated(self):
        response = SimpleNamespace(
            status_code=413,
            reason="Content Too Large",
            text="request body rejected\n" + ("x" * 3000),
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
                mock.patch.object(sav_cli, "structure_player", return_value=([], 0)),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                sav_cli.main()

        error = next(
            call.args[0]
            for call in log.call_args_list
            if call.args[0].startswith("Failed to sync players=")
        )
        self.assertNotIn("\n", error)
        self.assertIn("response body: request body rejected ", error)
        self.assertIn("<truncated>", error)
        self.assertLess(len(error), 800)

    def test_request_exception_is_logged_and_other_resource_still_syncs(self):
        class RequestError(Exception):
            pass

        success = SimpleNamespace(status_code=200, reason="OK", text="")
        requests = SimpleNamespace(
            RequestException=RequestError,
            put=mock.Mock(side_effect=[RequestError("connection refused"), success]),
        )

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
                mock.patch.object(sav_cli, "structure_player", return_value=([], 0)),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                exit_code = sav_cli.main()

        self.assertEqual(exit_code, 1)
        messages = "\n".join(call.args[0] for call in log.call_args_list)
        self.assertRegex(
            messages,
            r"Failed to sync players=0: RequestError: connection refused \(\d+\.\d{2}s\)",
        )
        self.assertRegex(messages, r"Synced guilds=0: HTTP 200 OK \(\d+\.\d{2}s\)")
        self.assertRegex(
            messages,
            r"Save sync failed: requests_failed=1 \(\d+\.\d{2}s\)",
        )

    def test_successful_sync_logs_concise_stage_summaries(self):
        response = SimpleNamespace(status_code=200, reason="OK", text='{"success":true}')
        requests = SimpleNamespace(put=mock.Mock(return_value=response))
        players = [
            {
                "player_uid": "player-1",
                "pals": [{"name": "Lamball"}, {"name": "Cattiva"}],
            }
        ]
        guilds = [
            {
                "players": [{"player_uid": "player-1", "last_online": 123}],
                "base_camp": [{"id": "camp-1"}],
            }
        ]

        with tempfile.NamedTemporaryFile(suffix="Level.sav") as save_file:
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
                mock.patch.object(
                    sav_cli, "structure_player", return_value=(players, 0)
                ),
                mock.patch.object(sav_cli, "structure_guild", return_value=guilds),
                mock.patch.object(sav_cli, "log") as log,
            ):
                exit_code = sav_cli.main()

        self.assertEqual(exit_code, 0)
        messages = "\n".join(call.args[0] for call in log.call_args_list)
        self.assertRegex(messages, r"Decoding save: .*Level\.sav")
        self.assertRegex(messages, r"Decoded save \(\d+\.\d{2}s\)")
        self.assertIn("Structuring save data", messages)
        self.assertRegex(
            messages,
            r"Structured save: players=1, pals=2, guilds=1, base_camps=1 \(\d+\.\d{2}s\)",
        )
        self.assertIn("Syncing save data: http://127.0.0.1:8080/api/", messages)
        self.assertRegex(messages, r"Synced players=1: HTTP 200 OK \(\d+\.\d{2}s\)")
        self.assertRegex(messages, r"Synced guilds=1: HTTP 200 OK \(\d+\.\d{2}s\)")
        self.assertRegex(messages, r"Save sync completed \(\d+\.\d{2}s\)")
        self.assertNotIn("Put players", messages)
        self.assertNotIn("Done in", messages)

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
                mock.patch.object(sav_cli, "structure_player", return_value=([], 0)),
                mock.patch.object(sav_cli, "structure_guild", return_value=[]),
                mock.patch.object(sav_cli, "log") as log,
            ):
                exit_code = sav_cli.main()

        self.assertEqual(exit_code, 1)
        errors = [
            call.args[0]
            for call in log.call_args_list
            if len(call.args) > 1 and call.args[1] == "ERROR"
        ]
        self.assertEqual(len(errors), 3)
        for message in errors[:2]:
            self.assertIn("HTTP 502 Bad Gateway", message)
            self.assertIn("empty response body", message)
        self.assertRegex(
            errors[-1],
            r"Save sync failed: requests_failed=2 \(\d+\.\d{2}s\)",
        )


if __name__ == "__main__":
    unittest.main()
