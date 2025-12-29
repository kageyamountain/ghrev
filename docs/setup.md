## 開発環境セットアップ
Windows環境を例に書いているが、Macもhomebrewで `mise` をインストールすれば同様の手順でセットアップ可能なはず。

### 前提

- Git
- WSL2 や Git Bash などのUnix系シェル実行環境
- Docker Desktop
- Pathは適宜通すこと

### インストール手順
パッケージ管理ツールの `Scoop` をインストールする。  
`Windows PowerShell` で以下のコマンドを実行する。
```shell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```
これ以降のコマンドはUnix系シェルでOK

環境管理ツールの `mise` をインストールする。
```shell
scoop install mise
```

ターミナルの環境変数に追加（git bashの場合は `~/.bashrc`）
```
# mise
export PATH="$HOME/scoop/apps/mise/current/bin:$PATH"
export PATH="$HOME/AppData/Local/mise/shims:$PATH"
```

環境変数設定を反映
```shell
source ~/.bashrc
```

miseの実験的機能の有効化
```shell
mise settings set experimental true

# 確認
mise settings get experimental
```

shimsを有効化
```shell
mise reshim

# shimsが作成されているか確認
ls ~/AppData/Local/mise/shims

# miseが動作するか確認
mise --version
```

インストール
```shell
cd /path/to/your/project
mise install

# 動作確認
go version
```

go mod tidyで依存関係の整理（コンテナ環境なので必須ではないが、IDE補完などの開発体験的に推奨）
```shell
go mod tidy
```

## ローカル開発環境でのghrev実行方法
```shell
# コンテナ起動
task compose-up
# ghrev実行
task exec -- --mode two-approve
```
基本的なコマンドは `task` コマンドで用意しているので `Taskfile.yml` 参照
