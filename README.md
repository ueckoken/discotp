# discotp

## En

### What is this?

this application is OTP(One Time Password) client works as Discord app.

## Jp

### これはなに

これは Discord のアプリケーションとして使えるワンタイムパスワードクライアントです。

### 設定方法

以下の環境変数を設定してください。

#### `DISCORD_APP_TOKEN`

(required)

Discord の Bot トークンを指定します。

#### `DISCORD_GUILD_ID`

(required)

メッセージを返答する Discord の GuildID を指定します。

#### `ALLOWED_REPLY_CHANNEL_IDS`

(required)

メッセージの返答を許可するチャンネルの ID です。分割文字を`,`として複数設定できます。

#### `TOTP_TOKENS`

(required)
TOTP トークンです。`サービス名1:トークン,サービス名2:トークン`のように指定してください。トークンの間にあるスペースは入れても入れなくても問題ないです。「QR コードを読み取れない」のような場所をクリックすると取得できます

#### `IS_DEVELOPMENT`

(optional)
`1`を指定すると開発モードになって、ログが json 形式ではなくなります。

