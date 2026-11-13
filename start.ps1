param(
    [string]$Command = "run",
    [string]$Arg = ""
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BuildDir = Join-Path $RootDir "build"
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

function Build-Host {
    Write-Host "Building MicroDeck host..."
    Push-Location (Join-Path $RootDir "host")
    go build -ldflags="-s -w" -o (Join-Path $BuildDir "microdeck.exe") main.go
    Pop-Location
}

function Build-App($AppName) {
    $AppSrc = Join-Path $RootDir "apps\$AppName"
    if (-not (Test-Path $AppSrc)) {
        Write-Error "App '$AppName' not found in apps/"
    }
    Write-Host "Building micro-app: $AppName..."
    Push-Location $AppSrc
    go build -ldflags="-s -w" -o (Join-Path $AppSrc "$AppName.exe") main.go
    Pop-Location
}

function Scaffold-App($NewName) {
    $TargetDir = Join-Path $RootDir "apps\$NewName"
    if (Test-Path $TargetDir) {
        Write-Error "App '$NewName' already exists."
    }
    New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null

    @"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://apple.com">
<plist version="1.0">
<dict>
    <key>CFBundleIdentifier</key>
    <string>com.microdeck.$NewName</string>
    <key>CFBundleDisplayName</key>
    <string>$NewName</string>
    <key>CFBundleExecutable</key>
    <string>$NewName</string>
    <key>MicroDeckDescription</key>
    <string>Custom MicroDeck utility</string>
</dict>
</plist>
"@ | Out-File -Encoding UTF8 (Join-Path $TargetDir "Info.plist")

    @"
package main

import "fmt"

func main() {
    fmt.Println("Hello from $NewName!")
}
"@ | Out-File -Encoding UTF8 (Join-Path $TargetDir "main.go")

    Write-Host "[OK] Created new micro-app at apps/$NewName"
}

switch ($Command) {
    "build" {
        Build-Host
        Get-ChildItem (Join-Path $RootDir "apps") -Directory | ForEach-Object {
            Build-App $_.Name
        }
    }
    "new" {
        if (-not $Arg) { Write-Error "Usage: .\start.ps1 new <app_name>" }
        Scaffold-App $Arg
    }
    "install" {
        if (-not (Test-Path (Join-Path $BuildDir "microdeck.exe"))) { Build-Host }
        if (-not $Arg) { Write-Error "Usage: .\start.ps1 install <app_name>" }
        Build-App $Arg
        & (Join-Path $BuildDir "microdeck.exe") install (Join-Path $RootDir "apps\$Arg")
    }
    "run" {
        if (-not (Test-Path (Join-Path $BuildDir "microdeck.exe"))) { Build-Host }
        & (Join-Path $BuildDir "microdeck.exe")
    }
    default {
        Write-Host "Usage: .\start.ps1 {build|new <name>|install <name>|run}"
    }
}
