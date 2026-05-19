# KidsFit 小勇士 Flutter项目环境配置指南

## 环境要求

- Flutter SDK >= 3.0.0
- Dart SDK >= 3.0.0
- Android Studio / Android SDK
- Xcode (仅iOS开发)
- macOS 10.14+ (使用Homebrew安装)

## 安装步骤

### 方式一：使用Homebrew安装（推荐 macOS）

```bash
# 更新Homebrew
brew update

# 安装Flutter
brew install flutter

# 验证安装
flutter --version

# 运行flutter doctor检查环境
flutter doctor
```

### 方式二：手动下载安装

1. 从 [Flutter官网](https://flutter.dev/docs/get-started/install/macos) 下载最新稳定版SDK
2. 解压到指定目录（如 `/Users/tynam/development/flutter`）
3. 添加到PATH：

```bash
# 编辑shell配置文件
nano ~/.zshrc  # 或 ~/.bashrc

# 添加以下内容
export PATH="$PATH:/Users/tynam/development/flutter/bin"

# 使配置生效
source ~/.zshrc

# 验证安装
flutter --version
```

### 方式三：使用FVM（Flutter Version Manager）

```bash
# 安装FVM
brew install fvm

# 安装指定版本Flutter
fvm install 3.22.0

# 在项目中使用
fvm flutter pub get
```

## Android环境配置

### 1. 安装Android Studio

从 [Android Developer官网](https://developer.android.com/studio) 下载安装Android Studio

### 2. 配置Android SDK

```bash
# 设置Android SDK路径
flutter config --android-sdk /Users/tynam/Library/Android/sdk

# 接受Android许可
flutter doctor --android-licenses
```

### 3. 配置Android环境变量

```bash
# 编辑配置文件
nano ~/.zshrc

# 添加以下内容
export ANDROID_HOME=/Users/tynam/Library/Android/sdk
export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools
```

## 项目依赖安装

```bash
# 进入项目目录
cd kidsfit_app

# 获取依赖
flutter pub get

# 运行代码生成（如有需要）
flutter pub run build_runner build --delete-conflicting-outputs
```

## 构建项目

### Debug版本

```bash
# 连接设备后
flutter run

# 或构建debug APK
flutter build apk --debug
```

### Release版本

```bash
# 构建release APK
flutter build apk --release

# APK位置: build/app/outputs/flutter-apk/app-release.apk
```

## 常见问题

### 1. Flutter命令找不到

确保Flutter已添加到PATH：
```bash
echo $PATH | grep flutter
```

如果没有，参考上面的安装步骤添加。

### 2. Android许可未接受

```bash
flutter doctor --android-licenses
```

### 3. 依赖下载失败

尝试更换Flutter镜像源：
```bash
flutter config --global pub https://pub.flutter-io.cn
```

### 4. iOS模拟器无法启动

```bash
# 列出可用模拟器
xcrun simctl list devices

# 启动指定模拟器
open -a Simulator
```

## 开发工具推荐

- **VS Code**: 轻量级IDE，Flutter插件优秀
- **Android Studio**: 完整Android开发环境
- **IntelliJ IDEA**: 强大的Dart/Flutter支持

### 推荐VS Code扩展

- Flutter
- Dart
- Awesome Flutter Snippets
- Flutter Intl (国际化)

## 下一步

1. 确保所有工具安装完成
2. 运行 `flutter doctor` 检查环境
3. 配置 Firebase（如需推送通知）
4. 开始开发！
