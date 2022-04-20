# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

- Fixed preset metrics handling for situation where requests metrics does not include All service metrics
- Add new option to enable error return status if service metrics are missing from requested metric list.  Useful for development of new presets or CI testing to catch when AWS extends a metrics set. 

## [0.0.1] - 2000-01-01

### Added
- Initial release
