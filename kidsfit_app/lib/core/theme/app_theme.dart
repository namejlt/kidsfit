import 'package:flutter/material.dart';
import '../constants/app_colors.dart';
import '../constants/app_dimensions.dart';

/// KidsFit应用主题配置
/// 遵循UI/UX设计规范中的色彩方案和组件规范
class AppTheme {
  AppTheme._();

  /// 明亮主题
  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,

      // 色彩配置
      colorScheme: const ColorScheme.light(
        primary: AppColors.primary,
        secondary: AppColors.secondary,
        tertiary: AppColors.accent,
        surface: AppColors.white,
        error: AppColors.error,
        onPrimary: AppColors.white,
        onSecondary: AppColors.white,
        onSurface: AppColors.textPrimary,
        onError: AppColors.white,
      ),

      // 脚手架背景色
      scaffoldBackgroundColor: AppColors.background,

      // AppBar主题
      appBarTheme: const AppBarTheme(
        elevation: 0,
        centerTitle: true,
        backgroundColor: AppColors.white,
        foregroundColor: AppColors.textPrimary,
        titleTextStyle: TextStyle(
          color: AppColors.textPrimary,
          fontSize: AppDimensions.fontSizeTitle,
          fontWeight: FontWeight.w600,
        ),
        iconTheme: IconThemeData(
          color: AppColors.textPrimary,
          size: AppDimensions.iconSizeMD,
        ),
      ),

