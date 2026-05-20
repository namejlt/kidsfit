import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/vision_record.dart';
import '../../data/datasources/vision_remote_data_source.dart';
import '../../core/network/api_client.dart';

/// 视力远程数据源Provider
final visionRemoteDataSourceProvider =
    Provider<VisionRemoteDataSource>((ref) {
  return VisionRemoteDataSource();
});

/// 视力状态数据类
class VisionState {
  /// 视力记录列表
  final List<VisionRecord> records;

  /// 今日户外活动
  final OutdoorActivity? todayOutdoor;

  /// 视力趋势
  final VisionTrendDTO? trend;

  /// 用眼提醒列表
  final List<EyeReminder> reminders;

  /// 是否加载中
  final bool isLoading;

  /// 错误消息
  final String? error;

  const VisionState({
    this.records = const [],
    this.todayOutdoor,
    this.trend,
    this.reminders = const [],
    this.isLoading = false,
    this.error,
  });

  VisionState copyWith({
    List<VisionRecord>? records,
    OutdoorActivity? todayOutdoor,
    VisionTrendDTO? trend,
    List<EyeReminder>? reminders,
    bool? isLoading,
    String? error,
  }) {
    return VisionState(
      records: records ?? this.records,
      todayOutdoor: todayOutdoor ?? this.todayOutdoor,
      trend: trend ?? this.trend,
      reminders: reminders ?? this.reminders,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// 视力Provider
final visionProvider =
    StateNotifierProvider<VisionNotifier, VisionState>((ref) {
  final dataSource = ref.watch(visionRemoteDataSourceProvider);
  return VisionNotifier(dataSource);
});

/// 视力Notifier
class VisionNotifier extends StateNotifier<VisionState> {
  /// 视力远程数据源
  final VisionRemoteDataSource _dataSource;

  VisionNotifier(this._dataSource) : super(const VisionState());

  /// 加载视力记录列表
  /// 调用API获取视力记录，支持按儿童ID筛选和分页
  Future<void> loadVisionRecords({
    String? childId,
    int page = 1,
    int pageSize = 20,
  }) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _dataSource.getVisionRecords(
        childId: childId,
        page: page,
        pageSize: pageSize,
      );
      if (response.isSuccess && response.data != null) {
        final records = response.data!.list
            .map((dto) => _convertVisionRecordDTO(dto))
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

  /// 添加视力记录
  /// 调用API创建视力记录，成功后刷新列表
  Future<bool> addVisionRecord(CreateVisionRecordRequest req) async {
    try {
      final response = await _dataSource.createVisionRecord(req);
      if (response.isSuccess && response.data != null) {
        final newRecord = _convertVisionRecordDTO(response.data!);
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

  /// 加载视力趋势
  /// 调用API获取视力变化趋势数据
  Future<void> loadVisionTrend({String? childId}) async {
    try {
      final response = await _dataSource.getVisionTrend(childId: childId);
      if (response.isSuccess && response.data != null) {
        state = state.copyWith(trend: response.data!);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载今日户外活动
  /// 调用API获取今日户外运动时长数据
  Future<void> loadTodayOutdoor() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _dataSource.getTodayOutdoor();
      if (response.isSuccess && response.data != null) {
        final outdoor = _convertOutdoorActivityDTO(response.data!);
        state = state.copyWith(
          isLoading: false,
          todayOutdoor: outdoor,
        );
      } else {
        state = state.copyWith(isLoading: false, error: response.message);
      }
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  /// 同步户外活动数据
  /// 调用API上报户外运动时长
  Future<bool> syncOutdoorData(int durationMin) async {
    try {
      final response = await _dataSource.syncOutdoorData(durationMin);
      if (response.isSuccess && response.data != null) {
        final outdoor = _convertOutdoorActivityDTO(response.data!);
        state = state.copyWith(todayOutdoor: outdoor);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载用眼提醒列表
  /// 调用API获取用眼提醒，支持分页
  Future<void> loadReminders({
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dataSource.getReminders(
        page: page,
        pageSize: pageSize,
      );
      if (response.isSuccess && response.data != null) {
        final reminders = response.data!.list
            .map((dto) => _convertEyeReminderDTO(dto))
            .toList();
        state = state.copyWith(reminders: reminders);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 确认提醒
  /// 调用API标记提醒为已确认
  Future<bool> ackReminder(String reminderId) async {
    try {
      final response = await _dataSource.ackReminder(reminderId);
      if (response.isSuccess) {
        // 更新本地提醒状态
        final updatedReminders = state.reminders.map((r) {
          if (r.id == reminderId) {
            return r.copyWith(acknowledged: true);
          }
          return r;
        }).toList();
        state = state.copyWith(reminders: updatedReminders);
        return true;
      }
      return false;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// OCR识别验光单
  /// 调用OCR API识别验光单图片，返回识别结果
  Future<VisionRecord?> recognizePrescription(String imagePath) async {
    try {
      // TODO: 待后端OCR接口实现后对接
      // 目前后端暂无OCR API端点，预留接口
      return null;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  // ==================== 私有辅助方法 ====================

  /// 将VisionRecordDTO转换为本地VisionRecord实体
  /// 注意：EyeDataDTO是嵌套结构，需要展平为本地扁平结构
  VisionRecord _convertVisionRecordDTO(VisionRecordDTO dto) {
    return VisionRecord(
      id: dto.id,
      userId: dto.userId,
      childId: dto.childId,
      date: dto.date.isNotEmpty ? DateTime.parse(dto.date) : DateTime.now(),
      // 将嵌套的EyeDataDTO转换为本地扁平的EyeData
      rightEye: _convertEyeDataDTO(dto.rightEye),
      leftEye: _convertEyeDataDTO(dto.leftEye),
      axialLengthRight: dto.axialLengthRight,
      axialLengthLeft: dto.axialLengthLeft,
      hyperopiaReserve: dto.hyperopiaReserve,
      source: dto.source == 'ocr'
          ? VisionDataSource.ocr
          : VisionDataSource.manual,
      imageUrl: dto.imageUrl,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将EyeDataDTO转换为本地EyeData实体
  /// 注意：DTO中axis为double（匹配后端float64），本地为int
  EyeData _convertEyeDataDTO(EyeDataDTO dto) {
    return EyeData(
      sph: dto.sph,
      cyl: dto.cyl,
      axis: dto.axis.round(), // double转int，展平嵌套结构
      va: dto.va,
    );
  }

  /// 将OutdoorActivityDTO转换为本地OutdoorActivity实体
  OutdoorActivity _convertOutdoorActivityDTO(OutdoorActivityDTO dto) {
    return OutdoorActivity(
      id: dto.id,
      userId: dto.userId,
      date: dto.date.isNotEmpty ? DateTime.parse(dto.date) : DateTime.now(),
      durationMin: dto.durationMin,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将EyeReminderDTO转换为本地EyeReminder实体
  EyeReminder _convertEyeReminderDTO(EyeReminderDTO dto) {
    return EyeReminder(
      id: dto.id,
      userId: dto.userId,
      type: _parseReminderType(dto.type),
      triggeredAt: dto.triggeredAt.isNotEmpty
          ? DateTime.parse(dto.triggeredAt)
          : DateTime.now(),
      acknowledged: dto.acknowledged,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 解析提醒类型字符串为枚举
  ReminderType _parseReminderType(String value) {
    return ReminderType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ReminderType.breakReminder,
    );
  }
}

/// 儿童列表Provider
final childrenListProvider = StateProvider<List>((ref) => []);
