工事中です。

## リリース方法

最新のコミットに対してタグを貼ってください。（自動化したいです）
そしてそのタグをプッシュすると自動的にリリースノートが生成されます。
Assets として(Win, Mac, Linux) x (x86_64, arm64, i386)のビルド成果物が生成されます。
また ghcr.io にコンテナイメージがプッシュされます。これは x86_64 と arm64 のみに対応です。

```bash
git tag v(SEM_VER)
git push --tags
```