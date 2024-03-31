# Changelog

## [1.4.0](https://github.com/yemaney/z/compare/v1.3.1...v1.4.0) (2024-03-31)


### Features

* **aws:** add aws feat to list, get, start, stop ec2 instances ([e8bd9d0](https://github.com/yemaney/z/commit/e8bd9d038ff39405793b844378d4c0972fade86a))
* **aws:** add cmds branches to create and delete ec2 instanes ([b1033bb](https://github.com/yemaney/z/commit/b1033bb03cc5871522989e58aa7022c5de1875ee))
* **aws:** add conf support and arg struct for create cmd ([fccf7af](https://github.com/yemaney/z/commit/fccf7af538091ca77e5a66fd733b54f8c0cc415d))
* **aws:** add options for create cmd ([3e4a4bc](https://github.com/yemaney/z/commit/3e4a4bc09664cb3a5690e85ec95d02333f759c5d))


### Bug Fixes

* add missing comma in commands slice ([285b0e5](https://github.com/yemaney/z/commit/285b0e56df79e91bb44aecb04322e63c8021823f))
* **aws:** handle cases when keyname is nil and add expected path to key ([45f7084](https://github.com/yemaney/z/commit/45f708482ce6ede4fae7e310d7d09c10764d8f11))

## [1.3.1](https://github.com/yemaney/z/compare/v1.3.0...v1.3.1) (2024-03-25)


### Bug Fixes

* **ssh:** enable patch host and create ~/.ssh dir if it doesn't exist ([ff2f3ea](https://github.com/yemaney/z/commit/ff2f3ea2737559283b6ce24fe50db257b8a44028))

## [1.3.0](https://github.com/yemaney/z/compare/v1.2.0...v1.3.0) (2024-02-24)


### Features

* **ssh:** add option to retrieve all sections in get command ([e31247b](https://github.com/yemaney/z/commit/e31247b08e13945c223604c350ba6423f9873db2))
* **ssh:** add patch command ([e10a78b](https://github.com/yemaney/z/commit/e10a78b873ae721c67bfaf45bb73dd4372998fb1))


### Bug Fixes

* **ssh:** remove unneeded params field in exported cmd ([b82b32a](https://github.com/yemaney/z/commit/b82b32a558886a4367b08aa172c00f48a5795cb0))

## [1.2.0](https://github.com/yemaney/z/compare/v1.1.0...v1.2.0) (2024-02-23)


### Features

* **ssh:** add get command ([6917b48](https://github.com/yemaney/z/commit/6917b481fc1e80c3ee8cd5ef87b9bc97b52484dd))


### Bug Fixes

* **cc:** remove uneeded commit types ([d768a22](https://github.com/yemaney/z/commit/d768a22c07b2c96f8423dd98b5911d3fd03115cd))

## [1.1.0](https://github.com/yemaney/z/compare/v1.0.0...v1.1.0) (2024-02-22)


### Features

* **ssh:** add delete section command ([e2ebcc3](https://github.com/yemaney/z/commit/e2ebcc394c9f90292350fd030acc286839ab09d9))

## [1.0.0](https://github.com/yemaney/z/compare/v0.2.0...v1.0.0) (2024-02-20)


### ⚠ BREAKING CHANGES

* upgrade to go 1.22

### Build System

* upgrade to go 1.22 ([3f97f9d](https://github.com/yemaney/z/commit/3f97f9d0da207f3742dc3824c829b226a9f24ce1))

## [0.2.0](https://github.com/yemaney/z/compare/v0.1.1...v0.2.0) (2024-02-20)


### Features

* **cc:** add option for breaking changes commit messages ([cd15fcf](https://github.com/yemaney/z/commit/cd15fcfbf8b73419bb726e64a1140255127e88ce))

## [0.1.1](https://github.com/yemaney/z/compare/v0.1.0...v0.1.1) (2024-02-19)


### Bug Fixes

* **completion:** add custom command completer ([f53728a](https://github.com/yemaney/z/commit/f53728a0e7c17c72997d6635051a26ff0e8cb7a4))

## [0.1.0](https://github.com/yemaney/z/compare/v0.0.1...v0.1.0) (2024-02-19)


### Features

* **ssh:** add new ssh command ([2062761](https://github.com/yemaney/z/commit/2062761c7323ca0f15899dba25d8bf820c6d6b2e))
