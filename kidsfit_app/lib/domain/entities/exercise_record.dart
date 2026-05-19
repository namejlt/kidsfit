import 'package:equatable/equatable.dart';

/// 运动类型枚举
enum ExerciseType {
  /// 跳绳
  jumpRope('jump_rope', '跳绳', 'images/jump_rope.png'),

  /// 开合跳
  jumpingJack('jumping_jack', '开合跳', 'images/jumping_jack.png'),

  /// 深蹲
  squat('squat', '深蹲', 'images/squat.png'),

  /// 仰卧起坐
  sitUp('sit_up', '仰卧起坐', 'images/sit_up.png'),

  /// 高抬腿
  highKnee('high_knee', '高抬腿', 'images/high_knee.png'),

  /// 俯卧撑
  pushUp('push_up', '俯卧撑', 'images/push_up.png');

  const ExerciseType(this.value, this.displayName, this.iconPath);
  final String value;
  final String displayName;
  final String iconPath;
}

/// 训练阶段枚举
enum ExercisePhase {
  /// 热身
  warmup('warmup', '热身'),

  /// 主训练
  main('main', '主训练'),

  /// 拉伸
  cooldown('cooldown', '拉伸');

  const ExercisePhase(this.value, this.displayName);
  final String value;
  final String displayName;
}

/// 计划状态枚举
enum PlanStatus {
  /// 待完成
  pending('pending'),

  /// 已完成
  completed('completed'),

  /// 已跳过
  skipped('skipped');

  const PlanStatus(this.value);
  final String value;
}

/// 运动记录实体
/// 记录儿童每次运动的详细数据
class ExerciseRecord extends Equatable {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 运动类型
  final ExerciseType type;

  /// 运动时长（秒）
  final int durationSeconds;

  /// 动作次数
  final int count;

  /// 综合评分（0-100）
  final int score;

  /// 节奏准确性评分（0-100）
  final int rhythmScore;

  /// 幅度规范性评分（0-100）
  final int amplitudeScore;

  /// 身体对称性评分（0-100）
  final int symmetryScore;

  /// 连贯性评分（0-100）
  final int continuityScore;

  /// 纠正建议列表
  final List<String> corrections;

  /// 是否离线记录
  final bool isOffline;

  /// 开始时间
  final DateTime startedAt;

  /// 完成时间
  final DateTime? completedAt;

  /// 创建时间
  final DateTime createdAt;

  const ExerciseRecord({
    required this.id,
    required this.userId,
    required this.type,
    required this.durationSeconds,
    required this.count,
    required this.score,
    this.rhythmScore = 0,
    this.amplitudeScore = 0,
    this.symmetryScore = 0,
    this.continuityScore = 0,
    this.corrections = const [],
    this.isOffline = false,
    required this.startedAt,
    this.completedAt,
    required this.createdAt,
  });

  /// 格式化时长显示
  String get formattedDuration {
    final minutes = durationSeconds ~/ 60;
    final seconds = durationSeconds % 60;
    return '${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';
  }

  /// 获取评分等级
  String get scoreGrade {
    if (score >= 90) return 'S';
    if (score >= 80) return 'A';
    if (score >= 70) return 'B';
    if (score >= 60) return 'C';
    return 'D';
  }

  /// 获取评分描述
  String get scoreDescription {
    if (score >= 90) return '完美！';
    if (score >= 80) return '很棒！';
    if (score >= 70) return '不错！';
    if (score >= 60) return '继续加油！';
    return '再接再厉！';
  }

  ExerciseRecord copyWith({
    String? id,
    String? userId,
    ExerciseType? type,
    int? durationSeconds,
    int? count,
    int? score,
    int? rhythmScore,
    int? amplitudeScore,
    int? symmetryScore,
    int? continuityScore,
    List<String>? corrections,
    bool? isOffline,
    DateTime? startedAt,
    DateTime? completedAt,
    DateTime? createdAt,
  }) {
    return ExerciseRecord(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      type: type ?? this.type,
      durationSeconds: durationSeconds ?? this.durationSeconds,
      count: count ?? this.count,
      score: score ?? this.score,
      rhythmScore: rhythmScore ?? this.rhythmScore,
      amplitudeScore: amplitudeScore ?? this.amplitudeScore,
      symmetryScore: symmetryScore ?? this.symmetryScore,
      continuityScore: continuityScore ?? this.continuityScore,
      corrections: corrections ?? this.corrections,
      isOffline: isOffline ?? this.isOffline,
      startedAt: startedAt ?? this.startedAt,
      completedAt: completedAt ?? this.completedAt,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        type,
        durationSeconds,
        count,
        score,
        rhythmScore,
        amplitudeScore,
        symmetryScore,
        continuityScore,
        corrections,
        isOffline,
        startedAt,
        completedAt,
        createdAt,
      ];
}

/// 训练项目实体
class ExerciseItem extends Equatable {
  /// 项目ID
  final String id;

  /// 所属计划ID
  final String planId;

  /// 运动类型
  final ExerciseType type;

  /// 项目名称
  final String name;

  /// 目标时长（秒）
  final int? durationSec;

