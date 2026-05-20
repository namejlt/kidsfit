import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/exercise_record.dart';
import '../../data/datasources/training_remote_data_source.dart';

/// 训练远程数据源Provider
final trainingRemoteDataSourceProvider =
    Provider<TrainingRemoteDataSource>((ref) {
  return TrainingRemoteDataSource();
});

/// 运动状态数据类
class ExerciseState {
  /// 运动记录列表
  final List<ExerciseRecord> records;

  /// 今日训练计划
  final TrainingPlan? todayPlan;

  /// 个人最佳记录
  final List<PersonalBestDTO> personalBest;

  /// 周统计
  final WeeklyStatsDTO? weeklyStats;

  /// 月统计
  final MonthlyStatsDTO? monthlyStats;

  /// 最新体能评估
  final FitnessAssessmentDTO? latestAssessment;

  /// 是否加载中
  final bool isLoading;

  /// 是否正在运动
  final bool isExercising;

  /// 当前运动计数
  final int count;

  /// 当前运动评分
  final int score;

  /// 当前纠正提示
  final String? currentCorrection;

  /// 错误消息
  final String? error;

  const ExerciseState({
    this.records = const [],
    this.todayPlan,
    this.personalBest = const [],
    this.weeklyStats,
    this.monthlyStats,
    this.latestAssessment,
    this.isLoading = false,
    this.isExercising = false,
    this.count = 0,
    this.score = 0,
    this.currentCorrection,
    this.error,
  });

  ExerciseState copyWith({
    List<ExerciseRecord>? records,
    TrainingPlan? todayPlan,
    List<PersonalBestDTO>? personalBest,
    WeeklyStatsDTO? weeklyStats,
    MonthlyStatsDTO? monthlyStats,
    FitnessAssessmentDTO? latestAssessment,
    bool? isLoading,
    bool? isExercising,
    int? count,
    int? score,
    String? currentCorrection,
    String? error,
  }) {
    return ExerciseState(
      records: records ?? this.records,
      todayPlan: todayPlan ?? this.todayPlan,
      personalBest: personalBest ?? this.personalBest,
      weeklyStats: weeklyStats ?? this.weeklyStats,
      monthlyStats: monthlyStats ?? this.monthlyStats,
      latestAssessment: latestAssessment ?? this.latestAssessment,
      isLoading: isLoading ?? this.isLoading,
      isExercising: isExercising ?? this.isExercising,
      count: count ?? this.count,
      score: score ?? this.score,
      currentCorrection: currentCorrection,
      error: error,
    );
  }
}

/// 运动Provider
final exerciseProvider =
    StateNotifierProvider<ExerciseNotifier, ExerciseState>((ref) {
  final dataSource = ref.watch(trainingRemoteDataSourceProvider);
  return ExerciseNotifier(dataSource);
});

/// 运动Notifier
class ExerciseNotifier extends StateNotifier<ExerciseState> {
  /// 训练远程数据源
  final TrainingRemoteDataSource _dataSource;

  ExerciseNotifier(this._dataSource) : super(const ExerciseState());

