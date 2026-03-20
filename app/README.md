# ohome

Flutter 客户端应用。

## 环境配置

环境文件：
- `assets/env/dev.json`
- `assets/env/prod.json`

运行时通过 `APP_ENV` 选择环境：
- 开发：`flutter run --dart-define=APP_ENV=dev`
- 生产：`flutter run --dart-define=APP_ENV=prod`
- 打包：`flutter build apk --split-per-abi --dart-define=APP_ENV=prod`
- 一键打包：`build_prod_split_abi.bat`（自动递增本地 `.build_number`，也可 `build_prod_split_abi.bat 42` 指定 build number）

运行在浏览器 
- 开发：`flutter run -d chrome`
