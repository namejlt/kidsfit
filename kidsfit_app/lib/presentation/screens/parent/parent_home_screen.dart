import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../core/router/app_router.dart';
import '../../../domain/entities/exercise_record.dart';
import '../../../domain/entities/vision_record.dart';
import '../../providers/auth_provider.dart';
import '../../providers/exercise_provider.dart';
import '../../providers/vision_provider.dart';
import '../../widgets/stat_card.dart';

/// 家长端首页
/// 展示数据看板、今日概览和快速入口
class ParentHomeScreen extends ConsumerStatefulWidget {
  const ParentHomeScreen({super.key});

  @override
  ConsumerState<ParentHomeScreen> createState() => _ParentHomeScreenState();
}

class _ParentHomeScreenState extends ConsumerState<ParentHomeScreen> {
  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    // 加载儿童列表
    final children = await ref.read(authStateProvider.notifier).getChildren();

    // 加载今日训练计划
    ref.read(exerciseProvider.notifier).loadTodayPlan();

    // 加载户外活动数据
    ref.read(visionProvider.notifier).loadTodayOutdoor();
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authStateProvider);
    final exerciseState = ref.watch(exerciseProvider);
    final visionState = ref.watch(visionProvider);
    final todayPlan = exerciseState.todayPlan;
    final todayOutdoor = visionState.todayOutdoor;
    final children = ref.watch(childrenListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.dataDashboard),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () {},
          ),
          IconButton(
            icon: const Icon(Icons.settings_outlined),
            onPressed: () => _showSettingsMenu(),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(AppDimensions.pagePadding),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 儿童切换器
              if (children.isNotEmpty) _buildChildSwitcher(children),

              const SizedBox(height: AppDimensions.spacingLG),

              // 今日概览
              _buildTodayOverview(todayPlan, todayOutdoor),

              const SizedBox(height: AppDimensions.spacingLG),

              // 快捷入口
              _buildQuickActions(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 周运动趋势
              _buildWeeklyTrend(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 能力雷达图
              _buildAbilityRadar(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 最近运动记录
              _buildRecentExercises(),
            ],
          ),
        ),
      ),
      bottomNavigationBar: _buildBottomNavBar(),
    );
  }

  /// 构建儿童切换器
  Widget _buildChildSwitcher(List children) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppDimensions.spacingMD,
        vertical: AppDimensions.spacingSM,
      ),
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
      child: Row(
        children: [
          const Icon(
            Icons.child_care,
            color: AppColors.parentPrimary,
          ),
          const SizedBox(width: AppDimensions.spacingSM),
          Expanded(
            child: DropdownButtonHideUnderline(
              child: DropdownButton<String>(
                value: children.isNotEmpty ? children.first.id : null,
                isExpanded: true,
                hint: const Text('选择儿童'),
                items: children.map((child) {
                  return DropdownMenuItem<String>(
                    value: child.id,
                    child: Text('${child.nickname} (${child.age}岁)'),
                  );
                }).toList(),
                onChanged: (value) {
                  // 切换儿童
                },
              ),
            ),
          ),
        ],
      ),
    );
  }

  /// 构建今日概览卡片
  Widget _buildTodayOverview(TrainingPlan? plan, OutdoorActivity? outdoor) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          AppStrings.todayOverview,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        Row(
          children: [
            Expanded(
              child: StatCard(
                title: '训练时长',
                value: plan?.totalDuration.toString() ?? '0',
                unit: '分钟',
                icon: Icons.timer,
                iconColor: AppColors.primary,
              ),
            ),
            const SizedBox(width: AppDimensions.spacingMD),
            Expanded(
              child: StatCard(
                title: '户外时间',
                value: outdoor?.durationMin.toString() ?? '0',
                unit: '分钟',
                icon: Icons.wb_sunny,
                iconColor: AppColors.success,
              ),
            ),
          ],
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        // 户外目标进度
        ProgressCard(
          title: '户外运动目标',
          progress: (outdoor?.targetProgress ?? 0) / 100,
          currentValue: outdoor?.formattedDuration ?? '0分钟',
          targetValue: '120分钟',
          progressColor: AppColors.success,
        ),
      ],
    );
  }

  /// 构建快捷入口
  Widget _buildQuickActions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '快捷入口',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        Row(
          children: [
            Expanded(
              child: _buildActionCard(
                icon: Icons.fitness_center,
                title: '训练计划',
                color: AppColors.primary,
                onTap: () => AppNavigator.goToTrainingPlan(context),
              ),
            ),
            const SizedBox(width: AppDimensions.spacingMD),
            Expanded(
              child: _buildActionCard(
                icon: Icons.visibility,
                title: '视力档案',
                color: AppColors.secondary,
                onTap: () => AppNavigator.goToVision(context),
              ),
            ),
          ],
        ),
        const SizedBox(height: AppDimensions.spacingMD),
        Row(
          children: [
            Expanded(
              child: _buildActionCard(
                icon: Icons.show_chart,
                title: '成长记录',
                color: AppColors.accent,
                onTap: () {},
              ),
            ),
            const SizedBox(width: AppDimensions.spacingMD),
            Expanded(
              child: _buildActionCard(
                icon: Icons.family_restroom,
                title: '亲子任务',
                color: AppColors.success,
                onTap: () {},
              ),
            ),
          ],
        ),
      ],
    );
  }

  /// 构建操作卡片
  Widget _buildActionCard({
    required IconData icon,
    required String title,
    required Color color,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
      child: Container(
        padding: const EdgeInsets.all(AppDimensions.cardPadding),
        decoration: BoxDecoration(
          color: color.withOpacity(0.1),
          borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        ),
        child: Column(
          children: [
            Icon(
              icon,
              size: 32,
              color: color,
            ),
            const SizedBox(height: 8),
            Text(
              title,
              style: TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w600,
                color: color,
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// 构建周运动趋势
  Widget _buildWeeklyTrend() {
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
          const Text(
            '本周运动趋势',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppDimensions.spacingMD),

          // 简化的柱状图
          SizedBox(
            height: 150,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              crossAxisAlignment: CrossAxisAlignment.end,
              children: List.generate(7, (index) {
                final heights = [60, 80, 50, 100, 70, 90, 40];
                final days = ['一', '二', '三', '四', '五', '六', '日'];
                return Column(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    Container(
                      width: 30,
                      height: heights[index].toDouble(),
                      decoration: BoxDecoration(
                        color: index == 3
                            ? AppColors.primary
                            : AppColors.primary.withOpacity(0.5),
                        borderRadius: BorderRadius.circular(4),
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      days[index],
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                );
              }),
            ),
          ),
        ],
      ),
    );
  }

  /// 构建能力雷达图
  Widget _buildAbilityRadar() {
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
          const Text(
            '能力雷达图',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppDimensions.spacingMD),

          // 雷达图占位
          SizedBox(
            height: 200,
            child: Center(
              child: CustomPaint(
                size: const Size(180, 180),
                painter: _RadarChartPainter(),
              ),
            ),
          ),

          const SizedBox(height: AppDimensions.spacingMD),

          // 能力标签
          Wrap(
            spacing: 16,
            runSpacing: 8,
            children: [
              _buildAbilityChip('耐力', AppColors.radarChartColors[0]),
              _buildAbilityChip('灵敏', AppColors.radarChartColors[1]),
              _buildAbilityChip('力量', AppColors.radarChartColors[2]),
              _buildAbilityChip('速度', AppColors.radarChartColors[3]),
              _buildAbilityChip('协调', AppColors.radarChartColors[4]),
              _buildAbilityChip('柔韧', AppColors.radarChartColors[5]),
            ],
          ),
        ],
      ),
    );
  }

  /// 构建能力标签
  Widget _buildAbilityChip(String label, Color color) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(
            color: color,
            shape: BoxShape.circle,
          ),
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  /// 构建最近运动记录
  Widget _buildRecentExercises() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '最近运动',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        // 示例数据
        _buildExerciseItem(
          type: '跳绳',
          time: '今天 10:30',
          count: 50,
          score: 85,
        ),
        _buildExerciseItem(
          type: '深蹲',
          time: '今天 09:15',
          count: 20,
          score: 78,
        ),
        _buildExerciseItem(
          type: '开合跳',
          time: '昨天 16:00',
          count: 30,
          score: 90,
        ),
      ],
    );
  }

  /// 构建运动记录项
  Widget _buildExerciseItem({
    required String type,
    required String time,
    required int count,
    required int score,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacingMD),
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Icon(
              Icons.fitness_center,
              color: AppColors.primary,
            ),
          ),
          const SizedBox(width: AppDimensions.spacingMD),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  type,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                Text(
                  time,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '$count次',
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: AppColors.success.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '$score分',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.success,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 显示设置菜单
  void _showSettingsMenu() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return Container(
          padding: const EdgeInsets.all(AppDimensions.pagePadding),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ListTile(
                leading: const Icon(Icons.person),
                title: const Text('账号设置'),
                onTap: () {},
              ),
              ListTile(
                leading: const Icon(Icons.child_care),
                title: const Text('儿童管理'),
                onTap: () {},
              ),
              ListTile(
                leading: const Icon(Icons.lock),
                title: const Text('隐私设置'),
                onTap: () {},
              ),
              ListTile(
                leading: const Icon(Icons.logout, color: AppColors.error),
                title: const Text(
                  '退出登录',
                  style: TextStyle(color: AppColors.error),
                ),
                onTap: () async {
                  await ref.read(authStateProvider.notifier).logout();
                  if (mounted) {
                    AppNavigator.goToLogin(context);
                  }
                },
              ),
            ],
          ),
        );
      },
    );
  }

  /// 构建底部导航栏
  Widget _buildBottomNavBar() {
    return BottomNavigationBar(
      currentIndex: 0,
      onTap: (index) {},
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.home),
          label: '首页',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.visibility),
          label: '视力',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.person),
          label: '我的',
        ),
      ],
    );
  }
}

