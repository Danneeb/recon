# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
