import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/reward.dart';
import '../../data/datasources/reward_remote_data_source.dart';

/// 激励远程数据源Provider
final rewardRemoteDataSourceProvider =
    Provider<RewardRemoteDataSource>((ref) {
  return RewardRemoteDataSource();
});

/// 激励状态数据类
class RewardState {
  /// 所有勋章列表
  final List<Badge> badges;

  /// 我的勋章列表
  final List<Badge> myBadges;

  /// 积分记录列表
  final List<PointRecord> pointRecords;

  /// 积分余额
  final int pointsBalance;

  /// 挑战列表
  final List<Challenge> challenges;

  /// 家庭排行榜
  final List<LeaderboardEntry> familyLeaderboard;

  /// 全球排行榜
  final List<LeaderboardEntry> globalLeaderboard;

  /// 是否加载中
  final bool isLoading;

  /// 错误消息
  final String? error;

  const RewardState({
    this.badges = const [],
    this.myBadges = const [],
    this.pointRecords = const [],
    this.pointsBalance = 0,
    this.challenges = const [],
    this.familyLeaderboard = const [],
    this.globalLeaderboard = const [],
    this.isLoading = false,
    this.error,
  });

  RewardState copyWith({
    List<Badge>? badges,
    List<Badge>? myBadges,
    List<PointRecord>? pointRecords,
    int? pointsBalance,
    List<Challenge>? challenges,
    List<LeaderboardEntry>? familyLeaderboard,
    List<LeaderboardEntry>? globalLeaderboard,
    bool? isLoading,
    String? error,
  }) {
    return RewardState(
      badges: badges ?? this.badges,
      myBadges: myBadges ?? this.myBadges,
      pointRecords: pointRecords ?? this.pointRecords,
      pointsBalance: pointsBalance ?? this.pointsBalance,
      challenges: challenges ?? this.challenges,
      familyLeaderboard: familyLeaderboard ?? this.familyLeaderboard,
      globalLeaderboard: globalLeaderboard ?? this.globalLeaderboard,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// 激励Provider
final rewardProvider =
    StateNotifierProvider<RewardNotifier, RewardState>((ref) {
  final dataSource = ref.watch(rewardRemoteDataSourceProvider);
  return RewardNotifier(dataSource);
});

/// 激励Notifier
class RewardNotifier extends StateNotifier<RewardState> {
  /// 激励远程数据源
  final RewardRemoteDataSource _dataSource;

  RewardNotifier(this._dataSource) : super(const RewardState());

  /// 加载勋章列表
  /// 调用API获取所有勋章，可按分类筛选
  Future<void> loadBadges({String? category}) async {
    try {
      final response = await _dataSource.getBadges(category: category);
      if (response.isSuccess && response.data != null) {
        final badges = response.data!
            .map((dto) => _convertBadgeDTO(dto))
            .toList();
        state = state.copyWith(badges: badges);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载我的勋章
  /// 调用API获取当前用户已获得的勋章
  Future<void> loadMyBadges() async {
    try {
      final response = await _dataSource.getMyBadges();
      if (response.isSuccess && response.data != null) {
        final myBadges = response.data!
            .map((dto) => _convertBadgeDTO(dto))
            .toList();
        state = state.copyWith(myBadges: myBadges);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载积分记录
  /// 调用API获取积分变动记录，支持分页
  Future<void> loadPointRecords({
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dataSource.getPoints(
        page: page,
        pageSize: pageSize,
      );
      if (response.isSuccess && response.data != null) {
        final records = response.data!.list
            .map((dto) => _convertPointRecordDTO(dto))
            .toList();
        state = state.copyWith(pointRecords: records);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载积分余额
  /// 调用API获取当前用户积分余额
  Future<void> loadPointsBalance() async {
    try {
      final response = await _dataSource.getPointsBalance();
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(pointsBalance: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 创建挑战
  /// 调用API发起新挑战
  Future<bool> createChallenge(CreateChallengeRequest req) async {
    try {
      final response = await _dataSource.createChallenge(req);
      if (response.isSuccess && response.data != null) {
        final challenge = _convertChallengeDTO(response.data!);
        state = state.copyWith(challenges: [challenge, ...state.challenges]);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载挑战列表
  /// 调用API获取挑战列表
  Future<void> loadChallenges() async {
    try {
      final response = await _dataSource.getChallenges();
      if (response.isSuccess && response.data != null) {
        final challenges = response.data!
            .map((dto) => _convertChallengeDTO(dto))
            .toList();
        state = state.copyWith(challenges: challenges);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 接受挑战
  /// 调用API接受指定挑战
  Future<bool> acceptChallenge(String challengeId) async {
    try {
      final response = await _dataSource.acceptChallenge(challengeId);
      if (response.isSuccess && response.data != null) {
        final updatedChallenge = _convertChallengeDTO(response.data!);
        final updatedChallenges = state.challenges.map((c) {
          if (c.id == challengeId) return updatedChallenge;
          return c;
        }).toList();
        state = state.copyWith(challenges: updatedChallenges);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 提交挑战得分
  /// 调用API提交挑战成绩
  Future<bool> submitChallengeScore(String challengeId, int score) async {
    try {
      final response =
          await _dataSource.submitChallengeScore(challengeId, score);
      if (response.isSuccess && response.data != null) {
        final updatedChallenge = _convertChallengeDTO(response.data!);
        final updatedChallenges = state.challenges.map((c) {
          if (c.id == challengeId) return updatedChallenge;
          return c;
        }).toList();
        state = state.copyWith(challenges: updatedChallenges);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载家庭排行榜
  /// 调用API获取家庭成员排行
  Future<void> loadFamilyLeaderboard() async {
    try {
      final response = await _dataSource.getFamilyLeaderboard();
      if (response.isSuccess && response.data != null) {
        final entries = response.data!
            .map((dto) => _convertLeaderboardDTO(dto))
            .toList();
        state = state.copyWith(familyLeaderboard: entries);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载全球排行榜
  /// 调用API获取全球排行，可按运动类型筛选
  Future<void> loadGlobalLeaderboard({String? exerciseType}) async {
    try {
      final response =
          await _dataSource.getGlobalLeaderboard(exerciseType: exerciseType);
      if (response.isSuccess && response.data != null) {
        final entries = response.data!
            .map((dto) => _convertLeaderboardDTO(dto))
            .toList();
        state = state.copyWith(globalLeaderboard: entries);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  // ==================== 私有辅助方法 ====================

  /// 将BadgeDTO转换为本地Badge实体
  Badge _convertBadgeDTO(BadgeDTO dto) {
    return Badge(
      id: dto.id,
      code: dto.code,
      name: dto.name,
      description: dto.description,
      category: _parseBadgeCategory(dto.category),
      icon: dto.icon.isNotEmpty ? dto.icon : null,
      points: dto.points,
      createdAt: DateTime.now(),
    );
  }

  /// 将PointRecordDTO转换为本地PointRecord实体
  PointRecord _convertPointRecordDTO(PointRecordDTO dto) {
    return PointRecord(
      id: dto.id,
      userId: dto.userId,
      points: dto.points,
      type: _parsePointType(dto.type),
      sourceId: dto.sourceId,
      sourceType: dto.sourceType,
      description: dto.description,
      balance: dto.balance,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将ChallengeDTO转换为本地Challenge实体
  Challenge _convertChallengeDTO(ChallengeDTO dto) {
    return Challenge(
      id: dto.id,
      type: _parseChallengeType(dto.type),
      initiatorId: dto.initiatorId,
      acceptorId: dto.acceptorId ?? '',
      exerciseType: dto.exerciseType,
      targetValue: dto.targetValue,
      initiatorScore: dto.initiatorScore,
      acceptorScore: dto.acceptorScore,
      winnerId: dto.winnerId,
      status: _parseChallengeStatus(dto.status),
      expiresAt: dto.expiresAt != null && dto.expiresAt!.isNotEmpty
          ? DateTime.parse(dto.expiresAt!)
          : DateTime.now().add(const Duration(days: 7)),
      completedAt: dto.completedAt != null && dto.completedAt!.isNotEmpty
          ? DateTime.parse(dto.completedAt!)
          : null,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将LeaderboardDTO转换为本地LeaderboardEntry实体
  LeaderboardEntry _convertLeaderboardDTO(LeaderboardDTO dto) {
    return LeaderboardEntry(
      rank: dto.rank,
      userId: dto.userId,
      nickname: dto.nickname,
      avatar: dto.avatar,
      score: dto.score.round(), // double转int
    );
  }

  /// 解析勋章分类字符串为枚举
  BadgeCategory _parseBadgeCategory(String value) {
    return BadgeCategory.values.firstWhere(
      (e) => e.value == value,
      orElse: () => BadgeCategory.milestone,
    );
  }

  /// 解析积分类型字符串为枚举
  PointType _parsePointType(String value) {
    return PointType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => PointType.exercise,
    );
  }

  /// 解析挑战类型字符串为枚举
  ChallengeType _parseChallengeType(String value) {
    return ChallengeType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ChallengeType.async,
    );
  }

  /// 解析挑战状态字符串为枚举
  ChallengeStatus _parseChallengeStatus(String value) {
    return ChallengeStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ChallengeStatus.pending,
    );
  }
}
