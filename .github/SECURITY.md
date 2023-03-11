# Security Policy

## Supported  Versions and Statuses

| Version/Section | Status | Note |
| :------ | :----- | :--- |
| Go 1.18, 1.19 ... latest | [![go1.18+](https://github.com/KEINOS/go-gisty/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/KEINOS/go-gisty/actions/workflows/unit-tests.yml "Unit tests on various Go versions") | Including Go 1.19 |
| Golangci-lint 1.51.2 or later | [![golangci-lint](https://github.com/KEINOS/go-gisty/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/go-gisty/actions/workflows/golangci-lint.yml) | |
| Security advisories | [Enabled](https://github.com/KEINOS/go-gisty/security/advisories) | |
| Dependabot alerts | [Enabled](https://github.com/KEINOS/go-gisty/security/dependabot) | (Viewable only for admins) |
| Code scanning alerts | [Enabled](https://github.com/KEINOS/go-gisty/security/code-scanning)<br>[![CodeQL-Analysis](https://github.com/KEINOS/go-gisty/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/go-gisty/actions/workflows/codeQL-analysis.yml) ||

## Update

- We [check the latest version of `go.mod` every week](https://github.com/KEINOS/go-gisty/blob/main/.github/workflows/weekly-update.yml) and update it when it has passed all tests.

## Reporting a Vulnerability, Bugs and etc

- [Issues](https://github.com/KEINOS/go-gisty/issues)
  - [![Opened Issues](https://img.shields.io/github/issues/KEINOS/go-gisty?color=lightblue&logo=github)](https://github.com/KEINOS/go-gisty/issues "opened issues")
  - Plase attach a simple test that replicates the issue.
