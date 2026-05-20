import 'package:flutter/material.dart';

/// KidsFit应用颜色配置
/// 遵循UI/UX设计规范中的色彩方案
class AppColors {
  AppColors._();

  // ==================== 主色调 ====================

  /// 活力橙 - 主按钮、重点强调
  static const Color primary = Color(0xFFFF6B35);

  /// 天空蓝 - 次要按钮、背景
  static const Color secondary = Color(0xFF4ECDC4);

  /// 阳光黄 - 高亮、警告、成就
  static const Color accent = Color(0xFFFFE66D);

  /// 健康绿 - 成功、完成、户外
  static const Color success = Color(0xFF96CEB4);

  /// 温暖红 - 错误、重要提醒
  static const Color error = Color(0xFFFF6B6B);

  // ==================== 中性色 ====================

  /// 深灰 - 文字
  static const Color textPrimary = Color(0xFF2C3E50);

  /// 中灰 - 次要文字
  static const Color textSecondary = Color(0xFF7F8C8D);

  /// 浅灰 - 提示文字
  static const Color textHint = Color(0xFFBDC3C7);

  /// 浅灰 - 背景、分隔线
  static const Color background = Color(0xFFECF0F1);

  /// 白色 - 卡片背景
  static const Color white = Color(0xFFFFFFFF);

  /// 黑色
  static const Color black = Color(0xFF000000);

  // ==================== 功能色 ====================

  /// 儿童端主题色 - 明亮活泼
  static const Color childPrimary = Color(0xFFFF6B35);
  static const Color childSecondary = Color(0xFFFFE66D);

  /// 家长端主题色 - 简洁专业
  static const Color parentPrimary = Color(0xFF4ECDC4);
  static const Color parentSecondary = Color(0xFF2C3E50);

  // ==================== 渐变色 ====================

  /// 主按钮渐变
  static const LinearGradient primaryGradient = LinearGradient(
    colors: [Color(0xFFFF6B35), Color(0xFFFF8E53)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  /// 成就渐变
  static const LinearGradient achievementGradient = LinearGradient(
    colors: [Color(0xFFFFE66D), Color(0xFFFF6B35)],
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  /// 户外运动渐变
  static const LinearGradient outdoorGradient = LinearGradient(
    colors: [Color(0xFF96CEB4), Color(0xFF4ECDC4)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  // ==================== 年龄适配色彩 ====================

  /// 3-6岁 - 更多暖色调
  static const Color ageGroup36Primary = Color(0xFFFF6B35);
  static const Color ageGroup36Secondary = Color(0xFFFFE66D);

  /// 7-9岁 - 平衡色调
  static const Color ageGroup79Primary = Color(0xFF4ECDC4);
  static const Color ageGroup79Secondary = Color(0xFFFFE66D);

  /// 10-12岁 - 完整色板
  static const Color ageGroup1012Primary = Color(0xFF2C3E50);
  static const Color ageGroup1012Secondary = Color(0xFF4ECDC4);

  // ==================== 图表颜色 ====================

  /// 雷达图颜色
  static const List<Color> radarChartColors = [
    Color(0xFFFF6B35), // 耐力
    Color(0xFF4ECDC4), // 灵敏
    Color(0xFFFFE66D), // 力量
    Color(0xFF96CEB4), // 速度
    Color(0xFFFF6B6B), // 协调
    Color(0xFF7F8C8D), // 柔韧
  ];

  // ==================== 状态颜色 ====================

  /// 在线/活跃
  static const Color online = Color(0xFF96CEB4);

  /// 离线/未激活
  static const Color offline = Color(0xFF7F8C8D);

  /// 警告
  static const Color warning = Color(0xFFFFE66D);

  /// 视力状态 - 良好
  static const Color visionGood = Color(0xFF96CEB4);

  /// 视力状态 - 中等
  static const Color visionMedium = Color(0xFFFFE66D);

  /// 视力状态 - 关注
  static const Color visionConcern = Color(0xFFFF6B6B);
}
