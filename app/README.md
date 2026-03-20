# ohome

Flutter 客户端应用。

## 环境配置

环境文件：
- `assets/env/dev.json`
- `assets/env/prod.json`

运行时通过 `APP_ENV` 选择环境：
- 开发：`flutter run --dart-define=APP_ENV=dev`
- 生产：`flutter run --dart-define=APP_ENV=prod`
- 打包：`flutter build apk --dart-define=APP_ENV=prod`
- 一键打包：`bash build_prod.sh`
- 发布版本号以 `pubspec.yaml` 的 `version: x.y.z+n` 为准
- GitHub 发版建议使用与 `versionName` 对应的 tag，例如 `version: 0.0.3+1` 时发布 `v0.0.3`

运行在浏览器 
- 开发：`flutter run -d chrome`
