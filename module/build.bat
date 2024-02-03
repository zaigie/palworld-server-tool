@echo off
start cmd /c pip install -r requirements.txt && pyinstaller --onefile .\sav_cli.py -n sav_cli_windows_x86_64.exe