      // 按钮主题
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          elevation: 0,
          minimumSize: const Size.fromHeight(AppDimensions.buttonHeight),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppDimensions.buttonRadius),
          ),
          textStyle: const TextStyle(
            fontSize: AppDimensions.fontSizeNormal,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),

      // 文本按钮主题
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          minimumSize: const Size.fromHeight(AppDimensions.buttonHeightSmall),
          textStyle: const TextStyle(
            fontSize: AppDimensions.fontSizeNormal,
            fontWeight: FontWeight.w500,
          ),
        ),
      ),

      // 轮廓按钮主题
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          minimumSize: const Size.fromHeight(AppDimensions.buttonHeightSecondary),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
          ),
          side: const BorderSide(
            color: AppColors.primary,
            width: AppDimensions.borderMedium,
          ),
          textStyle: const TextStyle(
            fontSize: AppDimensions.fontSizeNormal,
            fontWeight: FontWeight.w500,
          ),
        ),
      ),

      // 输入框主题
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.white,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppDimensions.spacingMD,
          vertical: AppDimensions.spacingMD,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppDimensions.inputRadius),
          borderSide: const BorderSide(
            color: AppColors.background,
            width: AppDimensions.borderThin,
          ),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppDimensions.inputRadius),
          borderSide: const BorderSide(
            color: AppColors.background,
            width: AppDimensions.borderThin,
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppDimensions.inputRadius),
          borderSide: const BorderSide(
            color: AppColors.primary,
            width: AppDimensions.borderMedium,
          ),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppDimensions.inputRadius),
          borderSide: const BorderSide(
            color: AppColors.error,
            width: AppDimensions.borderThin,
          ),
        ),
        labelStyle: const TextStyle(
          color: AppColors.textSecondary,
          fontSize: AppDimensions.inputLabelSize,
        ),
        hintStyle: const TextStyle(
          color: AppColors.textSecondary,
          fontSize: AppDimensions.inputTextSize,
        ),
      ),

      // 卡片主题
      cardTheme: CardTheme(
        elevation: AppDimensions.cardElevation,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        ),
        color: AppColors.white,
        margin: EdgeInsets.zero,
      ),

      // 底部导航栏主题
      bottomNavigationBarTheme: const BottomNavigationBarThemeData(
        type: BottomNavigationBarType.fixed,
        height: AppDimensions.bottomNavHeight,
        backgroundColor: AppColors.white,
        selectedItemColor: AppColors.primary,
        unselectedItemColor: AppColors.textSecondary,
        selectedLabelStyle: TextStyle(
          fontSize: AppDimensions.bottomNavLabelSize,
          fontWeight: FontWeight.w600,
        ),
        unselectedLabelStyle: TextStyle(
          fontSize: AppDimensions.bottomNavLabelSize,
          fontWeight: FontWeight.w400,
        ),
      ),

      // 对话框主题
      dialogTheme: DialogTheme(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
        ),
        backgroundColor: AppColors.white,
        titleTextStyle: const TextStyle(
          color: AppColors.textPrimary,
          fontSize: AppDimensions.fontSizeTitle,
          fontWeight: FontWeight.w600,
        ),
      ),

      // 浮动按钮主题
      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: AppColors.primary,
        foregroundColor: AppColors.white,
        elevation: 4,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
        ),
      ),

      // 进度指示器主题
      progressIndicatorTheme: const ProgressIndicatorThemeData(
        color: AppColors.primary,
        linearTrackColor: AppColors.background,
        circularTrackColor: AppColors.background,
      ),

      // 分隔线主题
      dividerTheme: const DividerThemeData(
        color: AppColors.background,
        thickness: AppDimensions.borderThin,
        space: AppDimensions.spacingMD,
      ),

      // 文字主题
      textTheme: const TextTheme(
        displayLarge: TextStyle(
          fontSize: AppDimensions.fontSizeDisplayXL,
          fontWeight: FontWeight.bold,
          color: AppColors.textPrimary,
        ),
        displayMedium: TextStyle(
          fontSize: AppDimensions.fontSizeDisplay,
          fontWeight: FontWeight.bold,
          color: AppColors.textPrimary,
        ),
        headlineLarge: TextStyle(
          fontSize: AppDimensions.fontSizeHeadline,
          fontWeight: FontWeight.w600,
          color: AppColors.textPrimary,
        ),
        headlineMedium: TextStyle(
          fontSize: AppDimensions.fontSizeTitle,
          fontWeight: FontWeight.w600,
          color: AppColors.textPrimary,
        ),
        titleLarge: TextStyle(
          fontSize: AppDimensions.fontSizeNormal,
          fontWeight: FontWeight.w600,
          color: AppColors.textPrimary,
        ),
        titleMedium: TextStyle(
          fontSize: AppDimensions.fontSizeNormal,
          fontWeight: FontWeight.w500,
          color: AppColors.textPrimary,
        ),
        bodyLarge: TextStyle(
          fontSize: AppDimensions.fontSizeNormal,
          color: AppColors.textPrimary,
        ),
        bodyMedium: TextStyle(
          fontSize: AppDimensions.fontSizeNormal,
          color: AppColors.textSecondary,
        ),
        labelLarge: TextStyle(
          fontSize: AppDimensions.fontSizeNormal,
          fontWeight: FontWeight.w500,
          color: AppColors.textPrimary,
        ),
      ),
    );
  }

  /// 暗黑主题
  static ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: ColorScheme.dark(
        primary: AppColors.primary,
        secondary: AppColors.secondary,
        tertiary: AppColors.accent,
        surface: const Color(0xFF1E1E1E),
        error: AppColors.error,
        onPrimary: AppColors.white,
        onSecondary: AppColors.white,
        onSurface: Colors.white,
        onError: AppColors.white,
      ),
      scaffoldBackgroundColor: const Color(0xFF121212),
    );
  }

  /// 儿童端主题（明亮活泼风格）
  static ThemeData get childTheme {
    return lightTheme.copyWith(
      colorScheme: lightTheme.colorScheme.copyWith(
        primary: AppColors.childPrimary,
        secondary: AppColors.childSecondary,
      ),
    );
  }

  /// 家长端主题（简洁专业风格）
  static ThemeData get parentTheme {
    return lightTheme.copyWith(
      colorScheme: lightTheme.colorScheme.copyWith(
        primary: AppColors.parentPrimary,
        secondary: AppColors.parentSecondary,
      ),
    );
  }
}
