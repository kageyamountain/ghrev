# ghrev

[日本語](README.md) | English

`ghrev` is a CLI tool for measuring Pull Request review metrics on GitHub repositories.
By visualizing code reviews from various angles, it helps improve review culture and surface how much time is spent on code review.

## Subcommands
### `approval`
For PRs created within the specified period, aggregates how long it took from PR open until the required number of approvals were collected.
```sh
ghrev approval \
  --owner <organization-or-user> \
  --name <repository> \
  --from <YYYYMMDD> \
  --to <YYYYMMDD> \
  --required-approvals <N> \
  --ignore-labels <label1,label2,...> \
  --assignees <user1,user2,...>
```

| Option | Required | Description |
| --- | :---: | --- |
| `--owner` | ✓ | Owner of the target repository (user or organization name) |
| `--name` | ✓ | Target repository name |
| `--from` | ✓ | Start date of aggregation (`YYYYMMDD` format) |
| `--to` | ✓ | End date of aggregation (`YYYYMMDD` format) |
| `--required-approvals` | ✓ | Number of approvals required to consider a review complete (integer ≥ 1) |
| `--ignore-labels` |   | Labels to exclude from aggregation (comma-separated) |
| `--assignees` |   | Include only PRs that have any of the specified assignees (comma-separated) |

### `first-review`
For PRs created within the specified period, aggregates how long it took from PR open until the first review reaction (any of approve / changes_requested / commented). DISMISSED reviews, bot reviews, and self-reviews by the PR author are not counted as reactions.
```sh
ghrev first-review \
  --owner <organization-or-user> \
  --name <repository> \
  --from <YYYYMMDD> \
  --to <YYYYMMDD> \
  --ignore-labels <label1,label2,...> \
  --assignees <user1,user2,...>
```

| Option | Required | Description |
| --- | :---: | --- |
| `--owner` | ✓ | Owner of the target repository (user or organization name) |
| `--name` | ✓ | Target repository name |
| `--from` | ✓ | Start date of aggregation (`YYYYMMDD` format) |
| `--to` | ✓ | End date of aggregation (`YYYYMMDD` format) |
| `--ignore-labels` |   | Labels to exclude from aggregation (comma-separated) |
| `--assignees` |   | Include only PRs that have any of the specified assignees (comma-separated) |

### `help`
Shows help for the available subcommands.
```sh
ghrev help
```

### `version`
Shows the ghrev version.
```sh
ghrev version
```

## Setup
`ghrev` calls the GitHub API and therefore requires credentials. Authentication is automatically picked up from the [GitHub CLI (`gh`)](https://cli.github.com/) login session.
Before using `ghrev`, please do the following.

### 1. Install `gh`
Follow the [official guide](https://github.com/cli/cli#installation) to install it.
For example, Homebrew on macOS, or Scoop / WinGet on Windows.

### 2. Authenticate with `gh auth login`
```sh
gh auth login
```
Follow the interactive prompts to log in to GitHub. The token is stored securely in your OS keychain or equivalent.
You can verify your login status with `gh auth status`.

```sh
$ gh auth status
github.com
  ✓ Logged in to github.com account <your-account> (keyring)
```

### Credential resolution order
`ghrev` resolves the token in the following order of priority:
1. `GH_TOKEN` environment variable
2. `GITHUB_TOKEN` environment variable
3. `gh` CLI login session

In environments where `gh` is unavailable, you can also issue a Personal Access Token (PAT) and set it as `GH_TOKEN`. However, PATs are long-lived and carry a higher leakage risk, so authentication via `gh auth login` is strongly recommended.
```sh
GH_TOKEN=<your-pat> ghrev approval --owner ... --name ... --from ... --to ... --required-approvals ...
```