  /// 目标次数
  final int? targetCount;

  /// 难度等级（1-5）
  final int difficulty;

  /// 要点提示
  final String? tips;

  /// 排序顺序
  final int order;

  /// 训练阶段
  final ExercisePhase phase;

  const ExerciseItem({
    required this.id,
    required this.planId,
    required this.type,
    required this.name,
    this.durationSec,
    this.targetCount,
    this.difficulty = 1,
    this.tips,
    required this.order,
    required this.phase,
  });

  @override
  List<Object?> get props => [
        id,
        planId,
        type,
        name,
        durationSec,
        targetCount,
        difficulty,
        tips,
        order,
        phase,
      ];
}

/// 训练计划实体
class TrainingPlan extends Equatable {
  /// 计划ID
  final String id;

  /// 用户ID
  final String userId;

  /// 计划日期
  final DateTime date;

  /// 计划状态
  final PlanStatus status;

  /// 计划总时长（分钟）
  final int totalDuration;

  /// 实际完成时长（分钟）
  final int? actualDuration;

  /// 热身项目列表
  final List<ExerciseItem> warmupItems;

  /// 主训练项目列表
  final List<ExerciseItem> mainItems;

  /// 拉伸项目列表
  final List<ExerciseItem> cooldownItems;

  /// 完成时间
  final DateTime? completedAt;

  /// 创建时间
  final DateTime createdAt;

  const TrainingPlan({
    required this.id,
    required this.userId,
    required this.date,
    this.status = PlanStatus.pending,
    required this.totalDuration,
    this.actualDuration,
    this.warmupItems = const [],
    this.mainItems = const [],
    this.cooldownItems = const [],
    this.completedAt,
    required this.createdAt,
  });

  /// 获取所有训练项目
  List<ExerciseItem> get allItems => [
        ...warmupItems,
        ...mainItems,
        ...cooldownItems,
      ];

  /// 获取已完成项目数量
  int get completedItemCount {
    return allItems.where((item) => item.phase == ExercisePhase.main).length;
  }

  /// 获取总项目数量
  int get totalItemCount => mainItems.length;

  /// 获取完成进度
  double get progress {
    if (totalItemCount == 0) return 0;
    return completedItemCount / totalItemCount;
  }

  /// 是否已完成
  bool get isCompleted => status == PlanStatus.completed;

  /// 是否可跳过
  bool get canSkip => status == PlanStatus.pending;

  TrainingPlan copyWith({
    String? id,
    String? userId,
    DateTime? date,
    PlanStatus? status,
    int? totalDuration,
    int? actualDuration,
    List<ExerciseItem>? warmupItems,
    List<ExerciseItem>? mainItems,
    List<ExerciseItem>? cooldownItems,
    DateTime? completedAt,
    DateTime? createdAt,
  }) {
    return TrainingPlan(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      date: date ?? this.date,
      status: status ?? this.status,
      totalDuration: totalDuration ?? this.totalDuration,
      actualDuration: actualDuration ?? this.actualDuration,
      warmupItems: warmupItems ?? this.warmupItems,
      mainItems: mainItems ?? this.mainItems,
      cooldownItems: cooldownItems ?? this.cooldownItems,
      completedAt: completedAt ?? this.completedAt,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        date,
        status,
        totalDuration,
        actualDuration,
        warmupItems,
        mainItems,
        cooldownItems,
        completedAt,
        createdAt,
      ];
}

/// 体能评估实体
class FitnessAssessment extends Equatable {
  /// 评估ID
  final String id;

  /// 用户ID
  final String userId;

  /// 耐力评分（1-10）
  final int endurance;

  /// 灵敏评分（1-10）
  final int agility;

  /// 力量评分（1-10）
  final int strength;

  /// 速度评分（1-10）
  final int speed;

  /// 协调评分（1-10）
  final int coordination;

  /// 平衡评分（1-10）
  final int balance;

  /// 柔韧评分（1-10）
  final int flexibility;

  /// 评估时间
  final DateTime assessedAt;

  /// 创建时间
  final DateTime createdAt;

  const FitnessAssessment({
    required this.id,
    required this.userId,
    this.endurance = 5,
    this.agility = 5,
    this.strength = 5,
    this.speed = 5,
    this.coordination = 5,
    this.balance = 5,
    this.flexibility = 5,
    required this.assessedAt,
    required this.createdAt,
  });

  /// 获取所有维度评分
  Map<String, int> get allScores => {
        '耐力': endurance,
        '灵敏': agility,
        '力量': strength,
        '速度': speed,
        '协调': coordination,
        '平衡': balance,
        '柔韧': flexibility,
      };

  /// 获取平均评分
  double get averageScore {
    return allScores.values.reduce((a, b) => a + b) / allScores.length;
  }

  /// 获取总分
  int get totalScore => allScores.values.reduce((a, b) => a + b);

  @override
  List<Object?> get props => [
        id,
        userId,
        endurance,
        agility,
        strength,
        speed,
        coordination,
        balance,
        flexibility,
        assessedAt,
        createdAt,
      ];
}
