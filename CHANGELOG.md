# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2026-03-09

### Added
- Branch dropdown on the detail view — select any local branch to view its commit history and stats

## [0.4.0] - 2026-03-04

### Added
- Live search/filter input on the By contributor chart — filters rows instantly as you type

### Changed
- By contributor list now shows all contributors instead of capping at 10
- By contributor list is now vertically scrollable instead of expanding the page
- Day of week chart now starts on Monday and ends on Sunday

## [0.3.0] - 2026-03-04

### Added
- CLI subcommands: `recon scan` to print discovered repos, `recon version` to print build info
- `--root` flag to set the scan root directory (defaults to current directory)
- `--port` flag to set the server port (defaults to 8484)
- Running `recon` with no arguments starts the server with defaults — no subcommand needed
- Startup banner showing version, commit, build date, root, and port
- `version`, `commit`, and `date` build variables injected via `-ldflags` (populated by `make build` and the release workflow)

### Changed
- README updated to reflect actual implemented features and correct CLI usage

## [0.2.3] - 2026-03-04

### Fixed
- Change default server port from 8080 to 8484 to avoid conflicts with other common local services

## [0.2.2] - 2026-03-03

### Fixed
- Commit log and activity chart now sort correctly on repos with merge commits by using committer-time ordering

## [0.2.1] - 2026-03-03

### Fixed
- Commit activity bar chart now scrolls horizontally for repos with long histories instead of compressing bars

## [0.2.0] - 2026-03-03

### Added
- Detail view with stats dashboard: commit activity bar chart (full history), commits by contributor, commits by day of week
- Tabbed UI on detail page — Overview (stats) and Commit Log tabs
- `GetRepoDetail` function using go-git to extract full commit history, contributors, and activity breakdowns
- `BarEntry` type for pre-computed chart data (label, count, percentage)
- Tests for `GetRepoDetail` covering commit count, contributors, date formatting, author ranking, monthly range, and day-of-week bucketing

## [0.1.2] - 2026-03-02

### Fixed
- Skip repos where HEAD cannot be resolved (e.g. empty repos) instead of aborting the scan
- Repo names no longer show as `.` when scanning from the current directory

## [0.1.1] - 2026-03-01

### Fixed
- Scan now starts from the current working directory (`./`) instead of the parent directory (`../`)

## [0.1.0] - 2026-03-01

### Added
- Recursive git repo discovery via `filepath.WalkDir` with configurable ignore list
- go-git integration for extracting branch metadata per repo
- Web server with repo list view rendered via `html/template`
- Tailwind CSS grid layout with dark/light theme toggle
- Navbar with placeholder routes for future pages (Overview, Settings)
- Clickable terminal hyperlink printed on startup
- Embedded templates via `//go:embed`
