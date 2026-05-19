import 'package:equatable/equatable.dart';

/// 勋章分类枚举
enum BadgeCategory {
  /// 里程碑
  milestone('milestone', '里程碑'),

  /// 技能认证
  skill('skill', '技能认证'),

  /// 连续打卡
  streak('streak', '连续打卡'),

  /// 挑战胜利
  challenge('challenge', '挑战胜利'),

  /// 亲子互动
  family('family', '亲子互动'),

  /// 视力守护
  vision('vision', '视力守护'),

  /// 特殊成就
  special('special', '特殊成就');

  const BadgeCategory(this.value, this.displayName);
  final String value;
  final String displayName;
}

/// 积分类型枚举
enum PointType {
  /// 运动积分
  exercise('exercise', '运动'),

  /// 记录突破积分
  recordBreak('record_break', '记录突破'),

  /// 亲子活动积分
  familyActivity('family_activity', '亲子活动'),

  /// 连续打卡积分
  streak('streak', '连续打卡'),

  /// 护眼任务积分
  visionTask('vision_task', '护眼任务'),

  /// 积分兑换
  redeem('redeem', '兑换');

  const PointType(this.value, this.displayName);
  final String value;
  final String displayName;
}

/// 挑战类型枚举
enum ChallengeType {
  /// 同步挑战
  sync('sync', '同步挑战'),

  /// 异步挑战
  async('async', '异步挑战'),

  /// 限时挑战
  timed('timed', '限时挑战');

  const ChallengeType(this.value, this.displayName);
  final String value;
  final String displayName;
}

/// 挑战状态枚举
enum ChallengeStatus {
  /// 待接受
  pending('pending'),

  /// 已接受
  accepted('accepted'),

  /// 已完成
  completed('completed'),

  /// 已过期
  expired('expired');

  const ChallengeStatus(this.value);
  final String value;
}

/// 勋章实体
class Badge extends Equatable {
  /// 勋章ID
  final String id;

  /// 勋章代码
  final String code;

  /// 勋章名称
  final String name;

  /// 勋章描述
  final String? description;

  /// 勋章分类
  final BadgeCategory category;

  /// 勋章图标URL
  final String? icon;

  /// 获得条件JSON
  final String? condition;

  /// 勋章积分
  final int points;

  /// 创建时间
  final DateTime createdAt;

  const Badge({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    required this.category,
    this.icon,
    this.condition,
    this.points = 0,
    required this.createdAt,
  });

  /// 判断是否为特殊勋章
  bool get isSpecial => category == BadgeCategory.special;

  /// 判断是否为技能勋章
  bool get isSkillBadge => category == BadgeCategory.skill;

  /// 判断是否为连续打卡勋章
  bool get isStreakBadge => category == BadgeCategory.streak;

  @override
  List<Object?> get props => [
        id,
        code,
        name,
        description,
        category,
        icon,
        condition,
        points,
        createdAt,
      ];
}

/// 用户勋章实体
class UserBadge extends Equatable {
  /// 用户勋章ID
  final String id;

  /// 用户ID
  final String userId;

  /// 勋章ID
  final String badgeId;

  /// 获得时间
  final DateTime earnedAt;

  /// 勋章详情
  final Badge? badge;

  /// 创建时间
  final DateTime createdAt;

  const UserBadge({
    required this.id,
    required this.userId,
    required this.badgeId,
    required this.earnedAt,
    this.badge,
    required this.createdAt,
  });

  /// 判断是否是新获得的（7天内）
  bool get isNew {
    return DateTime.now().difference(earnedAt).inDays <= 7;
  }

