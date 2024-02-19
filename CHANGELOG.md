# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- `check.Equal` and `assert.Equal` had the expected and actual values flipped in the test failure error message.

## [0.3.1] - 2024-02-11

### Fixed

- `check.Nil` and `assert.Nil` now handle nil pointers.

## [0.3.0] - 2024-02-01

### Added

- Added `check` module. This module contains the same functions as the
  `assert` module, but they return a boolean and mark the test as failed
  instead of failing the test immediately.

- `NotErrorIs`, `NotErrorContains`, and `NotDeepEqual` functions.

### Removed

- `LenEqual`, `Assert`, and `Check` functions.

## [0.2.0] - 2024-01-28

### Changed

- Moved assertion functions to a `assert` module and removed `Assert` suffix
  from many of the functions.

## [0.1.0] - 2024-01-28

- Initial Release
