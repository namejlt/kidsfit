import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';

import 'app.dart';
import 'core/constants/app_constants.dart';
import 'services/storage_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // 初始化Hive本地存储
  await Hive.initFlutter();

  // 打开Hive Boxes
  await StorageService.initialize();

  // 设置错误处理
  setupErrorHandling();

  runApp(
    const ProviderScope(
      child: KidsFitApp(),
    ),
  );
}

/// 设置全局错误处理
void setupErrorHandling() {
  // 生产环境错误处理
  FlutterError.onError = (FlutterErrorDetails details) {
    // 记录错误日志
    debugPrint('Flutter Error: ${details.exceptionAsString()}');
    // 可以发送到错误追踪服务
  };

  // 异步错误处理
  PlatformDispatcher.instance.onError = (Object error, StackTrace stack) {
    debugPrint('Platform Error: $error');
    return true;
  };
}
