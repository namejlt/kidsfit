import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/reward.dart';
import '../../services/storage_service.dart';

/// 用户勋章列表Provider
final userBadgesProvider = StateNotifierProvider<UserBadgesNotifier, List<UserBadge>>((ref) {
  return UserBadgesNotifier();
});

/// 用户勋章Notifier
class UserBadgesNotifier extends StateNotifier<List<UserBadge>> {
  UserBadgesNotifier() : super([]);

  /// 加载用户勋章
  Future<void> loadUserBadges(String userId) async {
    try {
      // TODO: 调用API获取用户勋章
      await Future.delayed(const Duration(milliseconds: 300));

      // 模拟数据
      final badges = [
        Badge(
          id: 'badge_1',
          code: 'first_exercise',
          name: '初次运动',
          description: '完成第一次运动',
          category: BadgeCategory.milestone,
          points: 10,
          createdAt: DateTime.now(),
        ),
        Badge(
          id: 'badge_2',
          code: 'jump_master',
          name: '跳绳达人',
          description: '跳绳得分超过80分',
          category: BadgeCategory.skill,
          points: 20,
          createdAt: DateTime.now(),
        ),
      ];

      state = badges.map((badge) => UserBadge(
        id: 'user_badge_${badge.id}',
        userId: userId,
        badgeId: badge.id,
        earnedAt: DateTime.now(),
        badge: badge,
        createdAt: DateTime.now(),
      )).toList();
    } catch (e) {
      state = [];
    }
  }

  /// 添加勋章
  Future<void> addBadge(UserBadge badge) async {
    state = [...state, badge];
  }

  /// 获取新勋章
  List<UserBadge> getNewBadges() {
    return state.where((b) => b.isNew).toList();
  }
}

/// 用户积分Provider
final userPointsProvider = StateProvider<int>((ref) => 0);

/// 积分记录列表Provider
final pointRecordsProvider =
    StateNotifierProvider<PointRecordsNotifier, List<PointRecord>>((ref) {
  return PointRecordsNotifier();
});

/// 积分记录Notifier
class PointRecordsNotifier extends StateNotifier<List<PointRecord>> {
  PointRecordsNotifier() : super([]);

  /// 加载积分记录
  Future<void> loadPointRecords(String userId) async {
    try {
      // TODO: 调用API获取积分记录
      state = [];
    } catch (e) {
      state = [];
    }
  }

  /// 添加积分记录
  Future<void> addPointRecord(PointRecord record) async {
    state = [record, ...state];
  }
}

/// 激励Provider
final rewardProvider =
    StateNotifierProvider<RewardNotifier, RewardState>((ref) {
  return RewardNotifier(ref);
});

/// 激励状态
class RewardState {
  final bool isLoading;
  final int totalPoints;
  final int badgeCount;
  final int streakDays;
  final String? newBadgeCode;
  final String? error;

  const RewardState({
    this.isLoading = false,
    this.totalPoints = 0,
    this.badgeCount = 0,
    this.streakDays = 0,
    this.newBadgeCode,
    this.error,
  });

  RewardState copyWith({
    bool? isLoading,
    int? totalPoints,
    int? badgeCount,
    int? streakDays,
    String? newBadgeCode,
    String? error,
  }) {
    return RewardState(
      isLoading: isLoading ?? this.isLoading,
      totalPoints: totalPoints ?? this.totalPoints,
      badgeCount: badgeCount ?? this.badgeCount,
      streakDays: streakDays ?? this.streakDays,
      newBadgeCode: newBadgeCode,
      error: error,
    );
  }
}

/// 激励Notifier
class RewardNotifier extends StateNotifier<RewardState> {
  final Ref _ref;

  RewardNotifier(this._ref) : super(const RewardState());

  /// 加载用户奖励信息
  Future<void> loadUserRewards() async {
    state = state.copyWith(isLoading: true);

    try {
      // TODO: 调用API获取用户奖励信息
      await Future.delayed(const Duration(milliseconds: 300));

      state = state.copyWith(
        isLoading: false,
        totalPoints: 100,
        badgeCount: 5,
        streakDays: 7,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 添加积分
  Future<void> addPoints({
    required int points,
    required PointType type,
    String? description,
  }) async {
    final newTotal = state.totalPoints + points;
    state = state.copyWith(totalPoints: newTotal);

    // 记录积分变动
    final record = PointRecord(
      id: 'point_${DateTime.now().millisecondsSinceEpoch}',
      userId: 'current_user',
      points: points,
      type: type,
      description: description,
      balance: newTotal,
      createdAt: DateTime.now(),
    );

    _ref.read(pointRecordsProvider.notifier).addPointRecord(record);
  }

  /// 检查并颁发新勋章
  Future<Badge?> checkAndAwardBadge(String badgeCode) async {
    // TODO: 调用API检查勋章条件并颁发
    return null;
  }

  /// 设置新勋章提示
  void setNewBadge(String badgeCode) {
    state = state.copyWith(
      newBadgeCode: badgeCode,
      badgeCount: state.badgeCount + 1,
    );
  }

  /// 清除新勋章提示
  void clearNewBadge() {
    state = state.copyWith(newBadgeCode: null);
  }
}

/// 挑战列表Provider
final challengesProvider =
    StateNotifierProvider<ChallengesNotifier, List<Challenge>>((ref) {
  return ChallengesNotifier();
});

/// 挑战Notifier
class ChallengesNotifier extends StateNotifier<List<Challenge>> {
  ChallengesNotifier() : super([]);

  /// 加载挑战列表
  Future<void> loadChallenges(String userId) async {
    try {
      // TODO: 调用API获取挑战列表
      state = [];
    } catch (e) {
      state = [];
    }
  }

  /// 添加挑战
  Future<void> addChallenge(Challenge challenge) async {
    state = [challenge, ...state];
  }

  /// 更新挑战状态
  Future<void> updateChallenge(Challenge challenge) async {
    state = state.map((c) => c.id == challenge.id ? challenge : c).toList();
  }
}
