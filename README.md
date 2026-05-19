# KidsFit 小勇士 - Flutter App

## 项目介绍

KidsFit 小勇士是一款面向3-12岁儿童的AI运动健康App，通过AI骨骼识别、游戏化训练、视力健康管理等功能，帮助儿童建立运动习惯、保护视力健康。

## 技术栈

- **框架**: Flutter 3.x
- **状态管理**: Riverpod
- **本地存储**: Hive + SQLite
- **AI识别**: Google ML Kit Pose Detection
- **架构**: 整洁架构 (Clean Architecture)

## 项目结构

```
lib/
├── main.dart
├── app.dart
├── core/                      # 核心模块
│   ├── constants/             # 常量定义
│   │   ├── app_colors.dart
│   │   ├── app_strings.dart
│   │   └── app_dimensions.dart
│   ├── theme/                # 主题配置
│   │   └── app_theme.dart
│   ├── router/               # 路由配置
│   │   └── app_router.dart
│   ├── utils/                # 工具类
│   │   ├── validators.dart
│   │   └── formatters.dart
│   └── errors/               # 错误处理
│       └── exceptions.dart
│
├── data/                      # 数据层
│   ├── datasources/          # 数据源
│   │   ├── local/            # 本地数据源
│   │   │   ├── hive_datasource.dart
│   │   │   └── sqlite_datasource.dart
│   │   └── remote/           # 远程数据源
│   │       └── api_datasource.dart
│   ├── models/               # 数据模型
│   │   ├── user_model.dart
│   │   ├── exercise_model.dart
│   │   ├── vision_model.dart
│   │   └── reward_model.dart
│   └── repositories/          # 仓储实现
│       ├── user_repository_impl.dart
│       ├── exercise_repository_impl.dart
│       └── vision_repository_impl.dart
│
├── domain/                    # 领域层（核心业务逻辑）
│   ├── entities/             # 领域实体
│   │   ├── user.dart
│   │   ├── exercise_record.dart
│   │   ├── training_plan.dart
│   │   ├── vision_record.dart
│   │   └── reward.dart
│   ├── repositories/         # 仓储接口
│   │   ├── user_repository.dart
│   │   ├── exercise_repository.dart
│   │   └── vision_repository.dart
│   └── usecases/             # 用例
│       ├── auth/
│       │   ├── login_usecase.dart
│       │   └── register_usecase.dart
│       ├── exercise/
│       │   ├── get_training_plan_usecase.dart
│       │   └── submit_exercise_usecase.dart
│       └── vision/
│           ├── get_vision_trend_usecase.dart
│           └── add_vision_record_usecase.dart
│
├── presentation/              # 表现层
│   ├── providers/             # Riverpod Provider
│   │   ├── auth_provider.dart
│   │   ├── user_provider.dart
│   │   ├── exercise_provider.dart
│   │   ├── vision_provider.dart
│   │   └── reward_provider.dart
│   ├── screens/              # 页面
│   │   ├── splash/
│   │   │   └── splash_screen.dart
│   │   ├── auth/
│   │   │   ├── login_screen.dart
│   │   │   └── register_screen.dart
│   │   ├── child/            # 儿童端
│   │   │   ├── child_home_screen.dart
│   │   │   ├── exercise_screen.dart
│   │   │   ├── challenge_screen.dart
│   │   │   ├── achievement_screen.dart
│   │   │   └── child_profile_screen.dart
│   │   └── parent/           # 家长端
│   │       ├── parent_home_screen.dart
│   │       ├── training_plan_screen.dart
│   │       ├── vision_screen.dart
│   │       ├── growth_screen.dart
│   │       ├── challenge_management_screen.dart
│   │       └── parent_settings_screen.dart
│   └── widgets/              # 通用组件
│       ├── buttons/
│       │   └── primary_button.dart
│       ├── cards/
│       │   └── exercise_card.dart
│       ├── charts/
│       │   └── radar_chart.dart
│       └── dialogs/
│           └── confirm_dialog.dart
│
├── ai/                        # AI模块
│   ├── pose_detector.dart
│   ├── exercise_analyzer.dart
│   ├── pose_painter.dart
│   └── exercise_counters/
│       ├── jump_rope_counter.dart
│       ├── squat_counter.dart
│       └── jumping_jack_counter.dart
│
└── services/                  # 服务层
    ├── storage_service.dart
    ├── camera_service.dart
    ├── location_service.dart
    ├── notification_service.dart
    └── tts_service.dart
```

## 功能模块

### 儿童端功能

1. **首页** - 今日任务、快速开始入口
2. **运动** - AI训练、游戏化课程
3. **挑战** - 竞技模式、亲子对战
4. **成就** - 勋章墙、成长记录
5. **我的** - 角色装扮、设置

### 家长端功能

1. **首页** - 数据看板、今日概览
2. **运动** - 训练计划、课程管理
3. **视力** - 视力档案、用眼分析
4. **成长** - 发育数据、能力雷达图
5. **我的** - 亲子任务、安全管控

## 隐私保护

- 敏感数据（手机号、密码）本地加密存储
- 摄像头数据实时处理，不存储原始视频
- 儿童个人信息最小化收集
- 家长完全控制权

## 开发指南

### 环境要求

- Flutter SDK: >= 3.0.0
- Dart SDK: >= 3.0.0
- Android SDK: API 26+ (Android 8.0)
- iOS: 12.0+

### 运行项目

```bash
# 获取依赖
flutter pub get

# 运行调试版本
flutter run

# 构建Android APK
flutter build apk --release

# 构建Android App Bundle
flutter build appbundle --release
```

### 代码规范

- 遵循 Dart 风格指南
- 使用 flutter_lints 进行静态分析
- 单元测试覆盖率 >= 80%
- 注释覆盖率 >= 80%

## 许可证

本项目仅供学习交流使用。