/// 雷达图绘制器
class _RadarChartPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width / 2 - 20;

    // 绘制背景六边形
    final bgPaint = Paint()
      ..color = AppColors.background
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

    for (var i = 1; i <= 5; i++) {
      final r = radius * i / 5;
      _drawPolygon(canvas, center, r, 6, bgPaint);
    }

    // 绘制轴线
    for (var i = 0; i < 6; i++) {
      final angle = -90.0 + i * 60.0;
      final radians = angle * 3.14159 / 180;
      final x = center.dx + radius * cos(radians);
      final y = center.dy + radius * sin(radians);
      canvas.drawLine(center, Offset(x, y), bgPaint);
    }

    // 绘制数据区域
    final dataPaint = Paint()
      ..color = AppColors.primary.withOpacity(0.3)
      ..style = PaintingStyle.fill;

    final path = Path();
    for (var i = 0; i < 6; i++) {
      final angle = -90.0 + i * 60.0;
      final radians = angle * 3.14159 / 180;
      // 示例数据
      final value = [0.8, 0.6, 0.7, 0.5, 0.9, 0.6][i];
      final x = center.dx + radius * value * cos(radians);
      final y = center.dy + radius * value * sin(radians);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    path.close();
    canvas.drawPath(path, dataPaint);

    // 绘制数据轮廓
    final linePaint = Paint()
      ..color = AppColors.primary
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;
    canvas.drawPath(path, linePaint);
  }

  void _drawPolygon(Canvas canvas, Offset center, double radius, int sides, Paint paint) {
    final path = Path();
    for (var i = 0; i < sides; i++) {
      final angle = -90.0 + i * (360 / sides);
      final radians = angle * 3.14159 / 180;
      final x = center.dx + radius * cos(radians);
      final y = center.dy + radius * sin(radians);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    path.close();
    canvas.drawPath(path, paint);
  }

  double cos(double radians) => math.cos(radians);
  double sin(double radians) => math.sin(radians);

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