  UserBadge copyWith({
    String? id,
    String? userId,
    String? badgeId,
    DateTime? earnedAt,
    Badge? badge,
    DateTime? createdAt,
  }) {
    return UserBadge(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      badgeId: badgeId ?? this.badgeId,
      earnedAt: earnedAt ?? this.earnedAt,
      badge: badge ?? this.badge,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [id, userId, badgeId, earnedAt, badge, createdAt];
}

/// 积分记录实体
class PointRecord extends Equatable {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 积分变动
  final int points;

  /// 积分类型
  final PointType type;

  /// 来源ID
  final String? sourceId;

  /// 来源类型
  final String? sourceType;

  /// 描述
  final String? description;

  /// 变动后余额
  final int balance;

  /// 创建时间
  final DateTime createdAt;

  const PointRecord({
    required this.id,
    required this.userId,
    required this.points,
    required this.type,
    this.sourceId,
    this.sourceType,
    this.description,
    required this.balance,
    required this.createdAt,
  });

  /// 判断是否为收入
  bool get isIncome => points > 0;

  /// 判断是否为支出
  bool get isExpense => points < 0;

  /// 获取格式化积分
  String get formattedPoints {
    if (points > 0) return '+$points';
    return '$points';
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        points,
        type,
        sourceId,
        sourceType,
        description,
        balance,
        createdAt,
      ];
}

/// 挑战实体
class Challenge extends Equatable {
  /// 挑战ID
  final String id;

  /// 挑战类型
  final ChallengeType type;

  /// 发起者ID
  final String initiatorId;

  /// 接受者ID
  final String acceptorId;

  /// 运动类型
  final String exerciseType;

  /// 目标值
  final int targetValue;

  /// 发起者成绩
  final int? initiatorScore;

  /// 接受者成绩
  final int? acceptorScore;

  /// 获胜者ID
  final String? winnerId;

  /// 挑战状态
  final ChallengeStatus status;

  /// 过期时间
  final DateTime expiresAt;

  /// 完成时间
  final DateTime? completedAt;

  /// 创建时间
  final DateTime createdAt;

  const Challenge({
    required this.id,
    required this.type,
    required this.initiatorId,
    required this.acceptorId,
    required this.exerciseType,
    required this.targetValue,
    this.initiatorScore,
    this.acceptorScore,
    this.winnerId,
    this.status = ChallengeStatus.pending,
    required this.expiresAt,
    this.completedAt,
    required this.createdAt,
  });

  /// 判断是否已过期
  bool get isExpired => DateTime.now().isAfter(expiresAt);

  /// 判断是否等待中
  bool get isPending => status == ChallengeStatus.pending;

  /// 判断是否已完成
  bool get isCompleted => status == ChallengeStatus.completed;

  /// 判断发起者是否获胜
  bool isInitiatorWinner(String userId) =>
      initiatorId == userId && winnerId == userId;

  /// 判断接受者是否获胜
  bool isAcceptorWinner(String userId) =>
      acceptorId == userId && winnerId == userId;

  /// 判断用户是否为发起者
  bool isInitiator(String userId) => initiatorId == userId;

  /// 判断用户是否为接受者
  bool isAcceptor(String userId) => acceptorId == userId;

  Challenge copyWith({
    String? id,
    ChallengeType? type,
    String? initiatorId,
    String? acceptorId,
    String? exerciseType,
    int? targetValue,
    int? initiatorScore,
    int? acceptorScore,
    String? winnerId,
    ChallengeStatus? status,
    DateTime? expiresAt,
    DateTime? completedAt,
    DateTime? createdAt,
  }) {
    return Challenge(
      id: id ?? this.id,
      type: type ?? this.type,
      initiatorId: initiatorId ?? this.initiatorId,
      acceptorId: acceptorId ?? this.acceptorId,
      exerciseType: exerciseType ?? this.exerciseType,
      targetValue: targetValue ?? this.targetValue,
      initiatorScore: initiatorScore ?? this.initiatorScore,
      acceptorScore: acceptorScore ?? this.acceptorScore,
      winnerId: winnerId ?? this.winnerId,
      status: status ?? this.status,
      expiresAt: expiresAt ?? this.expiresAt,
      completedAt: completedAt ?? this.completedAt,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        type,
        initiatorId,
        acceptorId,
        exerciseType,
        targetValue,
        initiatorScore,
        acceptorScore,
        winnerId,
        status,
        expiresAt,
        completedAt,
        createdAt,
      ];
}

/// 排行榜条目实体
class LeaderboardEntry extends Equatable {
  /// 排名
  final int rank;

  /// 用户ID
  final String userId;

  /// 用户昵称
  final String nickname;

  /// 用户头像
  final String? avatar;

  /// 积分/分数
  final int score;

  /// 额外数据
  final Map<String, dynamic>? extra;

  const LeaderboardEntry({
    required this.rank,
    required this.userId,
    required this.nickname,
    this.avatar,
    required this.score,
    this.extra,
  });

  /// 判断是否为第一名
  bool get isFirst => rank == 1;

  /// 判断是否为第二名
  bool get isSecond => rank == 2;

  /// 判断是否为第三名
  bool get isThird => rank == 3;

  @override
  List<Object?> get props => [rank, userId, nickname, avatar, score, extra];
}
