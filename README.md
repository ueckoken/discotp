# discotp

## En

### What is this?

this application is OTP(One Time Password) client works as Discord app.

## Jp

### これはなに

これは Discord のアプリケーションとして使えるワンタイムパスワードクライアントです。

### 設定方法

以下の環境変数を設定してください。

| キー                        | 値                                                                                               |
| --------------------------- | ------------------------------------------------------------------------------------------------ |
| `DISCORD_APP_TOKEN`         | Discord の Bot トークン                                                                          |
| `DISCORD_GUILD_ID`          | Discord の GuildID です。                                                                        |
| `ALLOWRD_REPLY_CHANNEL_IDS` | メッセージの返答を許可するチャンネルの ID です。分割文字は`,`で複数設定できます                  |
| `GOOGLE_TOTP_TOKEN`         | Google の TOTP トークンです。「QR コードを読み取れない」のような場所をクリックすると取得できます |

### リリース方法

最新のコミットに対してタグを貼ってください。（自動化したいです）
そしてそのタグをプッシュすると自動的にリリースノートが生成されます。
Assets として(Win, Mac, Linux) x (x86_64, arm64, i386)のビルド成果物が生成されます。
また ghcr.io にコンテナイメージがプッシュされます。これは x86_64 と arm64 のみに対応です。

```bash
git tag v(SEM_VER)
git push --tags
```
