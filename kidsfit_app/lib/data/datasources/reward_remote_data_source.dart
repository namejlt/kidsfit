import '../../core/network/api_client.dart';
import '../../core/network/api_response.dart';
import '../../core/constants/app_constants.dart';

// ==================== DTO类定义 ====================

/// 徽章DTO
class BadgeDTO {
  /// 徽章ID
  final String id;

  /// 徽章代码
  final String code;

  /// 徽章名称
  final String name;

  /// 徽章描述
  final String description;

  /// 徽章分类
  final String category;

  /// 徽章图标
  final String icon;

  /// 徽章积分
  final int points;

  /// 是否已获得
  final bool earned;

  const BadgeDTO({
    required this.id,
    required this.code,
    required this.name,
    this.description = '',
    this.category = '',
    this.icon = '',
    this.points = 0,
    this.earned = false,
  });

  /// 从JSON创建BadgeDTO
  factory BadgeDTO.fromJson(Map<String, dynamic> json) {
    return BadgeDTO(
      id: json['id'] as String? ?? '',
      code: json['code'] as String? ?? '',
      name: json['name'] as String? ?? '',
      description: json['description'] as String? ?? '',
      category: json['category'] as String? ?? '',
      icon: json['icon'] as String? ?? '',
      points: json['points'] as int? ?? 0,
      earned: json['earned'] as bool? ?? false,
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'code': code,
        'name': name,
        'description': description,
        'category': category,
        'icon': icon,
        'points': points,
        'earned': earned,
      };
}

/// 积分记录DTO
class PointRecordDTO {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 积分数量
  final int points;

  /// 类型（earn/spend）
  final String type;

  /// 来源ID
  final String? sourceId;

  /// 来源类型
  final String? sourceType;

  /// 描述
  final String description;

  /// 余额
  final int balance;

  /// 创建时间
  final String createdAt;

  const PointRecordDTO({
    required this.id,
    required this.userId,
    required this.points,
    required this.type,
    this.sourceId,
    this.sourceType,
    this.description = '',
    this.balance = 0,
    required this.createdAt,
  });

  /// 从JSON创建PointRecordDTO
  factory PointRecordDTO.fromJson(Map<String, dynamic> json) {
    return PointRecordDTO(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      points: json['points'] as int? ?? 0,
      type: json['type'] as String? ?? '',
      sourceId: json['source_id'] as String?,
      sourceType: json['source_type'] as String?,
      description: json['description'] as String? ?? '',
      balance: json['balance'] as int? ?? 0,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'user_id': userId,
        'points': points,
        'type': type,
        'source_id': sourceId,
        'source_type': sourceType,
        'description': description,
        'balance': balance,
        'created_at': createdAt,
      };
}

/// 挑战DTO
class ChallengeDTO {
  /// 挑战ID
  final String id;

  /// 挑战类型
  final String type;

  /// 发起者ID
  final String initiatorId;

  /// 接受者ID
  final String? acceptorId;

  /// 运动类型
  final String exerciseType;

  /// 目标值
  final int targetValue;

  /// 发起者得分
  final int? initiatorScore;

  /// 接受者得分
  final int? acceptorScore;

  /// 获胜者ID
  final String? winnerId;

  /// 挑战状态
  final String status;

  /// 过期时间
  final String? expiresAt;

  /// 完成时间
  final String? completedAt;

  /// 创建时间
  final String createdAt;

  const ChallengeDTO({
    required this.id,
    required this.type,
    required this.initiatorId,
    this.acceptorId,
    required this.exerciseType,
    required this.targetValue,
    this.initiatorScore,
    this.acceptorScore,
    this.winnerId,
    this.status = 'pending',
    this.expiresAt,
    this.completedAt,
    required this.createdAt,
  });

  /// 从JSON创建ChallengeDTO
  factory ChallengeDTO.fromJson(Map<String, dynamic> json) {
    return ChallengeDTO(
      id: json['id'] as String? ?? '',
      type: json['type'] as String? ?? '',
      initiatorId: json['initiator_id'] as String? ?? '',
      acceptorId: json['acceptor_id'] as String?,
      exerciseType: json['exercise_type'] as String? ?? '',
      targetValue: json['target_value'] as int? ?? 0,
      initiatorScore: json['initiator_score'] as int?,
      acceptorScore: json['acceptor_score'] as int?,
      winnerId: json['winner_id'] as String?,
      status: json['status'] as String? ?? 'pending',
      expiresAt: json['expires_at'] as String?,
      completedAt: json['completed_at'] as String?,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type,
        'initiator_id': initiatorId,
        'acceptor_id': acceptorId,
        'exercise_type': exerciseType,
        'target_value': targetValue,
        'initiator_score': initiatorScore,
        'acceptor_score': acceptorScore,
        'winner_id': winnerId,
        'status': status,
        'expires_at': expiresAt,
        'completed_at': completedAt,
        'created_at': createdAt,
      };
}

/// 创建挑战请求DTO
class CreateChallengeRequest {
  /// 挑战类型
  final String type;

  /// 接受者ID
  final String? acceptorId;

