import '../../core/network/api_client.dart';
import '../../core/network/api_response.dart';
import '../../core/constants/app_constants.dart';

// ==================== DTO类定义 ====================

/// 创建运动记录请求DTO
class CreateExerciseRequest {
  /// 运动类型
  final String type;

  /// 运动时长（秒）
  final int durationSeconds;

  /// 动作次数
  final int count;

  /// 综合评分
  final int score;

  /// 节奏评分
  final int rhythmScore;

  /// 幅度评分
  final int amplitudeScore;

  /// 对称性评分
  final int symmetryScore;

  /// 连贯性评分
  final int continuityScore;

  /// 纠正建议
  final List<String> corrections;

  /// 是否离线记录
  final bool isOffline;

  /// 开始时间
  final String startedAt;

  /// 完成时间
  final String? completedAt;

  const CreateExerciseRequest({
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
  });

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'type': type,
        'duration_seconds': durationSeconds,
        'count': count,
        'score': score,
        'rhythm_score': rhythmScore,
        'amplitude_score': amplitudeScore,
        'symmetry_score': symmetryScore,
        'continuity_score': continuityScore,
        'corrections': corrections,
        'is_offline': isOffline,
        'started_at': startedAt,
        'completed_at': completedAt,
      };
}

/// 运动记录DTO
class ExerciseDTO {
  /// 记录ID
  final String id;

  /// 运动类型
  final String type;

  /// 运动时长（秒）
  final int durationSeconds;

  /// 动作次数
  final int count;

  /// 综合评分
  final int score;

  /// 节奏评分
  final int rhythmScore;

  /// 幅度评分
  final int amplitudeScore;

  /// 对称性评分
  final int symmetryScore;

  /// 连贯性评分
  final int continuityScore;

  /// 纠正建议
  final List<String> corrections;

  /// 是否离线记录
  final bool isOffline;

  /// 开始时间
  final String startedAt;

  /// 完成时间
  final String? completedAt;

  /// 获得积分
  final int? pointsEarned;

  /// 获得徽章
  final List<BadgeEarnedDTO> badgesEarned;

  /// 是否打破记录
  final bool isRecordBroken;

  /// 创建时间
  final String createdAt;

  const ExerciseDTO({
    required this.id,
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
    this.pointsEarned,
    this.badgesEarned = const [],
    this.isRecordBroken = false,
    required this.createdAt,
  });

  /// 从JSON创建ExerciseDTO
  factory ExerciseDTO.fromJson(Map<String, dynamic> json) {
    final badgesRaw = json['badges_earned'] as List<dynamic>? ?? [];
    return ExerciseDTO(
      id: json['id'] as String? ?? '',
      type: json['type'] as String? ?? '',
      durationSeconds: json['duration_seconds'] as int? ?? 0,
      count: json['count'] as int? ?? 0,
      score: json['score'] as int? ?? 0,
      rhythmScore: json['rhythm_score'] as int? ?? 0,
      amplitudeScore: json['amplitude_score'] as int? ?? 0,
      symmetryScore: json['symmetry_score'] as int? ?? 0,
      continuityScore: json['continuity_score'] as int? ?? 0,
      corrections: (json['corrections'] as List<dynamic>? ?? [])
          .map((e) => e as String)
          .toList(),
      isOffline: json['is_offline'] as bool? ?? false,
      startedAt: json['started_at'] as String? ?? '',
      completedAt: json['completed_at'] as String?,
      pointsEarned: json['points_earned'] as int?,
      badgesEarned: badgesRaw
          .map((e) => BadgeEarnedDTO.fromJson(e as Map<String, dynamic>))
          .toList(),
      isRecordBroken: json['is_record_broken'] as bool? ?? false,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type,
        'duration_seconds': durationSeconds,
        'count': count,
        'score': score,
        'rhythm_score': rhythmScore,
        'amplitude_score': amplitudeScore,
        'symmetry_score': symmetryScore,
        'continuity_score': continuityScore,
        'corrections': corrections,
        'is_offline': isOffline,
        'started_at': startedAt,
        'completed_at': completedAt,
        'points_earned': pointsEarned,
        'badges_earned': badgesEarned.map((e) => e.toJson()).toList(),
        'is_record_broken': isRecordBroken,
        'created_at': createdAt,
      };
}

/// 训练项目DTO
class ExerciseItemDTO {
  /// 项目ID
  final String id;

