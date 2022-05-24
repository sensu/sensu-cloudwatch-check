# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.3.0] - 2022-05-24
### Changed
- Update to sensu-plugin-sdk v0.16
### Fixed
- Fixed Prometheus output formatting

## [0.2.0] - 2022-04-20

- Fixed preset metrics handling for situation where requests metrics does not include All service metrics
- Add new option to enable error return status if service metrics are missing from requested metric list.  Useful for development of new presets or CI testing to catch when AWS extends a metrics set. 

## [0.1.0] - 2022-04-19

### Added
- Initial release
