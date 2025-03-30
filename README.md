# Rpkgengine
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rsdate/rpkgengine/go.yml)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/rsdate/rpkgengine/latest)
![Codecov](https://img.shields.io/codecov/c/github/rsdate/rpkgengine)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=rsdate_rpkgengine&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=rsdate_rpkgengine)

Rpkgengine is the internal engine that [`rpkg`](github.com/rsdate/rpkg) uses. Right now, it contains the code that the `install` command will use. The API is not meant for outside use, having parameters specific to `rpkg`.
