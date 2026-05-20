import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../core/router/app_router.dart';
import '../../../domain/entities/exercise_record.dart';
import '../../providers/exercise_provider.dart';

/// 训练计划页面
/// 展示和管理训练计划
class TrainingPlanScreen extends ConsumerStatefulWidget {
  const TrainingPlanScreen({super.key});

  @override
  ConsumerState<TrainingPlanScreen> createState() =>
      _TrainingPlanScreenState();
}

class _TrainingPlanScreenState extends ConsumerState<TrainingPlanScreen> {
  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    ref.read(exerciseProvider.notifier).loadTodayPlan();
  }

  @override
  Widget build(BuildContext context) {
    final exerciseState = ref.watch(exerciseProvider);
    final todayPlan = exerciseState.todayPlan;

    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.trainingPlan),
        actions: [
          IconButton(
            icon: const Icon(Icons.calendar_today),
            onPressed: () => _showCalendarView(),
          ),
          IconButton(
            icon: const Icon(Icons.tune),
            onPressed: () => _showPlanSettings(),
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
              // 日期选择器
              _buildDateSelector(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 计划概览
              if (todayPlan != null) _buildPlanOverview(todayPlan),

              const SizedBox(height: AppDimensions.spacingLG),

              // 训练内容
              if (todayPlan != null) ...[
                // 热身
                _buildPhaseSection(
                  title: '热身运动',
                  icon: Icons.directions_walk,
                  color: AppColors.accent,
                  items: todayPlan.warmupItems,
                ),

                const SizedBox(height: AppDimensions.spacingMD),

                // 主训练
                _buildPhaseSection(
                  title: '主训练',
                  icon: Icons.fitness_center,
                  color: AppColors.primary,
                  items: todayPlan.mainItems,
                ),

                const SizedBox(height: AppDimensions.spacingMD),

                // 拉伸
                _buildPhaseSection(
                  title: '拉伸放松',
                  icon: Icons.self_improvement,
                  color: AppColors.secondary,
                  items: todayPlan.cooldownItems,
                ),
              ] else
                _buildEmptyPlan(),
            ],
          ),
        ),
      ),
      floatingActionButton: todayPlan != null && !todayPlan.isCompleted
          ? FloatingActionButton.extended(
              onPressed: () {
                if (todayPlan.mainItems.isNotEmpty) {
                  AppNavigator.goToExercise(
                    context,
                    todayPlan.mainItems.first.type.value,
                  );
                }
              },
              icon: const Icon(Icons.play_arrow),
              label: const Text('开始训练'),
              backgroundColor: AppColors.primary,
            )
          : null,
    );
  }

  /// 构建日期选择器
  Widget _buildDateSelector() {
    final now = DateTime.now();
    final dates = List.generate(7, (index) {
      return now.add(Duration(days: index - 3));
    });

    return Container(
      padding: const EdgeInsets.symmetric(vertical: AppDimensions.spacingMD),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: dates.map((date) {
          final isSelected = date.day == now.day;
          final isToday = date.day == now.day;

          return GestureDetector(
            onTap: () {},
            child: Container(
              width: 44,
              height: 64,
              decoration: BoxDecoration(
                color: isSelected ? AppColors.primary : Colors.transparent,
                borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    ['一', '二', '三', '四', '五', '六', '日'][date.weekday - 1],
                    style: TextStyle(
                      fontSize: 12,
                      color: isSelected ? Colors.white : AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${date.day}',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: isSelected
                          ? Colors.white
                          : AppColors.textPrimary,
                    ),
                  ),
                  if (isToday)
                    Container(
                      width: 4,
                      height: 4,
                      margin: const EdgeInsets.only(top: 4),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? Colors.white
                            : AppColors.primary,
                        shape: BoxShape.circle,
                      ),
                    ),
                ],
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  /// 构建计划概览
  Widget _buildPlanOverview(TrainingPlan plan) {
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
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    plan.isCompleted ? '已完成今日训练' : '今日训练计划',
                    style: const TextStyle(
                      color: Colors.white70,
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${plan.totalDuration}分钟',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 28,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
              _buildOverviewStat(
                '${plan.mainItems.length}',
                '个动作',
              ),
            ],
          ),
          const SizedBox(height: AppDimensions.spacingMD),

          // 进度条
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: LinearProgressIndicator(
              value: plan.progress,
              backgroundColor: Colors.white.withOpacity(0.3),
              valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
              minHeight: 8,
            ),
          ),

          const SizedBox(height: AppDimensions.spacingSM),

          Text(
            '完成度 ${(plan.progress * 100).toInt()}%',
            style: const TextStyle(
              color: Colors.white70,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  /// 构建概览统计项
  Widget _buildOverviewStat(String value, String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.2),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Text(
            value,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            label,
            style: const TextStyle(
              color: Colors.white70,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  /// 构建阶段区域
  Widget _buildPhaseSection({
    required String title,
    required IconData icon,
    required Color color,
    required List<ExerciseItem> items,
  }) {
    if (items.isEmpty) {
      return const SizedBox.shrink();
    }

    return Container(
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
        children: [
          // 标题
          Container(
            padding: const EdgeInsets.all(AppDimensions.cardPadding),
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: const BorderRadius.vertical(
                top: Radius.circular(AppDimensions.cardRadius),
              ),
            ),
            child: Row(
              children: [
                Icon(icon, color: color, size: 24),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
              ],
            ),
          ),

          // 动作列表
          ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: items.length,
            separatorBuilder: (context, index) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final item = items[index];
              return _buildExerciseItem(item, color);
            },
          ),
        ],
      ),
    );
  }

  /// 构建运动项
  Widget _buildExerciseItem(ExerciseItem item, Color color) {
    return InkWell(
      onTap: () => AppNavigator.goToExercise(context, item.type.value),
      child: Padding(
        padding: const EdgeInsets.all(AppDimensions.cardPadding),
        child: Row(
          children: [
            // 序号
            Container(
              width: 32,
              height: 32,
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Center(
                child: Text(
                  '${item.order}',
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),

            // 动作信息
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    item.name,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  if (item.tips != null)
                    Text(
                      item.tips!,
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                ],
              ),
            ),

            // 目标
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                if (item.targetCount != null)
                  Text(
                    '${item.targetCount}次',
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                if (item.durationSec != null)
                  Text(
                    '${item.durationSec! ~/ 60}分钟',
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textSecondary,
                    ),
                  ),
              ],
            ),

            const SizedBox(width: 8),

            // 难度
            Row(
              children: List.generate(5, (index) {
                return Container(
                  width: 6,
                  height: 6,
                  margin: const EdgeInsets.only(left: 2),
                  decoration: BoxDecoration(
                    color: index < item.difficulty
                        ? color
                        : AppColors.background,
                    shape: BoxShape.circle,
                  ),
                );
              }),
            ),

            const SizedBox(width: 8),

            const Icon(
              Icons.play_circle_outline,
              color: AppColors.textSecondary,
            ),
          ],
        ),
      ),
    );
  }

  /// 构建空计划
  Widget _buildEmptyPlan() {
    return Container(
      padding: const EdgeInsets.all(AppDimensions.cardPadding * 2),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        border: Border.all(color: AppColors.background),
      ),
      child: Center(
        child: Column(
          children: [
            Icon(
              Icons.event_note,
              size: 64,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            const Text(
              '今日暂无训练计划',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '系统将根据孩子的体能情况\n自动生成训练计划',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: AppColors.textSecondary,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 24),
            ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add),
              label: const Text('手动添加计划'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// 显示日历视图
  void _showCalendarView() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return Container(
          padding: const EdgeInsets.all(AppDimensions.pagePadding),
          height: 400,
          child: Column(
            children: [
              const Text(
                '训练日历',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 16),
              Expanded(
                child: GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 7,
                    childAspectRatio: 1,
                  ),
                  itemCount: 35,
                  itemBuilder: (context, index) {
                    return Center(
                      child: Container(
                        width: 36,
                        height: 36,
                        decoration: BoxDecoration(
                          color: index == 10
                              ? AppColors.primary
                              : Colors.transparent,
                          shape: BoxShape.circle,
                        ),
                        child: Center(
                          child: Text(
                            '${(index % 30) + 1}',
                            style: TextStyle(
                              color: index == 10
                                  ? Colors.white
                                  : AppColors.textPrimary,
                            ),
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  /// 显示计划设置
  void _showPlanSettings() {
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
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                '训练设置',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: AppDimensions.spacingLG),

              // 每日时长
              ListTile(
                leading: const Icon(Icons.timer),
                title: const Text('每日训练时长'),
                subtitle: const Text('当前: 20分钟'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),

              // 重点训练
              ListTile(
                leading: const Icon(Icons.fitness_center),
                title: const Text('重点训练部位'),
                subtitle: const Text('力量、耐力'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),

              // 运动偏好
              ListTile(
                leading: const Icon(Icons.sports),
                title: const Text('运动偏好'),
                subtitle: const Text('跳绳、跑步'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),

              const SizedBox(height: AppDimensions.spacingLG),

              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {},
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.primary,
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('保存设置'),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
