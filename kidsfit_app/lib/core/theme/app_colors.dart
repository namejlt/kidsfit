import 'package:flutter/material.dart';

/// KidsFit应用颜色配置
class KidsFitColors {
  KidsFitColors._();

  /// 主色调 - 活力橙
  static const Color primary = Color(0xFFFF6B35);

  /// 次要色 - 天空蓝
  static const Color secondary = Color(0xFF4ECDC4);

  /// 强调色 - 阳光黄
  static const Color accent = Color(0xFFFFE66D);

  /// 成功色 - 森林绿
  static const Color success = Color(0xFF2ECC71);

  /// 警告色 - 珊瑚橙
  static const Color warning = Color(0xFFF39C12);

  /// 错误色 - 玫红色
  static const Color error = Color(0xFFE74C3C);

  /// 文本颜色
  static const Color textPrimary = Color(0xFF2C3E50);
  static const Color textSecondary = Color(0xFF7F8C8D);
  static const Color textHint = Color(0xFFBDC3C7);

  /// 背景色
  static const Color background = Color(0xFFF5F6FA);
  static const Color surface = Color(0xFFFFFFFF);

  /// 视力状态颜色
  static const Color visionGood = Color(0xFF2ECC71);
  static const Color visionMedium = Color(0xFFF39C12);
  static const Color visionConcern = Color(0xFFE74C3C);

  /// 评分等级颜色
  static const Color gradeS = Color(0xFFFFD700);
  static const Color gradeA = Color(0xFFC0C0C0);
  static const Color gradeB = Color(0xFFCD7F32);
  static const Color gradeC = Color(0xFF808080);
  static const Color gradeD = Color(0xFF606060);

  /// 儿童模式渐变
  static const LinearGradient childGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [primary, Color(0xFFFF8C42)],
  );

  /// 家长模式渐变
  static const LinearGradient parentGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [secondary, Color(0xFF45B7D1)],
  );

  /// 户外运动渐变
  static const LinearGradient outdoorGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [Color(0xFF56AB2F), Color(0xFFA8E063)],
  );

  /// 夜间模式颜色
  static const Color darkBackground = Color(0xFF1A1A2E);
  static const Color darkSurface = Color(0xFF16213E);
  static const Color darkText = Color(0xFFE8E8E8);
}

/// 年龄组颜色
class AgeGroupColors {
  AgeGroupColors._();

  /// 3-6岁组
  static const Color age3to6 = Color(0xFFFF6B9D);

  /// 7-9岁组
  static const Color age7to9 = Color(0xFF4ECDC4);

  /// 10-12岁组
  static const Color age10to12 = Color(0xFF45B7D1);
}
