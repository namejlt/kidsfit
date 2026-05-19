import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../domain/entities/vision_record.dart';
import '../../providers/vision_provider.dart';
import '../../widgets/stat_card.dart';

/// 视力管理页面
/// 展示视力档案、户外时间和用眼提醒
class VisionScreen extends ConsumerStatefulWidget {
  const VisionScreen({super.key});

  @override
  ConsumerState<VisionScreen> createState() => _VisionScreenState();
}

class _VisionScreenState extends ConsumerState<VisionScreen> {
  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    // 加载今日户外
    ref.read(visionProvider.notifier).loadTodayOutdoor();
  }

  @override
  Widget build(BuildContext context) {
    final visionState = ref.watch(visionProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('视力健康'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showAddVisionDialog(),
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
              // 户外时间
              _buildOutdoorSection(visionState.todayOutdoor),

              const SizedBox(height: AppDimensions.spacingLG),

              // 视力档案
              _buildVisionRecordsSection(visionState.visionRecords),

              const SizedBox(height: AppDimensions.spacingLG),

              // 用眼提醒设置
              _buildEyeReminderSection(),

              const SizedBox(height: AppDimensions.spacingLG),

              // 护眼小贴士
              _buildEyeCareTips(),
            ],
          ),
        ),
      ),
    );
  }

  /// 构建户外时间区域
  Widget _buildOutdoorSection(OutdoorActivity? outdoor) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          AppStrings.outdoorTime,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        Container(
          padding: const EdgeInsets.all(AppDimensions.cardPadding),
          decoration: BoxDecoration(
            gradient: AppColors.outdoorGradient,
            borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
            boxShadow: [
              BoxShadow(
                color: AppColors.success.withOpacity(0.3),
                blurRadius: 10,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: Column(
            children: [
              Row(
                children: [
                  const Icon(
                    Icons.wb_sunny,
                    color: Colors.white,
                    size: 32,
                  ),
                  const SizedBox(width: 12),
                  const Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '今日户外运动',
                          style: TextStyle(
                            color: Colors.white70,
                            fontSize: 14,
                          ),
                        ),
                        Text(
                          '保护视力的关键',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  if (outdoor?.isTargetMet ?? false)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 6,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: const Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.check_circle,
                            color: Colors.white,
                            size: 16,
                          ),
                          SizedBox(width: 4),
                          Text(
                            '已达标',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                ],
              ),
              const SizedBox(height: AppDimensions.spacingMD),

              // 进度条
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: LinearProgressIndicator(
                  value: (outdoor?.targetProgress ?? 0) / 100,
                  backgroundColor: Colors.white.withOpacity(0.3),
                  valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
                  minHeight: 10,
                ),
              ),

              const SizedBox(height: AppDimensions.spacingSM),

              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    outdoor?.formattedDuration ?? '0分钟',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Text(
                    '目标: ${OutdoorActivity.targetMinutes}分钟',
                    style: const TextStyle(
                      color: Colors.white70,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// 构建视力记录区域
  Widget _buildVisionRecordsSection(List<VisionRecord> records) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              AppStrings.visionRecord,
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            TextButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.history, size: 18),
              label: const Text('历史记录'),
            ),
          ],
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        if (records.isEmpty)
          _buildEmptyVisionRecord()
        else
          ...records.take(2).map((record) => _buildVisionRecordCard(record)),
      ],
    );
  }

  /// 构建空视力记录
  Widget _buildEmptyVisionRecord() {
    return Container(
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
        border: Border.all(color: AppColors.background),
      ),
      child: Center(
        child: Column(
          children: [
            Icon(
              Icons.visibility_off,
              size: 48,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 12),
            const Text(
              '暂无视力档案',
              style: TextStyle(
                color: AppColors.textSecondary,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 8),
            ElevatedButton.icon(
              onPressed: () => _showAddVisionDialog(),
              icon: const Icon(Icons.add),
              label: const Text(AppStrings.addRecord),
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

  /// 构建视力记录卡片
  Widget _buildVisionRecordCard(VisionRecord record) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacingMD),
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
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '${record.date.month}月${record.date.day}日',
                style: const TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: _getVisionStatusColor(record.visionStatus)
                      .withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  record.visionStatus.displayName,
                  style: TextStyle(
                    fontSize: 12,
                    color: _getVisionStatusColor(record.visionStatus),
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: AppDimensions.spacingMD),

          // 双眼数据
          Row(
            children: [
              Expanded(
                child: _buildEyeDataColumn('右眼', record.rightEye),
              ),
              Container(
                width: 1,
                height: 50,
                color: AppColors.background,
              ),
              Expanded(
                child: _buildEyeDataColumn('左眼', record.leftEye),
              ),
            ],
          ),

          const SizedBox(height: AppDimensions.spacingMD),

          // 视力状态说明
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: AppColors.background.withOpacity(0.5),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              children: [
                const Icon(
                  Icons.info_outline,
                  size: 16,
                  color: AppColors.textSecondary,
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    _getVisionAdvice(record),
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// 构建单眼数据列
  Widget _buildEyeDataColumn(String label, EyeData eye) {
    return Column(
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '${eye.sph >= 0 ? '+' : ''}${eye.sph}',
          style: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        Text(
          '矫正视力: ${eye.va}',
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  /// 获取视力状态颜色
  Color _getVisionStatusColor(VisionStatus status) {
    switch (status) {
      case VisionStatus.good:
        return AppColors.visionGood;
      case VisionStatus.medium:
        return AppColors.visionMedium;
      case VisionStatus.concern:
        return AppColors.visionConcern;
    }
  }

  /// 获取视力建议
  String _getVisionAdvice(VisionRecord record) {
    if (record.averageVa >= 1.0) {
      return '视力发育良好，继续保持每天2小时户外活动';
    } else if (record.averageVa >= 0.8) {
      return '视力轻度下降，建议增加户外运动时间';
    } else {
      return '建议尽快到专业机构进行详细检查';
    }
  }

  /// 构建用眼提醒区域
  Widget _buildEyeReminderSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '用眼提醒',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        Container(
          padding: const EdgeInsets.all(AppDimensions.cardPadding),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: BorderRadius.circular(AppDimensions.cardRadius),
          ),
          child: Column(
            children: [
              _buildReminderItem(
                icon: Icons.timer,
                title: '20-20-20法则',
                subtitle: '每20分钟看近处，看20英尺外20秒',
                isEnabled: true,
              ),
              const Divider(),
              _buildReminderItem(
                icon: Icons.directions_walk,
                title: '户外提醒',
                subtitle: '连续室内45分钟后提醒远眺',
                isEnabled: true,
              ),
              const Divider(),
              _buildReminderItem(
                icon: Icons.sports,
                title: '户外运动提醒',
                subtitle: '每日120分钟户外运动目标',
                isEnabled: true,
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// 构建提醒项
  Widget _buildReminderItem({
    required IconData icon,
    required String title,
    required String subtitle,
    required bool isEnabled,
  }) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(
              icon,
              color: AppColors.primary,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                Text(
                  subtitle,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Switch(
            value: isEnabled,
            onChanged: (value) {},
            activeColor: AppColors.primary,
          ),
        ],
      ),
    );
  }

  /// 构建护眼小贴士
  Widget _buildEyeCareTips() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '护眼小贴士',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),

        _buildTipCard(
          icon: Icons.wb_sunny,
          title: '多晒太阳',
          description: '户外自然光可有效预防近视',
          color: AppColors.accent,
        ),
        _buildTipCard(
          icon: Icons.visibility,
          title: '正确用眼',
          description: '保持30-40cm观看距离',
          color: AppColors.secondary,
        ),
        _buildTipCard(
          icon: Icons.restaurant,
          title: '营养补充',
          description: '多吃胡萝卜、菠菜等护眼食物',
          color: AppColors.success,
        ),
      ],
    );
  }

  /// 构建提示卡片
  Widget _buildTipCard({
    required IconData icon,
    required String title,
    required String description,
    required Color color,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacingMD),
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
      ),
      child: Row(
        children: [
          Icon(
            icon,
            color: color,
            size: 32,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: color,
                  ),
                ),
                Text(
                  description,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// 显示添加视力记录对话框
  void _showAddVisionDialog() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.6,
          minChildSize: 0.4,
          maxChildSize: 0.9,
          expand: false,
          builder: (context, scrollController) {
            return Container(
              padding: const EdgeInsets.all(AppDimensions.pagePadding),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                        color: AppColors.background,
                        borderRadius: BorderRadius.circular(2),
                      ),
                    ),
                  ),
                  const SizedBox(height: AppDimensions.spacingMD),
                  const Text(
                    '添加视力记录',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: AppDimensions.spacingLG),

                  Expanded(
                    child: ListView(
                      controller: scrollController,
                      children: [
                        ListTile(
                          leading: Container(
                            width: 48,
                            height: 48,
                            decoration: BoxDecoration(
                              color: AppColors.primary.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: const Icon(
                              Icons.camera_alt,
                              color: AppColors.primary,
                            ),
                          ),
                          title: const Text('拍照识别验光单'),
                          subtitle: const Text('OCR自动识别'),
                          onTap: () {
                            // TODO: 打开相机拍照
                          },
                        ),
                        const Divider(),
                        ListTile(
                          leading: Container(
                            width: 48,
                            height: 48,
                            decoration: BoxDecoration(
                              color: AppColors.secondary.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: const Icon(
                              Icons.edit,
                              color: AppColors.secondary,
                            ),
                          ),
                          title: const Text('手动录入'),
                          subtitle: const Text('输入验光数据'),
                          onTap: () {
                            // TODO: 打开手动录入表单
                          },
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }
}
