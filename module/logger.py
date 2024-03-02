import logging
import sys
from datetime import datetime
from contextlib import contextmanager


class DateFileWrapper:

    def __init__(self, fileobj):
        self.fileobj = fileobj

    def write(self, data):
        timestamp = datetime.now().strftime("%Y/%m/%d - %H:%M:%S")
        if data.strip():
            data = f"{timestamp} | {data}"
        self.fileobj.write(data)

    def flush(self):
        self.fileobj.flush()

    def __getattr__(self, attr):
        return getattr(self.fileobj, attr)


@contextmanager
def redirect_stdout_stderr(
    stdout_path="save-tools-out.txt", stderr_path="save-tools-err.txt"
):
    original_stdout = sys.stdout
    original_stderr = sys.stderr

    try:
        with open(stdout_path, "a+", encoding="utf-8") as f_stdout, open(
            stderr_path, "a+", encoding="utf-8"
        ) as f_stderr:
            sys.stdout = DateFileWrapper(f_stdout)
            sys.stderr = DateFileWrapper(f_stderr)
            yield
    finally:
        sys.stdout = original_stdout
        sys.stderr = original_stderr


logging.basicConfig(
    level=logging.INFO,
    format="[SAV-CLI] %(asctime)s | %(levelname)s | %(message)s",
    datefmt="%Y/%m/%d - %H:%M:%S",
)


def log(text, level="INFO"):
    if level.upper() == "DEBUG":
        logging.debug(text)
    elif level.upper() == "INFO":
        logging.info(text)
    elif level.upper() == "WARNING":
        logging.warning(text)
    elif level.upper() == "ERROR":
        logging.error(text)
    elif level.upper() == "CRITICAL":
        logging.critical(text)
    else:
        logging.info(text)
