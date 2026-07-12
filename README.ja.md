<h1 align='center'>パルワールドサーバーツール</h1>

<p align="center">
  <a href="/README.md">简体中文</a> | <a href="/README.en.md">English</a> | <strong>日本語</strong>
</p>

<p align='center'>セーブ解析、公式 REST API、RCON を利用して Palworld 専用サーバーを管理する Web ダッシュボードです。</p>

![PC](./docs/img/pst-ja-1.png)

## 機能

- プレイヤー、ギルド、パル、インベントリの表示
- サーバー情報、メトリクス、オンラインプレイヤー
- キック、BAN、ブロードキャスト、正常終了
- マップとホワイトリスト管理
- RCON コマンドと定期実行
- セーブ同期、自動バックアップ、バックアップ管理
- PC・モバイル対応 UI
- 管理モード内の PST 設定画面

業務データは `pst.db`、PST 設定と管理者認証情報は別の `config.db` に保存されます。設定をリセットしても、プレイヤー、ギルド、RCON、バックアップ記録には影響しません。

## 公式 REST API と RCON の有効化

PST には Palworld サーバーの公式 REST API が必要です。RCON 機能を使う場合は RCON も有効にしてください。ゲームサーバーを停止し、[Pal-Conf](https://pal-conf.bluefissure.com/) で `PalWorldSettings.ini` または `WorldOption.sav` を設定します。ゲームサーバーの `AdminPassword` を設定してから REST API と RCON を有効にします。

![ADMIN](./docs/img/admin-ja.png)

![RCON_REST](./docs/img/rest-rcon-ja.png)

## インストール

`Level.sav` の解析時には短時間に約 1～3 GB のメモリを使用します。

### リリースファイル

1. [GitHub Releases](https://github.com/zaigie/palworld-server-tool/releases) から OS とアーキテクチャに合うファイルをダウンロードして展開します。
2. Linux/macOS では `pst` と `sav_cli` に実行権限を付けて `./pst` を実行します。Windows では `start.bat` または PowerShell から `.\pst.exe` を実行します。
3. `http://127.0.0.1:8080` または `http://サーバーアドレス:8080` を開き、PST Web 管理者を作成して設定画面を入力します。

初回起動はポート `8080` を使用します。Web/TLS 設定またはタスク間隔を変更し、画面に再起動が必要と表示された場合だけ PST を再起動してください。

> [!IMPORTANT]
> PST は `config.yaml`、`-config` 引数、PST 設定用環境変数を読み込みません。アップグレード時は以前の値を Web 設定へ手動でコピーし、古いファイルと変数を削除してください。

### Docker：ローカルセーブ

永続化する 2 つのデータベースファイルを先に作成します。

```bash
touch pst.db config.db
```

```bash
docker run -d --name pst \
  -p 8080:8080 \
  -v /path/to/your/Pal/Saved:/game \
  -v ./backups:/app/backups \
  -v ./pst.db:/app/pst.db \
  -v ./config.db:/app/config.db \
  jokerwho/palworld-server-tool:latest
```

PST 設定で「ローカルディレクトリ」を選択し、`/game` を指定します。RCON と REST のアドレスはコンテナから到達可能である必要があります。

### pst-agent：リモートセーブ

ゲームサーバー側で `pst-agent` を起動します。

```bash
docker run -d --name pst-agent \
  -p 8081:8081 \
  -v /path/to/your/Pal/Saved:/game \
  -e SAVED_DIR="/game" \
  jokerwho/palworld-server-tool-agent:latest
```

PST 本体には設定環境変数を渡しません。PST 設定で「pst-agent」を選択し、`http://ゲームサーバーアドレス:8081/sync` を入力します。ネイティブ版とコマンドラインオプションは [pst-agent ガイド](./README.agent.ja.md) を参照してください。

## 初回アクセスと設定

1. 最初の訪問者が PST Web 管理者パスワードを作成します。このパスワードは PST ダッシュボード専用で、ゲームサーバーの `AdminPassword` ではありません。
2. 初期化は一度だけ成功します。他の人が先に設定した場合は PST を停止し、`config.db` を削除して再起動してください。`pst.db` には影響しません。
3. サーバー側ファイルブラウザーでローカルディレクトリを選ぶか、`pst-agent` URL を入力します。
4. セーブ元と RCON には「未設定 / エラー / 正常」の状態が表示されます。RCON テストは公式の読み取り専用 `Info` コマンドを使用し、サーバー状態を変更しません。
5. セーブ元、RCON、REST、メッセージ、管理設定、パスワードは保存後すぐに反映されます。Web/TLS とタスク間隔のみ再起動が必要で、対象項目は画面に表示されます。

次の旧 PST 設定方法は互換パスなしで削除されました。

- `config.yaml`
- `-config` コマンドライン引数
- `WEB__*`、`RCON__*`、`REST__*`、`SAVE__*`、`TASK__*`、`MANAGE__*` 環境変数

`pst-agent` 自体は、引き続きディレクトリ指定オプションと `SAVED_DIR` を使用できます。

## API ドキュメント

- [APIFox ドキュメント](https://q4ly3bfcop.apifox.cn/)
- ローカル Swagger：`http://127.0.0.1:8080/swagger/index.html`

## ライセンス

[Apache License 2.0](LICENSE) に基づいて提供されます。
