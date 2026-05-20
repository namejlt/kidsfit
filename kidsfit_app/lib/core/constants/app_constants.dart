import 'dart:io';

/// 应用常量配置
/// 包含API地址、存储键名、超时时间等全局常量
class AppConstants {
  AppConstants._();

  // ==================== API配置 ====================

  /// Android模拟器访问本机的API基础URL
  static const String apiBaseUrl = 'http://10.0.2.2:8001';

  /// iOS模拟器访问本机的API基础URL
  static const String apiBaseUrlIOS = 'http://localhost:8001';

  /// 根据当前平台自动选择API基础URL
  static String get currentApiBaseUrl {
    if (Platform.isIOS) {
      return apiBaseUrlIOS;
    }
    return apiBaseUrl;
  }

  // ==================== Token存储键名 ====================

  /// 访问Token存储键
  static const String tokenKey = 'access_token';

  /// 刷新Token存储键
  static const String refreshTokenKey = 'refresh_token';

  /// Token过期时间存储键
  static const String tokenExpiryKey = 'token_expiry';

  // ==================== 网络配置 ====================

  /// 连接超时时间（毫秒）
  static const int connectionTimeout = 15000;

  /// 接收超时时间（毫秒）
  static const int receiveTimeout = 15000;

  /// 最大重试次数
  static const int maxRetryCount = 3;

  // ==================== API路径 ====================

  /// 认证相关路径
  static const String authRegister = '/api/v1/auth/register';
  static const String authLogin = '/api/v1/auth/login';
  static const String authRefresh = '/api/v1/auth/refresh';
  static const String authLogout = '/api/v1/auth/logout';

  /// 用户相关路径
  static const String userMe = '/api/v1/users/me';
  static const String userUpdate = '/api/v1/users/me';
  static const String userChildren = '/api/v1/users/children';
  static const String userAddChild = '/api/v1/users/children';
  static const String parentSettings = '/api/v1/users/parent-settings';

  /// 训练相关路径
  static const String exerciseRecords = '/api/v1/training/exercises';
  static const String personalBest = '/api/v1/training/personal-best';
  static const String todayPlan = '/api/v1/training/plan/today';
  static const String completePlan = '/api/v1/training/plan';
  static const String fitnessAssessment = '/api/v1/training/assessments';
  static const String latestAssessment =
      '/api/v1/training/assessments/latest';
  static const String weeklyStats = '/api/v1/training/stats/weekly';
  static const String monthlyStats = '/api/v1/training/stats/monthly';

  /// 视力相关路径
  static const String visionRecords = '/api/v1/vision/records';
  static const String visionTrend = '/api/v1/vision/trend';
  static const String todayOutdoor = '/api/v1/vision/outdoor/today';
  static const String syncOutdoor = '/api/v1/vision/outdoor/sync';
  static const String eyeReminders = '/api/v1/vision/reminders';

  /// 激励相关路径
  static const String badges = '/api/v1/rewards/badges';
  static const String myBadges = '/api/v1/rewards/badges/mine';
  static const String points = '/api/v1/rewards/points';
  static const String pointsBalance = '/api/v1/rewards/points/balance';
  static const String challenges = '/api/v1/rewards/challenges';
  static const String familyLeaderboard =
      '/api/v1/rewards/leaderboard/family';
  static const String globalLeaderboard =
      '/api/v1/rewards/leaderboard/global';
}
