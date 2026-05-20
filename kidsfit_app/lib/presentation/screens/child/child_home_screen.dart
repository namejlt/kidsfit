import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../core/router/app_router.dart';
import '../../../domain/entities/exercise_record.dart';
import '../../../domain/entities/reward.dart';
import '../../providers/auth_provider.dart';
import '../../providers/exercise_provider.dart';
import '../../providers/reward_provider.dart';
import '../../widgets/exercise_card.dart';
import '../../widgets/stat_card.dart';

/// 儿童端首页
/// 展示今日任务、快速开始入口和成就概览
class ChildHomeScreen extends ConsumerStatefulWidget {
  const ChildHomeScreen({super.key});

  @override
  ConsumerState<ChildHomeScreen> createState() => _ChildHomeScreenState();
}

class _ChildHomeScreenState extends ConsumerState<ChildHomeScreen> {
  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    // 加载今日计划
    ref.read(exerciseProvider.notifier).loadTodayPlan();
    // 加载积分和勋章信息
    ref.read(rewardProvider.notifier).loadPointsBalance();
    ref.read(rewardProvider.notifier).loadMyBadges();
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authStateProvider);
    final exerciseState = ref.watch(exerciseProvider);
    final rewardState = ref.watch(rewardProvider);
    final todayPlan = exerciseState.todayPlan;
    final userBadges = rewardState.myBadges;
    final userPoints = rewardState.pointsBalance;

    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(AppDimensions.pagePadding),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 顶部问候
              _buildGreeting(authState.user?.nickname ?? '小勇士'),

              const SizedBox(height: AppDimensions.spacingLG),

              // 成就卡片
              _buildAchievementCard(userPoints, userBadges.length),

              const SizedBox(height: AppDimensions.spacingLG),

              // 今日任务
              _buildTodayTask(todayPlan),

              const SizedBox(height: AppDimensions.spacingLG),

              // 快速开始
              _buildQuickStart(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 推荐运动
              _buildRecommendedExercises(),
            ],
          ),
        ),
      ),
      bottomNavigationBar: _buildBottomNavBar(),
    );
  }

  /// 构建问候语
  Widget _buildGreeting(String nickname) {
    final hour = DateTime.now().hour;
    String greeting;
    if (hour < 12) {
      greeting = '早上好';
    } else if (hour < 18) {
      greeting = '下午好';
    } else {
      greeting = '晚上好';
    }

    return Row(
      children: [
        // 头像
        CircleAvatar(
          radius: 28,
          backgroundColor: AppColors.primary,
          child: Text(
            nickname.substring(0, 1),
            style: const TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
        ),
        const SizedBox(width: AppDimensions.spacingMD),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '$greeting，$nickname',
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            const Text(
              '今天也要加油运动哦 💪',
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ],
    );
  }

  /// 构建成就卡片
  Widget _buildAchievementCard(int points, int badgeCount) {
    return Container(
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        gradient: AppColors.primaryGradient,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        boxShadow: [
          BoxShadow(
            color: AppColors.primary.withOpacity(0.3),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '我的成就',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.white70,
                  ),
                ),
                const SizedBox(height: 8),
                Row(
                  children: [
                    const Icon(
                      Icons.star,
                      color: AppColors.accent,
                      size: 28,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '$points',
                      style: const TextStyle(
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                    const Text(
                      ' 积分',
                      style: TextStyle(
                        fontSize: 16,
                        color: Colors.white70,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          GestureDetector(
            onTap: () => AppNavigator.goToAchievement(context),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.2),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Row(
                children: [
                  Text(
                    '$badgeCount',
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  ),
                  const SizedBox(width: 4),
                  const Icon(
                    Icons.emoji_events,
                    color: AppColors.accent,
                    size: 20,
                  ),
                  const Icon(
                    Icons.arrow_forward_ios,
                    color: Colors.white70,
                    size: 14,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  /// 构建今日任务
  Widget _buildTodayTask(TrainingPlan? plan) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              AppStrings.todayTask,
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            if (plan != null)
              TextButton(
                onPressed: () => AppNavigator.goToTrainingPlan(context),
                child: const Text('查看全部'),
              ),
          ],
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        if (plan == null)
          // 无任务时显示提示
          Container(
            padding: const EdgeInsets.all(AppDimensions.cardPadding),
            decoration: BoxDecoration(
              color: AppColors.white,
              borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
              border: Border.all(color: AppColors.background),
            ),
            child: const Center(
              child: Column(
                children: [
                  Icon(
                    Icons.calendar_today,
                    size: 48,
                    color: AppColors.textSecondary,
                  ),
                  SizedBox(height: 8),
                  Text(
                    '今日暂无训练计划',
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
          )
        else
          // 有任务时显示进度
          Container(
            padding: const EdgeInsets.all(AppDimensions.cardPadding),
            decoration: BoxDecoration(
              color: AppColors.white,
              borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.05),
                  blurRadius: 10,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          plan.isCompleted ? '已完成' : '进行中',
                          style: TextStyle(
                            fontSize: 12,
                            color: plan.isCompleted
                                ? AppColors.success
                                : AppColors.primary,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${plan.totalDuration}分钟训练',
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                            color: AppColors.textPrimary,
                          ),
                        ),
                      ],
                    ),
                    CircularProgressIndicator(
                      value: plan.progress,
                      backgroundColor: AppColors.background,
                      valueColor: const AlwaysStoppedAnimation<Color>(
                        AppColors.primary,
                      ),
                      strokeWidth: 6,
                    ),
                  ],
                ),
                if (!plan.isCompleted) ...[
                  const SizedBox(height: AppDimensions.spacingMD),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () {
                        // 开始今日训练
                        if (plan.mainItems.isNotEmpty) {
                          AppNavigator.goToExercise(
                            context,
                            plan.mainItems.first.type.value,
                          );
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.primary,
                        foregroundColor: Colors.white,
                      ),
                      child: Text(AppStrings.goExercise),
                    ),
                  ),
                ],
              ],
            ),
          ),
      ],
    );
  }

  /// 构建快速开始区域
  Widget _buildQuickStart() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          AppStrings.quickStart,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        // 快速运动选项
        GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: 2,
          mainAxisSpacing: AppDimensions.spacingMD,
          crossAxisSpacing: AppDimensions.spacingMD,
          childAspectRatio: 1.2,
          children: ExerciseType.values.take(4).map((type) {
            return ExerciseCard(
              exerciseType: type,
              onTap: () => AppNavigator.goToExercise(context, type.value),
            );
          }).toList(),
        ),
      ],
    );
  }

  /// 构建推荐运动区域
  Widget _buildRecommendedExercises() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '更多运动',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        SizedBox(
          height: 120,
          child: ListView(
            scrollDirection: Axis.horizontal,
            children: ExerciseType.values.skip(4).map((type) {
              return Container(
                width: 100,
                margin: const EdgeInsets.only(right: AppDimensions.spacingMD),
                decoration: BoxDecoration(
                  color: AppColors.white,
                  borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 10,
                    ),
                  ],
                ),
                child: InkWell(
                  onTap: () => AppNavigator.goToExercise(context, type.value),
                  borderRadius: BorderRadius.circular(AppDimensions.radiusLG),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        _getExerciseIcon(type),
                        size: 40,
                        color: AppColors.primary,
                      ),
                      const SizedBox(height: 8),
                      Text(
                        type.displayName,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                          color: AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }).toList(),
          ),
        ),
      ],
    );
  }

  /// 获取运动图标
  IconData _getExerciseIcon(ExerciseType type) {
    switch (type) {
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

  /// 构建底部导航栏
  Widget _buildBottomNavBar() {
    return BottomNavigationBar(
      currentIndex: 0,
      onTap: (index) {
        switch (index) {
          case 0:
            // 已在首页
            break;
          case 1:
            AppNavigator.goToChallenge(context);
            break;
          case 2:
            AppNavigator.goToAchievement(context);
            break;
        }
      },
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.home),
          label: '首页',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.emoji_events),
          label: '挑战',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.star),
          label: '成就',
        ),
      ],
    );
  }
}
