import 'package:flutter/material.dart' hide Badge;
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../domain/entities/reward.dart';
import '../../providers/reward_provider.dart';

/// 成就页面
/// 展示勋章墙和积分信息
class AchievementScreen extends ConsumerStatefulWidget {
  const AchievementScreen({super.key});

  @override
  ConsumerState<AchievementScreen> createState() => _AchievementScreenState();
}

class _AchievementScreenState extends ConsumerState<AchievementScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final rewardState = ref.watch(rewardProvider);
    final userBadges = rewardState.myBadges;

    return Scaffold(
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) {
          return [
            SliverAppBar(
              expandedHeight: 200,
              floating: false,
              pinned: true,
              flexibleSpace: FlexibleSpaceBar(
                background: Container(
                  decoration: const BoxDecoration(
                    gradient: AppColors.primaryGradient,
                  ),
                  child: SafeArea(
                    child: _buildHeader(rewardState),
                  ),
                ),
              ),
              bottom: TabBar(
                controller: _tabController,
                tabs: const [
                  Tab(text: '勋章墙'),
                  Tab(text: '积分商城'),
                ],
                labelColor: Colors.white,
                unselectedLabelColor: Colors.white70,
                indicatorColor: Colors.white,
              ),
            ),
          ];
        },
        body: TabBarView(
          controller: _tabController,
          children: [
            _buildBadgeWall(userBadges),
            _buildPointsMall(),
          ],
        ),
      ),
    );
  }

  /// 构建头部信息
  Widget _buildHeader(RewardState state) {
    return Padding(
      padding: const EdgeInsets.all(AppDimensions.pagePadding),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildStatItem(
                icon: Icons.star,
                value: '${state.pointsBalance}',
                label: '积分',
              ),
              _buildStatItem(
                icon: Icons.emoji_events,
                value: '${state.myBadges.length}',
                label: '勋章',
              ),
              _buildStatItem(
                icon: Icons.local_fire_department,
                value: '${state.challenges.length}',
                label: '挑战',
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// 构建统计项
  Widget _buildStatItem({
    required IconData icon,
    required String value,
    required String label,
  }) {
    return Column(
      children: [
        Icon(
          icon,
          color: Colors.white,
          size: 32,
        ),
        const SizedBox(height: 8),
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: Colors.white70,
            fontSize: 14,
          ),
        ),
      ],
    );
  }

  /// 构建勋章墙
  Widget _buildBadgeWall(List<Badge> badges) {
    if (badges.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.emoji_events,
              size: 80,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            const Text(
              '还没有获得勋章',
              style: TextStyle(
                fontSize: 18,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '快去运动获得勋章吧',
              style: TextStyle(
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    return GridView.builder(
      padding: const EdgeInsets.all(AppDimensions.pagePadding),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 3,
        mainAxisSpacing: 16,
        crossAxisSpacing: 16,
        childAspectRatio: 0.8,
      ),
      itemCount: badges.length,
      itemBuilder: (context, index) {
        final badge = badges[index];
        return _buildBadgeItem(badge);
      },
    );
  }

  /// 构建勋章项
  Widget _buildBadgeItem(Badge badge) {
    return Container(
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
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 60,
            height: 60,
            decoration: BoxDecoration(
              color: _getBadgeColor(badge.category).withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(
              _getBadgeIcon(badge.category),
              size: 32,
              color: _getBadgeColor(badge.category),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            badge.name,
            style: const TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
            textAlign: TextAlign.center,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }

  /// 获取勋章颜色
  Color _getBadgeColor(BadgeCategory? category) {
    switch (category) {
      case BadgeCategory.milestone:
        return AppColors.primary;
      case BadgeCategory.skill:
        return AppColors.secondary;
      case BadgeCategory.streak:
        return AppColors.accent;
      case BadgeCategory.challenge:
        return AppColors.success;
      case BadgeCategory.vision:
        return Colors.blue;
      default:
        return AppColors.primary;
    }
  }

  /// 获取勋章图标
  IconData _getBadgeIcon(BadgeCategory? category) {
    switch (category) {
      case BadgeCategory.milestone:
        return Icons.star;
      case BadgeCategory.skill:
        return Icons.military_tech;
      case BadgeCategory.streak:
        return Icons.local_fire_department;
      case BadgeCategory.challenge:
        return Icons.emoji_events;
      case BadgeCategory.family:
        return Icons.family_restroom;
      case BadgeCategory.vision:
        return Icons.visibility;
      case BadgeCategory.special:
        return Icons.auto_awesome;
      default:
        return Icons.emoji_events;
    }
  }

  /// 构建积分商城
  Widget _buildPointsMall() {
    return ListView(
      padding: const EdgeInsets.all(AppDimensions.pagePadding),
      children: [
        const Text(
          '虚拟兑换',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),
        _buildExchangeItem(
          title: '角色装扮',
          description: '解锁新服装',
          icon: Icons.checkroom,
          points: 100,
        ),
        _buildExchangeItem(
          title: '场景主题',
          description: '解锁新场景',
          icon: Icons.landscape,
          points: 200,
        ),
        _buildExchangeItem(
          title: '虚拟宠物',
          description: '领养小宠物',
          icon: Icons.pets,
          points: 300,
        ),
        const SizedBox(height: AppDimensions.spacingLG),
        const Text(
          '公益兑换',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),
        _buildExchangeItem(
          title: '体育器材捐赠',
          description: '为山区儿童捐赠体育用品',
          icon: Icons.volunteer_activism,
          points: 500,
          isCharity: true,
        ),
      ],
    );
  }

  /// 构建兑换项
  Widget _buildExchangeItem({
    required String title,
    required String description,
    required IconData icon,
    required int points,
    bool isCharity = false,
  }) {
    return Card(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacingMD),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: isCharity
                ? AppColors.success.withOpacity(0.1)
                : AppColors.primary.withOpacity(0.1),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Icon(
            icon,
            color: isCharity ? AppColors.success : AppColors.primary,
          ),
        ),
        title: Text(title),
        subtitle: Text(description),
        trailing: ElevatedButton(
          onPressed: () {},
          style: ElevatedButton.styleFrom(
            backgroundColor: isCharity ? AppColors.success : AppColors.primary,
            foregroundColor: Colors.white,
          ),
          child: Text('$points积分'),
        ),
      ),
    );
  }
}