  /// 运动类型
  final String type;

  /// 项目名称
  final String name;

  /// 目标时长（秒）
  final int? durationSec;

  /// 目标次数
  final int? targetCount;

  /// 难度等级
  final int difficulty;

  /// 要点提示
  final String? tips;

  /// 排序顺序
  final int order;

  /// 训练阶段
  final String phase;

  const ExerciseItemDTO({
    required this.id,
    required this.type,
    required this.name,
    this.durationSec,
    this.targetCount,
    this.difficulty = 1,
    this.tips,
    required this.order,
    required this.phase,
  });

  /// 从JSON创建ExerciseItemDTO
  factory ExerciseItemDTO.fromJson(Map<String, dynamic> json) {
    return ExerciseItemDTO(
      id: json['id'] as String? ?? '',
      type: json['type'] as String? ?? '',
      name: json['name'] as String? ?? '',
      durationSec: json['duration_sec'] as int?,
      targetCount: json['target_count'] as int?,
      difficulty: json['difficulty'] as int? ?? 1,
      tips: json['tips'] as String?,
      order: json['order'] as int? ?? 0,
      phase: json['phase'] as String? ?? 'main',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type,
        'name': name,
        'duration_sec': durationSec,
        'target_count': targetCount,
        'difficulty': difficulty,
        'tips': tips,
        'order': order,
        'phase': phase,
      };
}

/// 训练计划DTO
class TrainingPlanDTO {
  /// 计划ID
  final String id;

  /// 计划日期
  final String date;

  /// 计划状态
  final String status;

  /// 计划总时长（秒）
  final int totalDuration;

  /// 实际完成时长（秒）
  final int? actualDuration;

  /// 训练项目列表
  final List<ExerciseItemDTO> items;

  /// 完成时间
  final String? completedAt;

  /// 创建时间
  final String createdAt;

  const TrainingPlanDTO({
    required this.id,
    required this.date,
    this.status = 'pending',
    required this.totalDuration,
    this.actualDuration,
    this.items = const [],
    this.completedAt,
    required this.createdAt,
  });

  /// 从JSON创建TrainingPlanDTO
  factory TrainingPlanDTO.fromJson(Map<String, dynamic> json) {
    final itemsRaw = json['items'] as List<dynamic>? ?? [];
    return TrainingPlanDTO(
      id: json['id'] as String? ?? '',
      date: json['date'] as String? ?? '',
      status: json['status'] as String? ?? 'pending',
      totalDuration: json['total_duration'] as int? ?? 0,
      actualDuration: json['actual_duration'] as int?,
      items: itemsRaw
          .map((e) => ExerciseItemDTO.fromJson(e as Map<String, dynamic>))
          .toList(),
      completedAt: json['completed_at'] as String?,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'date': date,
        'status': status,
        'total_duration': totalDuration,
        'actual_duration': actualDuration,
        'items': items.map((e) => e.toJson()).toList(),
        'completed_at': completedAt,
        'created_at': createdAt,
      };
}

/// 个人最佳记录DTO
class PersonalBestDTO {
  /// 运动类型
  final String type;

  /// 最佳评分
  final int bestScore;

  /// 最佳次数
  final int bestCount;

  /// 达成时间
  final String achievedAt;

  const PersonalBestDTO({
    required this.type,
    required this.bestScore,
    required this.bestCount,
    required this.achievedAt,
  });

