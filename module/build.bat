@echo off
start cmd /c pip install -r requirements && pyinstaller --onefile .\sav_cli.py -n sav_cli_windows_x86.exe