# SiteRippers
![loc](https://sloc.xyz/github/nektro/SiteRippers)
[![license](https://img.shields.io/github/license/nektro/SiteRippers.svg)](https://github.com/nektro/SiteRippers/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)
[![paypal](https://img.shields.io/badge/donate-paypal-009cdf)](https://paypal.me/nektro)
[![circleci](https://circleci.com/gh/nektro/SiteRippers.svg?style=svg)](https://circleci.com/gh/nektro/SiteRippers)
[![release](https://img.shields.io/github/v/release/nektro/SiteRippers)](https://github.com/nektro/SiteRippers/releases/latest)
[![goreportcard](https://goreportcard.com/badge/github.com/nektro/SiteRippers)](https://goreportcard.com/report/github.com/nektro/SiteRippers)
[![codefactor](https://www.codefactor.io/repository/github/nektro/SiteRippers/badge)](https://www.codefactor.io/repository/github/nektro/SiteRippers)

A collection of Golang scripts to do entire rips of sites centralized in a single repo.

## Prerequisites
- Golang 1.12+

## Installing
```sh
$ go get -v -u github.com/nektro/SiteRippers
```

## Usage
```
Usage of ./SiteRippers:
      --concurrency int    Maximum number of tasks to run at once. Exactly how tasks are used varies slightly. (default 10)
      --save-dir string    Path to folder to save downloaded data to. (default "./data/")
      --site stringArray   List of domains of sites to rip. None passed means rip all.
```

## License
MIT
