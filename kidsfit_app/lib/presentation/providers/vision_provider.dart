import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/vision_record.dart';
import '../../services/storage_service.dart';

/// 今日户外活动Provider
final todayOutdoorProvider = StateProvider<OutdoorActivity?>((ref) => null);

/// 视力记录列表Provider
final visionRecordsProvider =
    StateNotifierProvider<VisionRecordsNotifier, List<VisionRecord>>((ref) {
  return VisionRecordsNotifier();
});

/// 视力记录Notifier
class VisionRecordsNotifier extends StateNotifier<List<VisionRecord>> {
  VisionRecordsNotifier() : super([]);

  /// 加载视力记录
  Future<void> loadRecords(String childId) async {
    try {
      final recordsData = await StorageService.getVisionRecordsByChild(childId);
      state = recordsData.map((data) => _parseRecord(data)).toList();
    } catch (e) {
      state = [];
    }
  }

  /// 添加视力记录
  Future<void> addRecord(VisionRecord record) async {
    await StorageService.saveVisionRecord(record.id, _toMap(record));
    state = [record, ...state];
  }

  /// 解析记录
  VisionRecord _parseRecord(Map<String, dynamic> data) {
    return VisionRecord(
      id: data['id'] as String,
      userId: data['user_id'] as String,
      childId: data['child_id'] as String,
      date: DateTime.parse(data['date'] as String),
      rightEye: EyeData(
        sph: (data['right_sph'] as num?)?.toDouble() ?? 0,
        cyl: (data['right_cyl'] as num?)?.toDouble() ?? 0,
        axis: data['right_axis'] as int? ?? 0,
        va: (data['right_va'] as num?)?.toDouble() ?? 1.0,
      ),
      leftEye: EyeData(
        sph: (data['left_sph'] as num?)?.toDouble() ?? 0,
        cyl: (data['left_cyl'] as num?)?.toDouble() ?? 0,
        axis: data['left_axis'] as int? ?? 0,
        va: (data['left_va'] as num?)?.toDouble() ?? 1.0,
      ),
      source: data['source'] == 'ocr'
          ? VisionDataSource.ocr
          : VisionDataSource.manual,
      createdAt: DateTime.parse(data['created_at'] as String),
    );
  }

  /// 转换为Map
  Map<String, dynamic> _toMap(VisionRecord record) {
    return {
      'id': record.id,
      'user_id': record.userId,
      'child_id': record.childId,
      'date': record.date.toIso8601String(),
      'right_sph': record.rightEye.sph,
      'right_cyl': record.rightEye.cyl,
      'right_axis': record.rightEye.axis,
      'right_va': record.rightEye.va,
      'left_sph': record.leftEye.sph,
      'left_cyl': record.leftEye.cyl,
      'left_axis': record.leftEye.axis,
      'left_va': record.leftEye.va,
      'source': record.source.value,
      'created_at': record.createdAt.toIso8601String(),
    };
  }
}

/// 视力Provider
final visionProvider =
    StateNotifierProvider<VisionNotifier, VisionState>((ref) {
  return VisionNotifier(ref);
});

/// 视力状态
class VisionState {
  final bool isLoading;
  final OutdoorActivity? todayOutdoor;
  final List<VisionRecord> visionRecords;
  final String? error;

  const VisionState({
    this.isLoading = false,
    this.todayOutdoor,
    this.visionRecords = const [],
    this.error,
  });

  VisionState copyWith({
    bool? isLoading,
    OutdoorActivity? todayOutdoor,
    List<VisionRecord>? visionRecords,
    String? error,
  }) {
    return VisionState(
      isLoading: isLoading ?? this.isLoading,
      todayOutdoor: todayOutdoor ?? this.todayOutdoor,
      visionRecords: visionRecords ?? this.visionRecords,
      error: error,
    );
  }
}

/// 视力Notifier
class VisionNotifier extends StateNotifier<VisionState> {
  final Ref _ref;

  VisionNotifier(this._ref) : super(const VisionState());

  /// 加载今日户外活动
  Future<void> loadTodayOutdoor() async {
    state = state.copyWith(isLoading: true);

    try {
      // TODO: 调用API获取今日户外活动
      await Future.delayed(const Duration(milliseconds: 300));

      // 模拟数据
      final outdoor = OutdoorActivity(
        id: 'outdoor_${DateTime.now().millisecondsSinceEpoch}',
        userId: 'current_user',
        date: DateTime.now(),
        durationMin: 45,
        createdAt: DateTime.now(),
      );

      state = state.copyWith(
        isLoading: false,
        todayOutdoor: outdoor,
      );

      // 更新Provider
      _ref.read(todayOutdoorProvider.notifier).state = outdoor;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 加载视力记录
  Future<void> loadVisionRecords(String childId) async {
    state = state.copyWith(isLoading: true);

    try {
      // TODO: 调用API获取视力记录
      await Future.delayed(const Duration(milliseconds: 300));

      // 模拟数据
      final records = [
        VisionRecord(
          id: 'vision_1',
          userId: 'user_1',
          childId: childId,
          date: DateTime.now(),
          rightEye: const EyeData(sph: 0.5, cyl: -0.5, axis: 180, va: 1.0),
          leftEye: const EyeData(sph: 0.75, cyl: -0.25, axis: 90, va: 1.0),
          source: VisionDataSource.manual,
          createdAt: DateTime.now(),
        ),
      ];

      state = state.copyWith(
        isLoading: false,
        visionRecords: records,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 添加视力记录
  Future<void> addVisionRecord({
    required String childId,
    required EyeData rightEye,
    required EyeData leftEye,
    VisionDataSource source = VisionDataSource.manual,
  }) async {
    final record = VisionRecord(
      id: 'vision_${DateTime.now().millisecondsSinceEpoch}',
      userId: 'current_user',
      childId: childId,
      date: DateTime.now(),
      rightEye: rightEye,
      leftEye: leftEye,
      source: source,
      createdAt: DateTime.now(),
    );

    await _ref.read(visionRecordsProvider.notifier).addRecord(record);

    // 重新加载记录
    await loadVisionRecords(childId);
  }

  /// OCR识别验光单
  Future<VisionRecord?> recognizePrescription(String imagePath) async {
    // TODO: 调用OCR服务识别验光单
    return null;
  }
}

/// 儿童列表Provider
final childrenListProvider = StateProvider<List>((ref) => []);