  /// 加载运动记录列表
  /// 调用API获取运动记录，支持分页和类型筛选
  Future<void> loadExerciseRecords({
    int page = 1,
    int pageSize = 20,
    String? type,
  }) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _dataSource.getExerciseRecords(
        page: page,
        pageSize: pageSize,
        type: type,
      );
      if (response.isSuccess && response.data != null) {
        final records = response.data!.list
            .map((dto) => _convertExerciseDTO(dto))
            .toList();
        state = state.copyWith(records: records, isLoading: false);
      } else {
        state = state.copyWith(
          isLoading: false,
          error: response.message,
        );
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 添加运动记录
  /// 调用API创建运动记录，成功后刷新列表
  Future<bool> addExerciseRecord(CreateExerciseRequest req) async {
    try {
      final response = await _dataSource.createExerciseRecord(req);
      if (response.isSuccess && response.data != null) {
        // 将新记录添加到列表头部
        final newRecord = _convertExerciseDTO(response.data!);
        state = state.copyWith(records: [newRecord, ...state.records]);
        return true;
      }
      state = state.copyWith(error: response.message);
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载个人最佳记录
  /// 调用API获取各运动类型的最佳成绩
  Future<void> loadPersonalBest() async {
    try {
      final response = await _dataSource.getPersonalBest();
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(personalBest: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载今日训练计划
  /// 调用API获取今日训练计划，转换为本地实体
  Future<void> loadTodayPlan() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _dataSource.getTodayPlan();
      if (response.isSuccess && response.data != null) {
        final plan = _convertTrainingPlanDTO(response.data!);
        state = state.copyWith(todayPlan: plan, isLoading: false);
      } else {
        state = state.copyWith(isLoading: false, error: response.message);
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 完成训练计划
  /// 调用API标记计划为已完成
  Future<bool> completePlan(String planId) async {
    try {
      final response = await _dataSource.completePlan(planId);
      if (response.isSuccess) {
        // 更新本地计划状态
        if (state.todayPlan != null && state.todayPlan!.id == planId) {
          final updatedPlan = state.todayPlan!.copyWith(
            status: PlanStatus.completed,
            completedAt: DateTime.now(),
          );
          state = state.copyWith(todayPlan: updatedPlan);
        }
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 创建体能评估
  /// 调用API提交体能评估数据
  Future<bool> createFitnessAssessment(FitnessAssessmentDTO dto) async {
    try {
      final response = await _dataSource.createFitnessAssessment(dto);
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(latestAssessment: response.data!);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载最新体能评估
  /// 调用API获取最新的体能评估结果
  Future<void> loadLatestAssessment() async {
    try {
      final response = await _dataSource.getLatestAssessment();
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(latestAssessment: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载周统计
  /// 调用API获取本周运动统计数据
  Future<void> loadWeeklyStats() async {
    try {
      final response = await _dataSource.getWeeklyStats();
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(weeklyStats: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载月统计
  /// 调用API获取本月运动统计数据
  Future<void> loadMonthlyStats() async {
    try {
      final response = await _dataSource.getMonthlyStats();
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(monthlyStats: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 开始运动
  /// 本地状态管理，标记进入运动模式
  void startExercise(ExerciseType type) {
    state = state.copyWith(
      isExercising: true,
      count: 0,
      score: 0,
      currentCorrection: null,
    );
  }

  /// 暂停运动
  /// 本地状态管理，标记运动暂停
  void pauseExercise() {
    state = state.copyWith(isExercising: false);
  }

  /// 恢复运动
  /// 本地状态管理，标记运动恢复
  void resumeExercise() {
    state = state.copyWith(isExercising: true);
  }

  /// 停止运动
  /// 本地状态管理，标记运动结束
  void stopExercise() {
    state = state.copyWith(isExercising: false);
  }

  /// 更新运动计数
  void updateCount(int count) {
    state = state.copyWith(count: count);
  }

  /// 更新运动评分
  void updateScore(int score) {
    state = state.copyWith(score: score);
  }

  /// 设置纠正提示
  void setCorrection(String? correction) {
    state = state.copyWith(currentCorrection: correction);
  }

  /// 重置状态
  void reset() {
    state = const ExerciseState();
  }

  // ==================== 私有辅助方法 ====================

  /// 将ExerciseDTO转换为本地ExerciseRecord实体
  ExerciseRecord _convertExerciseDTO(ExerciseDTO dto) {
    return ExerciseRecord(
      id: dto.id,
      userId: '', // DTO中没有userId，由服务端关联
      type: _parseExerciseType(dto.type),
      durationSeconds: dto.durationSeconds,
      count: dto.count,
      score: dto.score,
      rhythmScore: dto.rhythmScore,
      amplitudeScore: dto.amplitudeScore,
      symmetryScore: dto.symmetryScore,
      continuityScore: dto.continuityScore,
      corrections: dto.corrections,
      isOffline: dto.isOffline,
      startedAt: dto.startedAt.isNotEmpty
          ? DateTime.parse(dto.startedAt)
          : DateTime.now(),
      completedAt: dto.completedAt != null && dto.completedAt!.isNotEmpty
          ? DateTime.parse(dto.completedAt!)
          : null,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将TrainingPlanDTO转换为本地TrainingPlan实体
  TrainingPlan _convertTrainingPlanDTO(TrainingPlanDTO dto) {
    // 按阶段分组训练项目
    final warmupItems = <ExerciseItem>[];
    final mainItems = <ExerciseItem>[];
    final cooldownItems = <ExerciseItem>[];

    for (final itemDTO in dto.items) {
      final item = _convertExerciseItemDTO(itemDTO, dto.id);
      switch (itemDTO.phase) {
        case 'warmup':
          warmupItems.add(item);
          break;
        case 'cooldown':
          cooldownItems.add(item);
          break;
        default:
          mainItems.add(item);
          break;
      }
    }

    return TrainingPlan(
      id: dto.id,
      userId: '', // DTO中没有userId
      date: dto.date.isNotEmpty ? DateTime.parse(dto.date) : DateTime.now(),
      status: dto.status == 'completed'
          ? PlanStatus.completed
          : dto.status == 'skipped'
              ? PlanStatus.skipped
              : PlanStatus.pending,
      totalDuration: dto.totalDuration,
      actualDuration: dto.actualDuration,
      warmupItems: warmupItems,
      mainItems: mainItems,
      cooldownItems: cooldownItems,
      completedAt: dto.completedAt != null && dto.completedAt!.isNotEmpty
          ? DateTime.parse(dto.completedAt!)
          : null,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将ExerciseItemDTO转换为本地ExerciseItem实体
  ExerciseItem _convertExerciseItemDTO(ExerciseItemDTO dto, String planId) {
    return ExerciseItem(
      id: dto.id,
      planId: planId,
      type: _parseExerciseType(dto.type),
      name: dto.name,
      durationSec: dto.durationSec,
      targetCount: dto.targetCount,
      difficulty: dto.difficulty,
      tips: dto.tips,
      order: dto.order,
      phase: _parseExercisePhase(dto.phase),
    );
  }

  /// 解析运动类型字符串为枚举
  ExerciseType _parseExerciseType(String value) {
    return ExerciseType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ExerciseType.jumpRope,
    );
  }

  /// 解析训练阶段字符串为枚举
  ExercisePhase _parseExercisePhase(String value) {
    return ExercisePhase.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ExercisePhase.main,
    );
  }
}
