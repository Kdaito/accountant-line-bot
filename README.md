# Accountant Line Bot

Accountant Line Bot は、レシート画像を簡単にスプレッドシートにまとめることができる、LINE メッセージングボットです。ボットにレシート画像を送信すると、Chat GCP の API を使用して解析し、Google スプレッドシートとして作成後、指定の Drive フォルダに保存します。

この README では、以下の内容について説明します。

1. セットアップ
2. ローカルでの実行
3. デプロイ

注意: `3. デプロイ`を実行するためには、`2. ローカルでの実行`を完了している必要があります。

## 1. セットアップ

1. 外部サービスのセットアップ
2. ローカル実行用の環境変数の設定

### 1-1. 外部サービスのセットアップ

Accountant Line Bot で使用している外部サービスとその用途は以下の表の通りです。

| サービス                   | 用途                                   |
| -------------------------- | -------------------------------------- |
| Google Cloud Platform(GCP) | シートの作成と保存、デプロイ           |
| ChatGPT(Open AI)           | レシート画像の解析                     |
| LINE Messaging API         | レシート画像の受け取り、メッセージ返信 |

このセクションではこれらのセットアップと環境変数への設定について説明します。

#### Google Cloud Platform(GCP)

1. プロジェクトの作成
2. API の有効化: 以下の API を使用します。

- Google Sheets API
- Google Drive API
- Secret Manager API
- Cloud Run Admin API
- Artifact Registry API

3. OAuth 認証

- 「API とサービス」→ 「認証情報」から、Oauth2.0 の認証情報を作成し、作成したクライアントの認証情報を JSON としてダウンロードしてください。

4. JSON の内容を、`.env`ファイルに`ENV_CREDENTIALS_JSON`として設定してください。

#### ChatGPT(Open AI)

1. アカウントの作成

- [OpenAI Platform](https://platform.openai.com/docs/overview)にログインし、プロジェクトを作成してください

2. API キーの取得

- 「API keys」から、「SECRET KEY」を取得し、`.env`ファイルに`ENV_GPT_API_KEY`として設定してください。

#### LINE Messaging API

1. LINE Developers コンソール

- [LINE Developers](https://developers.line.biz/ja/)にアクセスし、プロバイダーを作成します。

2. チャネルの作成

- Messaging API チャネルを作成します。

3. チャネルシークレットとアクセストークンの取得

- チャネルシークレットとアクセストークンを取得し、それぞれ`.env`ファイルに、`ENV_LINE_CHANNEL_SECRET`および`ENV_LINE_CHANNEL_TOKEN`として設定します。

### 1-2. ローカル実行用の環境変数の設定

以下の通りに環境変数を保存する`.env`ファイルを作成してください。

```sh:.env
# .env

ENV_LINE_CHANNEL_SECRET # LINE Messaging APIのチャンネルシークレット
ENV_LINE_CHANNEL_TOKEN # LINE Messaging APIのチャンネルトークン
ENV_FOLDER_ID # 作成したスプシを保存したいGoogleDriveのフォルダID
ENV_GPT_API_URL # Chat GPT APIのURL
ENV_GPT_API_KEY # Chat GPT APIのAPIキー(シークレットキー)
ENV_PORT # ローカルで起動時のポート
ENV_CREDENTIALS_JSON # GCPのOAuth2クライアントの認証情報
ENV_TOKEN_JSON # GCPのOAuth2クライアントのアクセストークン(初回の認証時に作成される。後ほど設定)
```

## 2. ローカルでの実行
1. ローカルで実行する
2. 環境変数`ENV_TOKEN_JSON`を設定する

### 2-1. ローカルで実行する

以下のコマンドで、ローカルでアプリを起動します。(goがインストールされている必要があります。)

```sh
./run-local.sh
```

別のターミナルを開いて、以下のコマンドを実行すると、ngnrokを用いてline messaging apiのhooks用のurlを取得することができます。

```sh
./run-server.sh
```

ここに表示されたurlをline messaging apiのhooksに`[url]/callback`として設定することで、LINE botにメッセージが送信された場合、ローカルで実行しているアプリケーションの`callback`を呼び出すことができます。

### 2-2. 環境変数`ENV_TOKEN_JSON`を設定する

初回起動時に、GCPのOAuth2による認証が実行され、コマンドラインに認証用のリンクが表示されます。指示に従って認証に成功すると、token.jsonが作成されます。このtoken.jsonを新たに環境変数として保存してください。

```sh:.env
# .env
ENV_TOKEN_JSON # GCPのOAuth2クライアントのアクセストークン
```

これにより、2回目の起動以降は`ENV_TOKEN_JSON`から認証情報を取得するようになります。

## 3. デプロイ

1. 環境変数の設定
2. デプロイ

### 3-1. 環境変数の設定

デプロイ時にコンテナにSecret Managerから環境変数を注入するため、ローカルで起動時に作成していた環境変数を、Secret Managerに登録します。

CLIやGCPのコンソールなどから、`.env`の以下の項目を登録してください。

- LINE_CHANNEL_SECRET
- LINE_CHANNEL_TOKEN
- FOLDER_ID
- GPT_API_URL
- GPT_API_KEY
- CREDENTIALS_JSON
- TOKEN_JSON

GCPにデプロイするための環境変数を`.env.gcp`を作成して設定します。

```sh:.env.gcp
# .env.gcp

ENV_LOCATION # GCPのロケーション
ENV_PROJECT_ID # GCPのプロジェクトID
ENV_REPOSITORY # Artifact Registryのリポジトリ
ENV_IMAGE # Artifact Registryに登録する際のイメージ名
ENV_SERVICE # Cloud runのデプロイするサービス

SECRET_LINE_CHANNEL_SECRET # [LINE_CHANNEL_SECRETのsecret managerの名前]:latest
SECRET_LINE_CHANNEL_TOKEN # [LINE_CHANNEL_TOKENのsecret managerの名前]:latest
SECRET_FOLDER_ID # [FOLDER_IDのsecret managerの名前]:latest
SECRET_GPT_API_URL # [GPT_API_URLのsecret managerの名前]:latest
SECRET_GPT_API_KEY # [GPT_API_KEYのsecret managerの名前]:latest
SECRET_CREDENTIALS_JSON # [CREDENTIALS_JSONのsecret managerの名前]:latest
SECRET_TOKEN_JSON # [TOKEN_JSONのsecret managerの名前]:latest
```

### 3-2. デプロイ

Artifact RegistryとCloud runサービスの準備、確認が完了した後に、以下のコマンドを実行してください。

```sh
./run-deploy.sh
```

実行されたCloud runのurlを`url/callback`としてLINE Messaging APIのhooks urlとして登録すれば完了です。