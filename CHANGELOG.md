# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] -- 2020-07-18
### Breaking
- [protocol]
  - `LookupKey` is now a struct, not a string
  - `LookupKey`s now support the `json.Marshaler` and `json.Unmarshsler` interfaces
  - `LookupKey`s are now formatted with their type and a : before the shasum
    - Types can be short
    - `"h:#{partial-key}"`, `"host:#{partial-key}"`
  - Signed Certificates now exist under their own subprefix in S3.
### Changed
- [protocol]
  - `CaKeyPair` (`cakp`) added to `CertType` for referencing specific Ca KeyPairs
- [deps] - Update aws-sdk-go:v1.33.x
### Added
- [protocol]
  - `S3Object.LoadObject`
    - populate fields of various s3 object structs with data from s3

## [0.5.1] - 2020-06-15
### Changed
- [deps] - Update aws-sdk-go:v1.32.2

## [0.5.0] - 2020-06-12
### Added
- [protocol] - Ability to expand partial `LookupKey`s
### Changed
- [protocol] - `CAPublicKeyS3Object` now supports `KeyFingerprint`

[protocol]: <https://src.doom.fm/schism/commonLib/-/tree/v0.6.0-dev/protocol>
[deps]:<https://src.doom.fm/schism/commonLib/-/blob/v0.6.0-dev/go.mod>

[0.5.0]:<https://src.doom.fm/schism/commonLib/-/tags/v0.5.0>
[0.5.1]:<https://src.doom.fm/schism/commonLib/-/tags/v0.5.1>
[0.6.0]:<https://src.doom.fm/schism/commonLib/-/tags/v0.6.0>
[Unreleased]:<https://src.doom.fm/schism/commonLib/-/compare/v0.5.1...v0.6.0-dev>
