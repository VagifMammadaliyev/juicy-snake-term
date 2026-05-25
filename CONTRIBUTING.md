# Contributing to Juicy Snake

Thanks for taking the time to contribute! This document outlines how to propose changes, report issues, and get your patch merged.

By participating in this project you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Ways to contribute

- Report bugs by opening an issue using the **Bug report** template.
- Suggest features or improvements by opening an issue.
- Improve documentation, examples, or comments.
- Fix open issues and submit a pull request.

If you are unsure whether a change is welcome, open an issue first to discuss it before investing time in a pull request.

## Development setup

1. Install [Go 1.26](https://go.dev/dl/) or later.
2. Fork and clone the repository:
   ```sh
   git clone git@github.com:<your-username>/juicy-snake-term.git
   cd juicy-snake-term
   ```
3. Build the server and client:
   ```sh
   go build ./cmd/server
   go build ./cmd/game
   ```
4. Run the test suite in race check mode:
   ```sh
   go test -race ./...
   ```

## Making changes

1. Create a topic branch off `master`:
   ```sh
   git checkout -b feat/short-description
   ```
   Use a prefix that reflects intent: `feat/`, `fix/`, `chore/`, `docs/`, `refactor/`, `test/`.
2. Keep commits focused and write descriptive commit messages.
3. Run `go fmt ./...` and `go vet ./...` before pushing.
4. Add or update tests when you change behavior. Logic-heavy code lives under `internal/logic/` and should be covered by unit tests.
5. Update the README or other docs if your change affects how the game is built, configured, or played.

## Use of AI tools

We want the game logic to be understood and owned by the people who write it, so the use of large language models (LLMs) is intentionally limited:

- **Features and bug fixes: prohibited.** Do not use an LLM to generate the code for features or bug fixes. This is the core of the project and should be written and reasoned about by a human.
- **Tests: allowed, with care.** LLMs may help write tests, but never blindly. Decide which cases you want to cover first, then write a focused prompt describing those cases. Review every generated test carefully — a test you do not understand is worse than no test at all.
- **Documentation: allowed for prose outside the code.** LLMs are fine for READMEs, docs files, and similar standalone documentation. They are **not** for code comments or doc comments, which should be written by hand alongside the code they describe.

If you use an LLM where it is allowed, you are still fully responsible for the result and must review it as if you had written it yourself.

## Submitting a pull request

1. Push your branch and open a pull request against `master`.
2. Fill out the pull request template and link any related issues.
3. Make sure CI passes (build, vet, tests).
4. Implement the changes agreed upon in the pull request discussion.

## Reporting security issues

Please **do not** open public issues for security vulnerabilities. Instead, contact the maintainer directly at <vagifmammadaliyev@outlook.com> with details and reproduction steps.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