  /// 运动类型
  final String exerciseType;

  /// 目标值
  final int targetValue;

  const CreateChallengeRequest({
    required this.type,
    this.acceptorId,
    required this.exerciseType,
    required this.targetValue,
  });

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'type': type,
        'acceptor_id': acceptorId,
        'exercise_type': exerciseType,
        'target_value': targetValue,
      };
}

/// 排行榜DTO
class LeaderboardDTO {
  /// 排名
  final int rank;

  /// 用户ID
  final String userId;

  /// 昵称
  final String nickname;

  /// 头像URL
  final String? avatar;

  /// 得分（double匹配后端float64）
  final double score;

  const LeaderboardDTO({
    required this.rank,
    required this.userId,
    required this.nickname,
    this.avatar,
    required this.score,
  });

  /// 从JSON创建LeaderboardDTO
  factory LeaderboardDTO.fromJson(Map<String, dynamic> json) {
    return LeaderboardDTO(
      rank: json['rank'] as int? ?? 0,
      userId: json['user_id'] as String? ?? '',
      nickname: json['nickname'] as String? ?? '',
      avatar: json['avatar'] as String?,
      score: (json['score'] as num?)?.toDouble() ?? 0,
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'rank': rank,
        'user_id': userId,
        'nickname': nickname,
        'avatar': avatar,
        'score': score,
      };
}

// ==================== 激励远程数据源 ====================

/// 激励远程数据源
/// 负责与后端激励/奖励相关API通信
class RewardRemoteDataSource {
  /// API客户端实例
  final ApiClient _apiClient;

  RewardRemoteDataSource({ApiClient? apiClient})
      : _apiClient = apiClient ?? ApiClient.instance;

  /// 获取徽章列表
  /// GET /api/v1/rewards/badges
  Future<ApiResponse<List<BadgeDTO>>> getBadges({String? category}) async {
    final queryParams = <String, dynamic>{};
    if (category != null) {
      queryParams['category'] = category;
    }

    return _apiClient.get(
      AppConstants.badges,
      queryParameters: queryParams,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => BadgeDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 获取我的徽章
  /// GET /api/v1/rewards/badges/mine
  Future<ApiResponse<List<BadgeDTO>>> getMyBadges() async {
    return _apiClient.get(
      AppConstants.myBadges,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => BadgeDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 获取积分记录（分页）
  /// GET /api/v1/rewards/points
  Future<ApiResponse<PaginatedData<PointRecordDTO>>> getPoints({
    int page = 1,
    int pageSize = 20,
  }) async {
    return _apiClient.getPaged(
      AppConstants.points,
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
      fromJsonT: (data) =>
          PointRecordDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取积分余额
  /// GET /api/v1/rewards/points/balance
  Future<ApiResponse<int>> getPointsBalance() async {
    return _apiClient.get(
      AppConstants.pointsBalance,
      fromJsonT: (data) => data as int? ?? 0,
    );
  }

  /// 创建挑战
  /// POST /api/v1/rewards/challenges
  Future<ApiResponse<ChallengeDTO>> createChallenge(
      CreateChallengeRequest req) async {
    return _apiClient.post(
      AppConstants.challenges,
      data: req.toJson(),
      fromJsonT: (data) =>
          ChallengeDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取挑战列表
  /// GET /api/v1/rewards/challenges
  Future<ApiResponse<List<ChallengeDTO>>> getChallenges() async {
    return _apiClient.get(
      AppConstants.challenges,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => ChallengeDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 接受挑战
  /// POST /api/v1/rewards/challenges/{challengeId}/accept
  Future<ApiResponse<ChallengeDTO>> acceptChallenge(
      String challengeId) async {
    return _apiClient.post(
      '${AppConstants.challenges}/$challengeId/accept',
      fromJsonT: (data) =>
          ChallengeDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 提交挑战得分
  /// POST /api/v1/rewards/challenges/{challengeId}/score
  Future<ApiResponse<ChallengeDTO>> submitChallengeScore(
      String challengeId, int score) async {
    return _apiClient.post(
      '${AppConstants.challenges}/$challengeId/score',
      data: {'score': score},
      fromJsonT: (data) =>
          ChallengeDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取家庭排行榜
  /// GET /api/v1/rewards/leaderboard/family
  Future<ApiResponse<List<LeaderboardDTO>>> getFamilyLeaderboard() async {
    return _apiClient.get(
      AppConstants.familyLeaderboard,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => LeaderboardDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 获取全球排行榜
  /// GET /api/v1/rewards/leaderboard/global
  Future<ApiResponse<List<LeaderboardDTO>>> getGlobalLeaderboard({
    String? exerciseType,
    int? limit,
  }) async {
    final queryParams = <String, dynamic>{};
    if (exerciseType != null) {
      queryParams['exercise_type'] = exerciseType;
    }
    if (limit != null) {
      queryParams['limit'] = limit;
    }

    return _apiClient.get(
      AppConstants.globalLeaderboard,
      queryParameters: queryParams,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => LeaderboardDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }
}
