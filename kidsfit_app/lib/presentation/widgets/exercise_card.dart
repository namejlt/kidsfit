import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_dimensions.dart';
import '../../domain/entities/exercise_record.dart';

/// 运动卡片组件
/// 展示单个运动项目的快速入口卡片
class ExerciseCard extends StatelessWidget {
  /// 运动类型
  final ExerciseType exerciseType;

  /// 点击回调
  final VoidCallback? onTap;

  const ExerciseCard({
    super.key,
    required this.exerciseType,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
        child: Padding(
          padding: const EdgeInsets.all(AppDimensions.cardPadding),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              // 运动图标
              Container(
                width: 56,
                height: 56,
                decoration: BoxDecoration(
                  color: AppColors.primary.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(14),
                ),
                child: Icon(
                  _getExerciseIcon(),
                  size: 32,
                  color: AppColors.primary,
                ),
              ),
              const SizedBox(height: AppDimensions.spacingSM),

              // 运动名称
              Text(
                exerciseType.displayName,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
                textAlign: TextAlign.center,
              ),

              const SizedBox(height: 4),

              // 难度指示
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(5, (index) {
                  return Container(
                    width: 6,
                    height: 6,
                    margin: const EdgeInsets.symmetric(horizontal: 1),
                    decoration: BoxDecoration(
                      color: index < 2
                          ? AppColors.primary
                          : AppColors.background,
                      shape: BoxShape.circle,
                    ),
                  );
                }),
              ),
            ],
          ),
        ),
      ),
    );
  }

  /// 获取运动图标
  IconData _getExerciseIcon() {
    switch (exerciseType) {
      case ExerciseType.jumpRope:
        return Icons.hardware;
      case ExerciseType.jumpingJack:
        return Icons.accessibility_new;
      case ExerciseType.squat:
        return Icons.fitness_center;
      case ExerciseType.sitUp:
        return Icons.airline_seat_flat;
      case ExerciseType.highKnee:
        return Icons.directions_run;
      case ExerciseType.pushUp:
        return Icons.sports_gymnastics;
    }
  }
}

/// 统计卡片组件
/// 展示单个统计指标的卡片
class StatCard extends StatelessWidget {
  /// 标题
  final String title;

  /// 数值
  final String value;

  /// 单位
  final String? unit;

  /// 图标
  final IconData icon;

  /// 图标颜色
  final Color? iconColor;

  /// 背景色
  final Color? backgroundColor;

  /// 点击回调
  final VoidCallback? onTap;

  const StatCard({
    super.key,
    required this.title,
    required this.value,
    this.unit,
    required this.icon,
    this.iconColor,
    this.backgroundColor,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: backgroundColor ?? AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 图标
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: (iconColor ?? AppColors.primary).withOpacity(0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                icon,
                size: 24,
                color: iconColor ?? AppColors.primary,
              ),
            ),
            const SizedBox(height: AppDimensions.spacingMD),

            // 数值
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  value,
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                if (unit != null) ...[
                  const SizedBox(width: 4),
                  Padding(
                    padding: const EdgeInsets.only(bottom: 2),
                    child: Text(
                      unit!,
                      style: const TextStyle(
                        fontSize: 14,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ),
                ],
              ],
            ),

            const SizedBox(height: 4),

            // 标题
            Text(
              title,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// 进度卡片组件
class ProgressCard extends StatelessWidget {
  /// 标题
  final String title;

  /// 当前进度
  final double progress;

  /// 当前值
  final String currentValue;

  /// 目标值
  final String targetValue;

  /// 进度颜色
  final Color? progressColor;

  /// 右侧额外内容
  final Widget? trailing;

  const ProgressCard({
    super.key,
    required this.title,
    required this.progress,
    required this.currentValue,
    required this.targetValue,
    this.progressColor,
    this.trailing,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                title,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              if (trailing != null) trailing!,
            ],
          ),
          const SizedBox(height: AppDimensions.spacingMD),

          // 进度条
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: LinearProgressIndicator(
              value: progress.clamp(0, 1),
              backgroundColor: AppColors.background,
              valueColor: AlwaysStoppedAnimation<Color>(
                progressColor ?? AppColors.primary,
              ),
              minHeight: 8,
            ),
          ),

          const SizedBox(height: AppDimensions.spacingSM),

          // 数值
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                currentValue,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: progressColor ?? AppColors.primary,
                ),
              ),
              Text(
                '目标: $targetValue',
                style: const TextStyle(
                  fontSize: 12,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
