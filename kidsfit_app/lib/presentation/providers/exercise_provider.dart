import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/exercise_record.dart';
import '../../services/storage_service.dart';

/// 今日训练计划Provider
final todayPlanProvider = StateProvider<TrainingPlan?>((ref) => null);

/// 运动记录列表Provider
final exerciseRecordsProvider =
    StateNotifierProvider<ExerciseRecordsNotifier, List<ExerciseRecord>>(
        (ref) {
  return ExerciseRecordsNotifier();
});

/// 运动记录Notifier
class ExerciseRecordsNotifier
    extends StateNotifier<List<ExerciseRecord>> {
  ExerciseRecordsNotifier() : super([]);

  /// 加载运动记录
  Future<void> loadRecords(String userId) async {
    try {
      final recordsData = await StorageService.getExerciseRecordsByUser(userId);
      state = recordsData.map((data) => _parseRecord(data)).toList();
    } catch (e) {
      state = [];
    }
  }

  /// 添加运动记录
  Future<void> addRecord(ExerciseRecord record) async {
    await StorageService.saveExerciseRecord(record.id, _toMap(record));
    state = [...state, record];
  }

  /// 解析记录
  ExerciseRecord _parseRecord(Map<String, dynamic> data) {
    return ExerciseRecord(
      id: data['id'] as String,
      userId: data['user_id'] as String,
      type: ExerciseType.values.firstWhere(
        (e) => e.value == data['type'],
        orElse: () => ExerciseType.jumpRope,
      ),
      durationSeconds: data['duration_seconds'] as int,
      count: data['count'] as int,
      score: data['score'] as int,
      rhythmScore: data['rhythm_score'] as int? ?? 0,
      amplitudeScore: data['amplitude_score'] as int? ?? 0,
      symmetryScore: data['symmetry_score'] as int? ?? 0,
      continuityScore: data['continuity_score'] as int? ?? 0,
      corrections: (data['corrections'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
      isOffline: data['is_offline'] as bool? ?? false,
      startedAt: DateTime.parse(data['started_at'] as String),
      completedAt: data['completed_at'] != null
          ? DateTime.parse(data['completed_at'] as String)
          : null,
      createdAt: DateTime.parse(data['created_at'] as String),
    );
  }

  /// 转换为Map
  Map<String, dynamic> _toMap(ExerciseRecord record) {
    return {
      'id': record.id,
      'user_id': record.userId,
      'type': record.type.value,
      'duration_seconds': record.durationSeconds,
      'count': record.count,
      'score': record.score,
      'rhythm_score': record.rhythmScore,
      'amplitude_score': record.amplitudeScore,
      'symmetry_score': record.symmetryScore,
      'continuity_score': record.continuityScore,
      'corrections': record.corrections,
      'is_offline': record.isOffline,
      'started_at': record.startedAt.toIso8601String(),
      'completed_at': record.completedAt?.toIso8601String(),
      'created_at': record.createdAt.toIso8601String(),
    };
  }
}

/// 运动Provider
final exerciseProvider =
    StateNotifierProvider<ExerciseNotifier, ExerciseState>((ref) {
  return ExerciseNotifier();
});

/// 运动状态
class ExerciseState {
  final bool isLoading;
  final bool isExercising;
  final int count;
  final int score;
  final String? currentCorrection;
  final String? error;

  const ExerciseState({
    this.isLoading = false,
    this.isExercising = false,
    this.count = 0,
    this.score = 0,
    this.currentCorrection,
    this.error,
  });

  ExerciseState copyWith({
    bool? isLoading,
    bool? isExercising,
    int? count,
    int? score,
    String? currentCorrection,
    String? error,
  }) {
    return ExerciseState(
      isLoading: isLoading ?? this.isLoading,
      isExercising: isExercising ?? this.isExercising,
      count: count ?? this.count,
      score: score ?? this.score,
      currentCorrection: currentCorrection,
      error: error,
    );
  }
}

/// 运动Notifier
class ExerciseNotifier extends StateNotifier<ExerciseState> {
  ExerciseNotifier() : super(const ExerciseState());

  /// 加载今日计划
  Future<void> loadTodayPlan() async {
    state = state.copyWith(isLoading: true);

    try {
      // TODO: 调用API获取今日训练计划
      await Future.delayed(const Duration(milliseconds: 500));

      // 模拟数据
      final plan = TrainingPlan(
        id: 'plan_${DateTime.now().millisecondsSinceEpoch}',
        userId: 'current_user',
        date: DateTime.now(),
        totalDuration: 20,
        mainItems: [
          ExerciseItem(
            id: 'item_1',
            planId: 'plan_1',
            type: ExerciseType.jumpRope,
            name: '跳绳练习',
            targetCount: 50,
            difficulty: 2,
            tips: '保持节奏，手腕摇绳',
            order: 1,
            phase: ExercisePhase.main,
          ),
          ExerciseItem(
            id: 'item_2',
            planId: 'plan_1',
            type: ExerciseType.squat,
            name: '深蹲训练',
            targetCount: 20,
            difficulty: 2,
            tips: '膝盖不要内扣',
            order: 2,
            phase: ExercisePhase.main,
          ),
        ],
        createdAt: DateTime.now(),
      );

      ref.read(todayPlanProvider.notifier).state = plan;
    } catch (e) {
      state = state.copyWith(error: e.toString());
    } finally {
      state = state.copyWith(isLoading: false);
    }
  }

  /// 开始运动
  void startExercise() {
    state = state.copyWith(isExercising: true, count: 0, score: 0);
  }

  /// 更新计数
  void updateCount(int count) {
    state = state.copyWith(count: count);
  }

  /// 更新评分
  void updateScore(int score) {
    state = state.copyWith(score: score);
  }

  /// 设置纠正提示
  void setCorrection(String? correction) {
    state = state.copyWith(currentCorrection: correction);
  }

  /// 暂停运动
  void pauseExercise() {
    state = state.copyWith(isExercising: false);
  }

  /// 恢复运动
  void resumeExercise() {
    state = state.copyWith(isExercising: true);
  }

  /// 停止运动
  void stopExercise() {
    state = state.copyWith(isExercising: false);
  }

  /// 重置状态
  void reset() {
    state = const ExerciseState();
  }
}
