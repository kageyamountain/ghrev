# ghrev

日本語 | [English](README.en.md)

`ghrev` は GitHub リポジトリの Pull Request レビューに関するメトリクスを計測する CLI ツールです。  
コードレビューを様々な切り口で可視化することで、レビュー文化の改善や、コードレビューに関わる時間の見える化を支援します。

## サブコマンド
### `approval`
指定期間内に作成された PR を対象に、PR オープンから指定件数の Approve が揃うまでにかかった時間を集計します。
```sh
ghrev approval \
  --owner <organization-or-user> \
  --name <repository> \
  --from <YYYYMMDD> \
  --to <YYYYMMDD> \
  --required-approvals <N> \
  --ignore-labels <label1,label2,...>
```

| オプション | 必須 | 説明 |
| --- | :---: | --- |
| `--owner` | ✓ | 対象リポジトリのオーナー(ユーザー名または Organization 名) |
| `--name` | ✓ | 対象リポジトリ名 |
| `--from` | ✓ | 集計開始日(`YYYYMMDD` 形式) |
| `--to` | ✓ | 集計終了日(`YYYYMMDD` 形式) |
| `--required-approvals` | ✓ | レビュー完了とみなす Approve の件数(1 以上の整数) |
| `--ignore-labels` |   | 集計対象から除外するラベル(カンマ区切り) |

### `first-review`
指定期間内に作成された PR を対象に、PR オープンから最初のレビュー反応(approve / changes_requested / commented のいずれか)までにかかった時間を集計します。DISMISSED されたレビューは反応とみなしません。
```sh
ghrev first-review \
  --owner <organization-or-user> \
  --name <repository> \
  --from <YYYYMMDD> \
  --to <YYYYMMDD> \
  --ignore-labels <label1,label2,...>
```

| オプション | 必須 | 説明 |
| --- | :---: | --- |
| `--owner` | ✓ | 対象リポジトリのオーナー(ユーザー名または Organization 名) |
| `--name` | ✓ | 対象リポジトリ名 |
| `--from` | ✓ | 集計開始日(`YYYYMMDD` 形式) |
| `--to` | ✓ | 集計終了日(`YYYYMMDD` 形式) |
| `--ignore-labels` |   | 集計対象から除外するラベル(カンマ区切り) |

### `help`
利用可能なサブコマンドのヘルプを表示します。
```sh
ghrev help
```

### `version`
ghrevのバージョンを表示します。
```sh
ghrev version
```

## セットアップ
`ghrev` は GitHub の API を呼び出すため、認証情報を必要とします。認証は [GitHub CLI (`gh`)](https://cli.github.com/) のログインセッションから自動で取得します。
そのため、利用前に以下を実施してください。

### 1. `gh` のインストール
[公式のガイド](https://github.com/cli/cli#installation)を参考にインストールしてください。  
MacであればHomebrew、WindowsであればScoopやWinGetなど。

### 2. `gh auth login` で認証
```sh
gh auth login
```
対話的なプロンプトに従って GitHub にログインしてください。トークンは OS のキーチェーン等に安全に保管されます。  
ログイン状態は `gh auth status` で確認できます。

```sh
$ gh auth status
github.com
  ✓ Logged in to github.com account <your-account> (keyring)
```

### 認証情報の解決順序
`ghrev` は以下の優先順位でトークンを解決します。
1. 環境変数 `GH_TOKEN`
2. 環境変数 `GITHUB_TOKEN`
3. `gh` CLI のログインセッション

`gh` を利用できない環境では、Personal Access Token (PAT) を発行して `GH_TOKEN` にセットすることでも動作しますが、PAT は長命で漏洩リスクが高いため、`gh auth login` による認証を強く推奨します。
```sh
GH_TOKEN=<your-pat> ghrev approval --owner ... --name ... --from ... --to ... --required-approvals ...
```
