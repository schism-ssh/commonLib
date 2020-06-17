# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- BREAKING [protocol] -
  - `LookupKey` is now a struct, not a string
  - Keys now (Un)Marshal in/out of JSON
  - Keys now are formatted with their type and a : before the shasum

## [0.5.1] - 2020-06-15
### Changed
- [deps] - Update aws-sdk-go@1.32.2
## [0.5.0] - 2020-06-12
### Added
- [protocol] - Ability to expand partial `LookupKey`s
### Changed
- [protocol] - `CAPublicKeyS3Object` now supports `KeyFingerprint`

