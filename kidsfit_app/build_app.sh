#!/bin/bash

# KidsFit 快速构建脚本

set -e

echo "========================================"
echo "  KidsFit 小勇士 - APK构建脚本"
echo "========================================"
echo ""

# 检查Flutter
echo "🔍 检查Flutter环境..."
if ! command -v flutter &> /dev/null; then
  echo "❌ Flutter未安装，请先运行 setup_env.sh"
  exit 1
fi

flutter --version
echo ""

# 检查Android SDK
echo "📱 检查Android SDK..."
flutter doctor --android-licenses 2>/dev/null || true
flutter doctor
echo ""

# 进入项目目录
if [ -d "kidsfit_app" ]; then
  cd kidsfit_app
fi

# 清理旧的构建文件
echo "🧹 清理旧的构建文件..."
flutter clean

# 获取依赖
echo "📥 获取项目依赖..."
flutter pub get

# 代码生成（如需要）
echo "⚙️  运行代码生成..."
flutter pub run build_runner build --delete-conflicting-outputs 2>/dev/null || true

echo ""
echo "🔨 开始构建..."
echo ""

# 构建Debug版本
if [ "$1" == "debug" ]; then
  echo "构建Debug APK..."
  flutter build apk --debug
  echo ""
  echo "✅ Debug APK构建完成！"
  echo "📁 文件位置: build/app/outputs/flutter-apk/app-debug.apk"
fi

# 构建Release版本（默认）
if [ "$1" == "release" ] || [ -z "$1" ]; then
  echo "构建Release APK..."
  flutter build apk --release
  echo ""
  echo "✅ Release APK构建完成！"
  echo "📁 文件位置: build/app/outputs/flutter-apk/app-release.apk"
fi

# 构建Bundle
if [ "$1" == "bundle" ]; then
  echo "构建App Bundle..."
  flutter build appbundle --release
  echo ""
  echo "✅ App Bundle构建完成！"
  echo "📁 文件位置: build/app/outputs/bundle/release/app-release.aab"
fi

echo ""
echo "========================================"
echo "  构建完成！"
echo "========================================"
