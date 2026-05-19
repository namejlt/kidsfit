# KidsFit 小勇士 - APK构建指南

## ⚠️ 重要提示

**Flutter SDK需要用户手动安装！**

由于Flutter SDK（约2.2GB）下载速度较慢，建议用户使用以下方式之一安装：

### 方式一：使用自动化脚本（推荐）

```bash
# 1. 进入项目目录
cd kidsfit_app

# 2. 运行环境配置脚本
./setup_env.sh

# 3. 按照提示完成Flutter安装
```

### 方式二：使用Homebrew

```bash
# 安装Flutter
brew install flutter

# 验证安装
flutter --version
```

### 方式三：手动下载

1. 从 [Flutter官网](https://flutter.dev/docs/get-started/install/macos) 下载最新SDK
2. 解压到 `~/development/flutter`
3. 添加到PATH：`export PATH="$PATH:$HOME/development/flutter/bin"`
4. 运行 `flutter doctor` 验证

## 环境要求

- Flutter SDK >= 3.0.0
- Dart SDK >= 3.0.0
- Android Studio / Android SDK (API Level 26+)
- macOS 10.14+

## 构建APK

### 快速构建

```bash
# 进入项目目录
cd kidsfit_app

# 运行构建脚本
./build_app.sh
```

### 手动构建

#### Debug版本

```bash
# 获取依赖
flutter pub get

# 构建Debug APK
flutter build apk --debug
```

**APK文件位置**: `build/app/outputs/flutter-apk/app-debug.apk`

#### Release版本

```bash
# 获取依赖
flutter pub get

# 构建Release APK
flutter build apk --release
```

**APK文件位置**: `build/app/outputs/flutter-apk/app-release.apk`

#### App Bundle（Google Play发布用）

```bash
# 构建App Bundle
flutter build appbundle --release
```

**Bundle文件位置**: `build/app/outputs/bundle/release/app-release.aab`

## 常见问题

### 1. Flutter命令找不到

```bash
# 检查Flutter是否在PATH中
which flutter

# 如果没有，手动添加到 ~/.zshrc
echo 'export PATH="$PATH:$HOME/development/flutter/bin"' >> ~/.zshrc
source ~/.zshrc
```

### 2. Android许可未接受

```bash
flutter doctor --android-licenses
```

### 3. 依赖下载失败

```bash
# 更换Flutter镜像源
flutter config --global pub https://pub.flutter-io.cn

# 重新获取依赖
flutter pub get
```

### 4. Android SDK未找到

```bash
# 设置Android SDK路径
flutter config --android-sdk /Users/你的用户名/Library/Android/sdk
```

## 构建验证

构建成功后，可以使用以下命令验证APK：

```bash
# 查看APK信息
ls -lh build/app/outputs/flutter-apk/

# 检查APK签名
keytool -printcert -jarfile build/app/outputs/flutter-apk/app-release.apk
```

## 下一步

1. 确保所有环境检查通过（`flutter doctor`）
2. 安装Firebase配置（用于推送通知，可选）
3. 将APK安装到设备测试
4. 部署到应用商店

## 技术支持

- 详细环境配置：参见 [FLUTTER_SETUP.md](./FLUTTER_SETUP.md)
- 项目文档：参见 [../docs/README.md](../docs/README.md)
- 需求文档：参见 [../docs/PRD.md](../docs/PRD.md)
- 技术设计：参见 [../docs/TECH-DESIGN.md](../docs/TECH-DESIGN.md)
