import contextlib
import io
import logging
import sys
import unittest
from pathlib import Path


sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

import logger


class LoggerTests(unittest.TestCase):
    def test_configure_logging_uses_one_format_and_hides_dependency_info(self):
        output = io.StringIO()

        with contextlib.redirect_stderr(output):
            logger.configure_logging()
            logger.log("Decoding save")
            logging.getLogger("palsav.compressor.oozlib").info(
                "Starting decompression process with palooz"
            )

        lines = output.getvalue().splitlines()
        self.assertEqual(len(lines), 1)
        self.assertRegex(
            lines[0],
            r"^\[SAV-CLI\] \d{4}/\d{2}/\d{2} - \d{2}:\d{2}:\d{2} \| INFO \| Decoding save$",
        )

    def test_verbose_dependency_logs_keep_the_same_format(self):
        output = io.StringIO()

        with contextlib.redirect_stderr(output):
            logger.configure_logging(verbose=True)
            logging.getLogger("palsav.compressor.oozlib").info(
                "Starting decompression process with palooz"
            )

        self.assertRegex(
            output.getvalue().strip(),
            r"^\[SAV-CLI\] \d{4}/\d{2}/\d{2} - \d{2}:\d{2}:\d{2} \| INFO \| Starting decompression process with palooz$",
        )


if __name__ == "__main__":
    unittest.main()
