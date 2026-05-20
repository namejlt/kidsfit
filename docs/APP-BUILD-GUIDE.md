# KidsFit 小勇士 - App客户端编译安装文档

## 1. 环境要求

| 工具 | 最低版本 | 说明 |
|------|---------|------|
| Flutter SDK | 3.0.0+ | 推荐 3.22+，含 Dart SDK |
| Android Studio | 2023.1+ | 或仅安装 Android SDK |
| Android SDK | API 26+ | compileSdk / targetSdk 34 |
| JDK | 11+ | 推荐 17 |
| Xcode | 15.0+ | 仅 iOS 构建需要 |
| macOS | 13.0+ | 开发环境 |

## 2. 安装 Flutter SDK

### 2.1 macOS（Homebrew）

```bash
brew install flutter
flutter --version
flutter doctor
```

### 2.2 macOS（手动安装）

```bash
# 下载最新稳定版
curl -O https://storage.googleapis.com/flutter_infra_release/releases/stable/macos/flutter_macos_3.22.0-stable.zip

# 解压到指定目录
unzip flutter_macos_3.22.0-stable.zip -d ~/development

# 添加到 PATH
echo 'export PATH="$PATH:$HOME/development/flutter/bin"' >> ~/.zshrc
source ~/.zshrc

# 验证
flutter --version
```

### 2.3 配置 Flutter 国内镜像（加速下载）

```bash
# 编辑 ~/.zshrc，添加以下内容
export PUB_HOSTED_URL=https://pub.flutter-io.cn
export FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn
source ~/.zshrc
```

## 3. 配置 Android SDK

### 3.1 安装 Android Studio

