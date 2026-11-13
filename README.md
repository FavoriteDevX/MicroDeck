# MicroDeck

> *"Today, we are introducing three revolutionary concepts to the whole GitHub community.*

> *The first is a **zero-bloat native app launcher**.*
> *The second is a **cross-platform CLI package installer**.*
> *And the third is a **decentralized, zero-App-Store ecosystem**.*

<img width="1368" height="768" alt="image" src="https://github.com/user-attachments/assets/6c4e50de-2738-47b2-a047-f871afd021f9" />

>
> *A native launcher. A CLI installer. A zero-store ecosystem.*
> *Are you getting it?
>
>These are not three separate tools.
>
>This is one unified runtime.*

> *And we are calling it **MicroDeck**."*

<img width="1195" height="896" alt="image" src="https://github.com/user-attachments/assets/c0883fd6-99d6-48a0-bb91-03941e86d517" />

---

**MicroDeck** lets you run, build, and install lightweight, self-contained micro-apps on **Apple macOS** and **Windows** — no **App Store**, no **Microsoft Store**, and no 300 MB Electron wrappers required.

Inspired by classic embedded homebrew launchers, MicroDeck brings fast, single-binary utilities to modern desktop operating systems. Every micro-app lives inside a clean, standardized directory structure (`/Apps`, `/Data`, `/Library`) with an Apple-style `Info.plist` manifest and launches in milliseconds.

## Why MicroDeck?

- **Zero-Store Installation (`install`)**: Distribute and install micro-apps directly from your terminal or GitHub without signing certificates, `.dmg` wizards, or App Store review queues.
- **Native Speed on macOS & Windows**: Written in Go to compile into tiny, standalone native executables for Apple Silicon, Intel Macs, and Windows 10/11.
- **Standardized `.microapp` Bundles**: Just like classic Mac bundles, every tool is a self-contained `.microapp` folder pairing an executable with an `Info.plist` identity file.
- **One-Command Workflow**: Scaffold, compile, install, and launch tools using a single entry point (`./start` on macOS/Linux or `.\start.ps1` on Windows).

## Included Micro-Apps

Out of the box, MicroDeck ships with ready-to-use native utilities:

- **System Info (`sysinfo`)**: Instant hardware, OS architecture, CPU core, and memory diagnostics.
- **QuickPass (`quickpass`)**: Cryptographically secure, zero-dependency password generator.

## Quick Start

### Apple macOS / Linux
```bash
git clone https://github.com
cd MicroDeck
chmod +x start

# Build the host and sample micro-apps
./start build

# Install a micro-app into ~/MicroDeck/Apps
./start install sysinfo
./start install quickpass

# Launch the interactive MicroDeck menu
./start run
```

### Windows (PowerShell)
```powershell
git clone https://github.com
cd MicroDeck

# Build the host and sample micro-apps
.\start.ps1 build

# Install a micro-app into %USERPROFILE%\MicroDeck\Apps
.\start.ps1 install sysinfo
.\start.ps1 install quickpass

# Launch the interactive MicroDeck menu
.\start.ps1 run
```

## Build & Install Your Own Micro-App

Creating a new utility takes seconds. Use the `new` command to scaffold a fresh `.microapp` project with its own `Info.plist` and Go entry point:

```bash
./start new my_tool
```

Edit `apps/my_tool/main.go`, then compile and install it straight to your MicroDeck launcher:

```bash
./start install my_tool
```

### Anatomy of a `.microapp` Bundle

Every installed app resides under `~/MicroDeck/Apps/<AppName>.microapp/` (or `%USERPROFILE%\MicroDeck\Apps\` on Windows) with two core files:

```text
System Info.microapp/
├── Info.plist    # Bundle ID, Display Name, Executable name, and Description
└── sysinfo       # Self-contained native binary (sysinfo.exe on Windows)
```

Persistent app files are automatically isolated under `~/MicroDeck/Data/<CFBundleIdentifier>/` via the `MICRODECK_DATA_DIR` environment variable and the included `sdk` package.

## Contributing

Go and visit GitHub Multiverse repo lol, why not?
