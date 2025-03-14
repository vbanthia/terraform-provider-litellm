# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.4] - 2025-03-13

### Added
- Added new `thinking` capability to model resource with configurable parameters:
  - `thinking_enabled` - Boolean to enable/disable thinking capability (default: false)
  - `thinking_budget_tokens` - Integer to set token budget for thinking (default: 1024)

### Fixed
- Fixed issue where `thinking_budget_tokens` was being added to models that don't have `thinking_enabled = true`

## [0.2.2] - 2025-02-06

### Added
- Added new `reasoning_effort` parameter to model resource with values: "low", "medium", "high"
- Added "chat" mode to model resource

### Changed
- Updated model mode options to: "completion", "embedding", "image_generation", "chat", "moderation", "audio_transcription"

## [1.0.0] - 2024-01-17

### Added
- Initial release of the LiteLLM Terraform Provider
- Support for managing LiteLLM models
- Support for managing teams and team members
- Comprehensive documentation for all resources