从 [Android Developer](https://developer.android.com/studio) 下载安装。

### 3.2 安装必要 SDK 组件

打开 Android Studio → Settings → SDK Manager：

- SDK Platforms：勾选 Android 14 (API 34)
- SDK Tools：勾选 Android SDK Build-Tools 34、Android SDK Command-line Tools、Android SDK Platform-Tools

### 3.3 接受 Android 许可

```bash
flutter doctor --android-licenses
# 逐个输入 y 接受
```

### 3.4 配置环境变量

```bash
echo 'export ANDROID_HOME=$HOME/Library/Android/sdk' >> ~/.zshrc
echo 'export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools' >> ~/.zshrc
source ~/.zshrc
```

## 4. 项目配置

### 4.1 进入项目目录

```bash
cd kidsfit_app
```

### 4.2 获取依赖

```bash
flutter pub get
```

### 4.3 生成代码（如使用 json_serializable）

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### 4.4 配置 API 地址

编辑 `lib/core/constants/app_constants.dart`：

```dart
/// Android 模拟器访问本机
static const String apiBaseUrl = 'http://10.0.2.2:8001';

/// iOS 模拟器访问本机
static const String apiBaseUrlIOS = 'http://localhost:8001';

/// 真机访问（替换为服务端实际IP）
/// static const String apiBaseUrl = 'http://192.168.1.100:8001';
```

### 4.5 配置 Android 签名

#### Debug 签名（默认）

项目默认使用 debug 签名，无需额外配置。

#### Release 签名

1. 创建签名密钥：

```bash
keytool -genkey -v -keystore ~/kidsfit-release.jks \
  -keyalg RSA -keysize 2048 -validity 10000 \
  -alias kidsfit
```

2. 创建 `android/key.properties`：

```properties
storePassword=你的密钥库密码
keyPassword=你的密钥密码
keyAlias=kidsfit
storeFile=/Users/你的用户名/kidsfit-release.jks
```

3. 修改 `android/app/build.gradle` 中的签名配置（替换 `signingConfigs.debug`）。

## 5. 编译构建

### 5.1 检查环境

```bash
flutter doctor -v
```

确保以下项全部通过：
- Flutter
- Android toolchain
- Android Studio

### 5.2 构建 Debug APK

```bash
flutter build apk --debug
```

输出位置：`build/app/outputs/flutter-apk/app-debug.apk`

### 5.3 构建 Release APK

```bash
flutter build apk --release
```

输出位置：`build/app/outputs/flutter-apk/app-release.apk`

### 5.4 构建 App Bundle（Google Play 发布）

```bash
flutter build appbundle --release
```

输出位置：`build/app/outputs/bundle/release/app-release.aab`

### 5.5 构建 iOS（仅 macOS）

```bash
flutter build ios --release
```

## 6. 安装运行

### 6.1 使用模拟器

```bash
# 列出可用模拟器
flutter emulators

# 启动 Android 模拟器
flutter emulators --launch <emulator_id>

# 运行应用
flutter run
```

### 6.2 使用真机

1. 手机开启 **开发者模式** 和 **USB 调试**
2. USB 连接电脑
3. 确认设备已识别：

```bash
flutter devices
```

4. 安装运行：

```bash
flutter run --release
```

### 6.3 手动安装 APK

将构建好的 APK 文件传输到手机：

```bash
# 通过 adb 安装
adb install build/app/outputs/flutter-apk/app-release.apk

# 或直接将 APK 文件拷贝到手机，点击安装
```

## 7. 项目结构

```
kidsfit_app/
├── android/                     # Android 平台配置
│   ├── app/
│   │   ├── build.gradle         # 应用级构建配置
│   │   └── src/main/AndroidManifest.xml  # 权限声明
│   ├── build.gradle             # 项目级构建配置（含阿里云镜像）
│   └── settings.gradle
├── ios/                         # iOS 平台配置
├── lib/
│   ├── main.dart                # 应用入口
│   ├── app.dart                 # App 根组件
│   ├── core/
│   │   ├── constants/           # 常量（API地址、颜色、尺寸、字符串）
│   │   ├── network/             # 网络层（Dio封装、JWT拦截器、统一响应）
│   │   └── theme/               # 主题配置
│   ├── data/
│   │   └── datasources/         # 远程数据源（4个API调用层）
│   ├── domain/
│   │   └── entities/            # 领域实体（User/ExerciseRecord/VisionRecord/Reward）
│   ├── presentation/
│   │   ├── providers/           # Riverpod 状态管理（4个Provider）
│   │   └── screens/             # UI 页面（儿童端/家长端/认证/视力）
│   └── services/                # 服务层（ML Kit/存储/运动控制器）
├── assets/                      # 静态资源
├── pubspec.yaml                 # 依赖配置
└── analysis_options.yaml        # 代码规范
```

## 8. 核心依赖说明

| 依赖 | 版本 | 用途 |
|------|------|------|
| flutter_riverpod | ^2.4.9 | 状态管理 |
| dio | ^5.3.3 | HTTP 网络请求 |
| google_mlkit_pose_detection | ^0.11.0 | AI 骨骼检测 |
| camera | ^0.10.5+9 | 摄像头控制 |
| hive + hive_flutter | ^2.2.3 | 本地离线存储 |
| go_router | ^13.0.1 | 路由导航 |
| fl_chart | ^0.65.0 | 数据图表 |
| geolocator | ^10.1.0 | GPS 定位（户外时间） |
| permission_handler | ^11.1.0 | 权限管理 |

## 9. 常见问题

### 9.1 Flutter 命令找不到

```bash
which flutter
# 如果为空，检查 PATH 配置
echo $PATH | grep flutter
```

### 9.2 Gradle 下载慢

项目已配置阿里云镜像（`android/build.gradle`）：

```groovy
maven { url 'https://maven.aliyun.com/repository/google' }
maven { url 'https://maven.aliyun.com/repository/central' }
maven { url 'https://maven.aliyun.com/repository/public' }
```

### 9.3 依赖下载失败

```bash
# 清理缓存
flutter pub cache clean
flutter clean
flutter pub get
```

### 9.4 Android 许可未接受

```bash
flutter doctor --android-licenses
```

### 9.5 构建失败：minSdkVersion 冲突

确保 `android/app/build.gradle` 中 `minSdkVersion 26`。

### 9.6 ML Kit 模型下载

首次运行骨骼检测时，Google ML Kit 会自动下载模型文件。确保网络通畅。

### 9.7 真机无法连接后端

- 确保手机和服务端在同一局域网
- 修改 `app_constants.dart` 中的 `apiBaseUrl` 为服务端实际 IP
- 检查防火墙设置

## 10. 性能指标

| 指标 | 目标值 |
|------|--------|
| 冷启动时间 | < 3秒 |
| 页面切换 | < 300ms |
| AI骨骼识别延迟 | < 200ms |
| 运行时内存峰值 | < 300MB |
| 安装包大小 | < 150MB |
| 崩溃率 | < 0.1% |
