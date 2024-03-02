$pythonScript = "sav_cli.py"
$distName = "sav_cli_windows_x86_64.exe"
$pyPipCommand = "pip install -r requirements.txt"
$pyInstallerCommand = "pyinstaller --onefile " + $pythonScript + " -n " + $distName
Invoke-Expression $pyPipCommand
Invoke-Expression $pyInstallerCommand
Write-Host "sav_cli.exe has been built successfully."
