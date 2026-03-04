# recon

A fast, single-binary CLI tool that scans your local filesystem for git repositories and serves them as a browsable web UI. Think of it as a local service catalog for navigating a large codebase.

## Features

- Recursively discovers all git repositories from a configurable root
- Browsable list view with repo summaries
- Detail view per repo: branch, last commit, remote URL, rendered README
- Zero runtime dependencies — single binary, pure Go
- Cross-platform: macOS, Linux, Windows

## Installation

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Danneeb/recon/master/install.sh | bash
```

Detects Intel or Apple Silicon automatically and installs to `/usr/local/bin`.

### Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Danneeb/recon/master/install.sh | bash
```

Installs to `/usr/local/bin` (uses `sudo` if needed). Supports x86-64 and arm64.

### Windows

The shell installer requires bash (Git Bash or WSL). If you have either, the curl one-liner above works as-is.

**Without bash** — download the zip from the [Releases](https://github.com/Danneeb/recon/releases) page, then in PowerShell:

```powershell
Expand-Archive recon_vX.Y.Z_windows_amd64.zip -DestinationPath C:\tools
# Add C:\tools to your PATH, or move recon.exe wherever you prefer
```

### go install (all platforms)

```bash
go install github.com/Danneeb/recon/cmd/recon@latest
```

### Manual download

| OS      | Architecture             | File                               |
| ------- | ------------------------ | ---------------------------------- |
| macOS   | Apple Silicon (M1/M2/M3) | `recon_vX.Y.Z_darwin_arm64.tar.gz` |
| macOS   | Intel                    | `recon_vX.Y.Z_darwin_amd64.tar.gz` |
| Linux   | x86-64                   | `recon_vX.Y.Z_linux_amd64.tar.gz`  |
| Linux   | ARM 64-bit               | `recon_vX.Y.Z_linux_arm64.tar.gz`  |
| Windows | x86-64                   | `recon_vX.Y.Z_windows_amd64.zip`   |
| Windows | ARM 64-bit               | `recon_vX.Y.Z_windows_arm64.zip`   |

Extract the archive and place the binary somewhere on your `PATH`.

## Usage

```bash
# Scan repos under the current directory and start the UI
recon serve

# Scan from a specific root
recon serve --root /code

# Just scan and print discovered repos (no server)
recon scan --root /code

# Change the port (default: 7474)
recon serve --port 8484

# Print version info
recon --version
```

Open [http://localhost:7474](http://localhost:7474) in your browser after starting the server.

## Building from source

```bash
git clone https://github.com/Danneeb/recon.git
cd recon
make build
./bin/recon --version
```

Run tests:

```bash
make test
```

## Release process

```bash
make release VERSION=1.2.3
git push origin master v1.2.3
```

This renames the `[Unreleased]` section in `CHANGELOG.md` to `[1.2.3] - YYYY-MM-DD`, commits, and creates the tag locally. Pushing the tag triggers the release workflow which builds all platform binaries and publishes the GitHub Release.

## License

MIT — see [LICENSE](LICENSE)