  /// 从JSON创建PersonalBestDTO
  factory PersonalBestDTO.fromJson(Map<String, dynamic> json) {
    return PersonalBestDTO(
      type: json['type'] as String? ?? '',
      bestScore: json['best_score'] as int? ?? 0,
      bestCount: json['best_count'] as int? ?? 0,
      achievedAt: json['achieved_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'type': type,
        'best_score': bestScore,
        'best_count': bestCount,
        'achieved_at': achievedAt,
      };
}

/// 体能评估DTO
class FitnessAssessmentDTO {
  /// 评估ID
  final String id;

  /// 耐力评分
  final int endurance;

  /// 灵敏评分
  final int agility;

  /// 力量评分
  final int strength;

  /// 速度评分
  final int speed;

  /// 协调评分
  final int coordination;

  /// 平衡评分
  final int balance;

  /// 柔韧评分
  final int flexibility;

  /// 评估时间
  final String assessedAt;

  const FitnessAssessmentDTO({
    required this.id,
    this.endurance = 5,
    this.agility = 5,
    this.strength = 5,
    this.speed = 5,
    this.coordination = 5,
    this.balance = 5,
    this.flexibility = 5,
    required this.assessedAt,
  });

  /// 从JSON创建FitnessAssessmentDTO
  factory FitnessAssessmentDTO.fromJson(Map<String, dynamic> json) {
    return FitnessAssessmentDTO(
      id: json['id'] as String? ?? '',
      endurance: json['endurance'] as int? ?? 5,
      agility: json['agility'] as int? ?? 5,
      strength: json['strength'] as int? ?? 5,
      speed: json['speed'] as int? ?? 5,
      coordination: json['coordination'] as int? ?? 5,
      balance: json['balance'] as int? ?? 5,
      flexibility: json['flexibility'] as int? ?? 5,
      assessedAt: json['assessed_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'endurance': endurance,
        'agility': agility,
        'strength': strength,
        'speed': speed,
        'coordination': coordination,
        'balance': balance,
        'flexibility': flexibility,
        'assessed_at': assessedAt,
      };
}

/// 周统计DTO
class WeeklyStatsDTO {
  /// 总运动次数
  final int totalExercises;

  /// 总运动时长（秒）
  final int totalDuration;

  /// 平均评分
  final double averageScore;

  /// 运动类型分布
  final Map<String, int> exerciseTypeBreakdown;

  const WeeklyStatsDTO({
    this.totalExercises = 0,
    this.totalDuration = 0,
    this.averageScore = 0,
    this.exerciseTypeBreakdown = const {},
  });

  /// 从JSON创建WeeklyStatsDTO
  factory WeeklyStatsDTO.fromJson(Map<String, dynamic> json) {
    final breakdownRaw = json['exercise_type_breakdown'] as Map<String, dynamic>? ?? {};
    return WeeklyStatsDTO(
      totalExercises: json['total_exercises'] as int? ?? 0,
      totalDuration: json['total_duration'] as int? ?? 0,
      averageScore: (json['average_score'] as num?)?.toDouble() ?? 0,
      exerciseTypeBreakdown: breakdownRaw.map((k, v) => MapEntry(k, v as int)),
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'total_exercises': totalExercises,
        'total_duration': totalDuration,
        'average_score': averageScore,
        'exercise_type_breakdown': exerciseTypeBreakdown,
      };
}

/// 月统计DTO
class MonthlyStatsDTO {
  /// 总运动次数
  final int totalExercises;

  /// 总运动时长（秒）
  final int totalDuration;

  /// 平均评分
  final double averageScore;

  /// 运动类型分布
  final Map<String, int> exerciseTypeBreakdown;

  const MonthlyStatsDTO({
    this.totalExercises = 0,
    this.totalDuration = 0,
    this.averageScore = 0,
    this.exerciseTypeBreakdown = const {},
  });

  /// 从JSON创建MonthlyStatsDTO
  factory MonthlyStatsDTO.fromJson(Map<String, dynamic> json) {
    final breakdownRaw = json['exercise_type_breakdown'] as Map<String, dynamic>? ?? {};
    return MonthlyStatsDTO(
      totalExercises: json['total_exercises'] as int? ?? 0,
      totalDuration: json['total_duration'] as int? ?? 0,
      averageScore: (json['average_score'] as num?)?.toDouble() ?? 0,
      exerciseTypeBreakdown: breakdownRaw.map((k, v) => MapEntry(k, v as int)),
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'total_exercises': totalExercises,
        'total_duration': totalDuration,
        'average_score': averageScore,
        'exercise_type_breakdown': exerciseTypeBreakdown,
      };
}

/// 徽章获得DTO
class BadgeEarnedDTO {
  /// 徽章ID
  final String id;

  /// 徽章名称
  final String name;

  /// 徽章图标
  final String icon;

  const BadgeEarnedDTO({
    required this.id,
    required this.name,
    required this.icon,
  });

  /// 从JSON创建BadgeEarnedDTO
  factory BadgeEarnedDTO.fromJson(Map<String, dynamic> json) {
    return BadgeEarnedDTO(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      icon: json['icon'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'icon': icon,
      };
}

// ==================== 训练远程数据源 ====================

/// 训练远程数据源
/// 负责与后端训练相关API通信
class TrainingRemoteDataSource {
  /// API客户端实例
  final ApiClient _apiClient;

  TrainingRemoteDataSource({ApiClient? apiClient})
      : _apiClient = apiClient ?? ApiClient.instance;

  /// 创建运动记录
  /// POST /api/v1/training/exercises
  Future<ApiResponse<ExerciseDTO>> createExerciseRecord(
      CreateExerciseRequest req) async {
    return _apiClient.post(
      AppConstants.exerciseRecords,
      data: req.toJson(),
      fromJsonT: (data) => ExerciseDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取运动记录列表（分页）
  /// GET /api/v1/training/exercises
  Future<ApiResponse<PaginatedData<ExerciseDTO>>> getExerciseRecords({
    int page = 1,
    int pageSize = 20,
    String? type,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (type != null) {
      queryParams['type'] = type;
    }

    return _apiClient.getPaged(
      AppConstants.exerciseRecords,
      queryParameters: queryParams,
      fromJsonT: (data) => ExerciseDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取个人最佳记录
  /// GET /api/v1/training/personal-best
  Future<ApiResponse<List<PersonalBestDTO>>> getPersonalBest() async {
    return _apiClient.get(
      AppConstants.personalBest,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => PersonalBestDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 获取今日训练计划
  /// GET /api/v1/training/plan/today
  Future<ApiResponse<TrainingPlanDTO>> getTodayPlan() async {
    return _apiClient.get(
      AppConstants.todayPlan,
      fromJsonT: (data) =>
          TrainingPlanDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 完成训练计划
  /// PUT /api/v1/training/plan/{planId}/complete
  Future<ApiResponse<void>> completePlan(String planId) async {
    return _apiClient.putVoid(
      '${AppConstants.completePlan}/$planId/complete',
    );
  }

  /// 创建体能评估
  /// POST /api/v1/training/assessments
  Future<ApiResponse<FitnessAssessmentDTO>> createFitnessAssessment(
      FitnessAssessmentDTO dto) async {
    return _apiClient.post(
      AppConstants.fitnessAssessment,
      data: dto.toJson(),
      fromJsonT: (data) =>
          FitnessAssessmentDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取最新体能评估
  /// GET /api/v1/training/assessments/latest
  Future<ApiResponse<FitnessAssessmentDTO>> getLatestAssessment() async {
    return _apiClient.get(
      AppConstants.latestAssessment,
      fromJsonT: (data) =>
          FitnessAssessmentDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取周统计
  /// GET /api/v1/training/stats/weekly
  Future<ApiResponse<WeeklyStatsDTO>> getWeeklyStats() async {
    return _apiClient.get(
      AppConstants.weeklyStats,
      fromJsonT: (data) =>
          WeeklyStatsDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取月统计
  /// GET /api/v1/training/stats/monthly
  Future<ApiResponse<MonthlyStatsDTO>> getMonthlyStats() async {
    return _apiClient.get(
      AppConstants.monthlyStats,
      fromJsonT: (data) =>
          MonthlyStatsDTO.fromJson(data as Map<String, dynamic>),
    );
  }
}
