# Contributing to alert-tune

Thanks for helping improve alert-tune. Bug reports, documentation fixes, tests, and code changes are welcome.

The project is in early development. The CLI currently supports `version` and `--version`; the README describes planned backtesting and tuning features.

## Before you start

- Search existing issues and pull requests for related work.
- Open an issue to discuss large changes or new features before implementing them.
- Report security vulnerabilities privately using [SECURITY.md](SECURITY.md).

## Set up your checkout

Install Git, Make, and the Go version specified in [go.mod](go.mod).

Fork the repository on GitHub, then replace `YOUR-USERNAME` below with your GitHub username:

```sh
git clone https://github.com/YOUR-USERNAME/alert-tune.git
cd alert-tune
git remote add upstream https://github.com/kang0225/alert-tune.git
git fetch upstream
git switch -c fix/describe-your-change upstream/main
go mod download
```

Use a descriptive branch name, such as `fix/version-output` or `docs/contributing`.

## Make and validate changes

Keep each contribution focused on one problem. Follow the surrounding code style, format Go files with `gofmt`, and add meaningful tests when introducing or changing behavior. Update documentation when behavior or usage changes.

Run the shared checks before submitting:

```sh
make check
```

This checks Go formatting with `gofmt -l .`, runs `go vet ./...` and `go test ./...`, builds `bin/alert-tune`, and smoke-tests that `version` and `--version` print the same nonempty version output. There are currently no Go test files; add tests alongside new behavior.

To run the current CLI manually:

```sh
make build
./bin/alert-tune version
./bin/alert-tune --version
```

## Submit a pull request

Write commit messages in English using Conventional Commit style, such as `fix(cli): handle invalid arguments` or `docs: clarify setup`. Add a `Signed-off-by` trailer with `git commit -s`, following the repository's existing commit style:

Replace `PATH_TO_CHANGED_FILES` below with the paths of the files you changed.

```sh
git add PATH_TO_CHANGED_FILES
git commit -s -m "docs: clarify setup"
git push -u origin HEAD
```

Open a pull request against `kang0225/alert-tune`'s `main` branch. Fill in the pull request template with the problem, resulting behavior, actual validation commands and results, and any risks or limitations. Link related issues and respond to review feedback.

If a check cannot run, explain what prevented it and which validation you completed. Only mark a validation checkbox after performing that check.
