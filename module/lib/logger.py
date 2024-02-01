from datetime import datetime


def log(text, level="INFO"):
    prefix = "[SAV-CLI]"
    current = datetime.now().strftime("%Y/%m/%d - %H:%M:%S")
    print(f"{prefix} {current} | {level.upper()} | {text}")
