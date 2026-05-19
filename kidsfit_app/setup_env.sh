#!/bin/bash

# KidsFit Flutter Environment Setup Script
# 自动配置Flutter开发环境

set -e

echo "========================================"
echo "  KidsFit Flutter 环境配置脚本"
echo "========================================"
echo ""

# 检查是否以root运行
if [ "$EUID" -eq 0 ]; then
  echo "❌ 请不要以root用户运行此脚本"
  exit 1
fi

# 1. 检查Homebrew
echo "📦 检查Homebrew..."
if ! command -v brew &> /dev/null; then
  echo "⚠️  Homebrew未安装，正在安装..."
  /bin/bash -c "$(curl -fsSL https://mirrors.ustc.edu.cn/brew.sh)"
else
  echo "✅ Homebrew已安装"
fi

# 2. 检查Flutter
echo ""
echo "🔍 检查Flutter..."
if command -v flutter &> /dev/null; then
  FLUTTER_VERSION=$(flutter --version)
  echo "✅ Flutter已安装: $FLUTTER_VERSION"
  flutter doctor
else
  echo "⚠️  Flutter未安装"
  echo ""
  echo "请选择安装方式："
  echo "1. 使用Homebrew安装 (brew install flutter)"
  echo "2. 手动下载：https://flutter.dev/docs/get-started/install/macos"
  echo ""
  read -p "输入您的选择 [1/2]: " choice

  case $choice in
    1)
      echo "使用Homebrew安装Flutter..."
      brew install flutter
      ;;
    2)
      echo "请手动下载Flutter SDK并解压到指定目录"
      echo "然后将以下内容添加到 ~/.zshrc："
      echo "export PATH=\"\$PATH:/path/to/flutter/bin\""
      ;;
    *)
      echo "无效选择，请手动安装Flutter"
      ;;
  esac
fi

# 3. 检查Android SDK
echo ""
echo "📱 检查Android SDK..."
if [ -z "$ANDROID_HOME" ] && [ -z "$ANDROID_SDK_ROOT" ]; then
  if [ -d "/Users/$USER/Library/Android/sdk" ]; then
    export ANDROID_HOME="/Users/$USER/Library/Android/sdk"
    echo "✅ Android SDK已找到: $ANDROID_HOME"
  else
    echo "⚠️  Android SDK未找到"
    echo "请安装Android Studio: https://developer.android.com/studio"
  fi
else
  echo "✅ Android SDK已配置"
fi

# 4. 配置环境变量
echo ""
echo "⚙️  配置环境变量..."
SHELL_RC="$HOME/.zshrc"
if [ -f "$SHELL_RC" ]; then
  if ! grep -q "flutter" "$SHELL_RC"; then
    echo 'export PATH="$PATH:$HOME/development/flutter/bin"' >> "$SHELL_RC"
    echo "✅ 已添加到 $SHELL_RC"
  fi

  if ! grep -q "ANDROID_HOME" "$SHELL_RC"; then
    echo 'export ANDROID_HOME="$HOME/Library/Android/sdk"' >> "$SHELL_RC"
    echo "✅ Android SDK路径已添加到 $SHELL_RC"
  fi

  source "$SHELL_RC"
  echo "✅ 环境变量已生效"
else
  echo "⚠️  未找到 $SHELL_RC，请手动配置"
fi

# 5. 运行flutter doctor
echo ""
echo "🔧 运行flutter doctor..."
flutter doctor

# 6. 获取项目依赖
echo ""
echo "📥 获取项目依赖..."
if [ -d "kidsfit_app" ]; then
  cd kidsfit_app
  flutter pub get
  echo "✅ 依赖安装完成"
else
  echo "⚠️  未找到kidsfit_app目录"
fi

echo ""
echo "========================================"
echo "  环境配置完成！"
echo "========================================"
echo ""
echo "下一步："
echo "1. 确保所有环境检查通过"
echo "2. 运行: flutter build apk --release"
echo "3. APK将生成在: build/app/outputs/flutter-apk/"
echo ""
