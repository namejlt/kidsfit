import '../../core/network/api_client.dart';
import '../../core/network/api_response.dart';
import '../../core/constants/app_constants.dart';

// ==================== DTO类定义 ====================

/// 单眼视力数据DTO
/// 匹配后端EyeData嵌套结构，axis使用double匹配后端float64
class EyeDataDTO {
  /// 球镜度数
  final double sph;

  /// 柱镜度数
  final double cyl;

  /// 轴位（double匹配后端float64）
  final double axis;

  /// 矫正视力
  final double va;

  const EyeDataDTO({
    this.sph = 0.0,
    this.cyl = 0.0,
    this.axis = 0.0,
    this.va = 1.0,
  });

  /// 从JSON创建EyeDataDTO
  factory EyeDataDTO.fromJson(Map<String, dynamic> json) {
    return EyeDataDTO(
      sph: (json['sph'] as num?)?.toDouble() ?? 0.0,
      cyl: (json['cyl'] as num?)?.toDouble() ?? 0.0,
      axis: (json['axis'] as num?)?.toDouble() ?? 0.0,
      va: (json['va'] as num?)?.toDouble() ?? 1.0,
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'sph': sph,
        'cyl': cyl,
        'axis': axis,
        'va': va,
      };
}

/// 创建视力记录请求DTO
class CreateVisionRecordRequest {
  /// 儿童ID
  final String childId;

  /// 记录日期
  final String date;

  /// 右眼数据
  final EyeDataDTO rightEye;

  /// 左眼数据
  final EyeDataDTO leftEye;

  /// 右眼眼轴长度
  final double? axialLengthRight;

  /// 左眼眼轴长度
  final double? axialLengthLeft;

  /// 远视储备
  final double? hyperopiaReserve;

  /// 数据来源
  final String source;

  /// 验光单照片URL
  final String? imageUrl;

  const CreateVisionRecordRequest({
    required this.childId,
    required this.date,
    required this.rightEye,
    required this.leftEye,
    this.axialLengthRight,
    this.axialLengthLeft,
    this.hyperopiaReserve,
    this.source = 'manual',
    this.imageUrl,
  });

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'child_id': childId,
        'date': date,
        'right_eye': rightEye.toJson(),
        'left_eye': leftEye.toJson(),
        'axial_length_right': axialLengthRight,
        'axial_length_left': axialLengthLeft,
        'hyperopia_reserve': hyperopiaReserve,
        'source': source,
        'image_url': imageUrl,
      };
}

/// 视力记录DTO
class VisionRecordDTO {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 儿童ID
  final String childId;

  /// 记录日期
  final String date;

  /// 右眼数据
  final EyeDataDTO rightEye;

  /// 左眼数据
  final EyeDataDTO leftEye;

  /// 右眼眼轴长度
  final double? axialLengthRight;

  /// 左眼眼轴长度
  final double? axialLengthLeft;

  /// 远视储备
  final double? hyperopiaReserve;

  /// 数据来源
  final String source;

  /// 验光单照片URL
  final String? imageUrl;

  /// 视力状态
  final String? visionStatus;

  /// 屈光状态
  final String? refractiveStatus;

  /// 创建时间
  final String createdAt;

  const VisionRecordDTO({
    required this.id,
    required this.userId,
    required this.childId,
    required this.date,
    this.rightEye = const EyeDataDTO(),
    this.leftEye = const EyeDataDTO(),
    this.axialLengthRight,
    this.axialLengthLeft,
    this.hyperopiaReserve,
    this.source = 'manual',
    this.imageUrl,
    this.visionStatus,
    this.refractiveStatus,
    required this.createdAt,
  });

  /// 从JSON创建VisionRecordDTO
  factory VisionRecordDTO.fromJson(Map<String, dynamic> json) {
    return VisionRecordDTO(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      childId: json['child_id'] as String? ?? '',
      date: json['date'] as String? ?? '',
      rightEye: json['right_eye'] != null
          ? EyeDataDTO.fromJson(json['right_eye'] as Map<String, dynamic>)
          : const EyeDataDTO(),
      leftEye: json['left_eye'] != null
          ? EyeDataDTO.fromJson(json['left_eye'] as Map<String, dynamic>)
          : const EyeDataDTO(),
      axialLengthRight:
          (json['axial_length_right'] as num?)?.toDouble(),
      axialLengthLeft:
          (json['axial_length_left'] as num?)?.toDouble(),
      hyperopiaReserve:
          (json['hyperopia_reserve'] as num?)?.toDouble(),
      source: json['source'] as String? ?? 'manual',
      imageUrl: json['image_url'] as String?,
      visionStatus: json['vision_status'] as String?,
      refractiveStatus: json['refractive_status'] as String?,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'user_id': userId,
        'child_id': childId,
        'date': date,
        'right_eye': rightEye.toJson(),
        'left_eye': leftEye.toJson(),
        'axial_length_right': axialLengthRight,
        'axial_length_left': axialLengthLeft,
        'hyperopia_reserve': hyperopiaReserve,
        'source': source,
        'image_url': imageUrl,
        'vision_status': visionStatus,
        'refractive_status': refractiveStatus,
        'created_at': createdAt,
      };
}

/// 户外活动DTO
class OutdoorActivityDTO {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 活动日期
  final String date;

  /// 时长（分钟）
  final int durationMin;

  /// 是否达标
  final bool isTargetMet;

  /// 目标完成进度
  final double targetProgress;

  /// 创建时间
  final String createdAt;

  const OutdoorActivityDTO({
    required this.id,
    required this.userId,
    required this.date,
    this.durationMin = 0,
    this.isTargetMet = false,
    this.targetProgress = 0,
    required this.createdAt,
  });

  /// 从JSON创建OutdoorActivityDTO
  factory OutdoorActivityDTO.fromJson(Map<String, dynamic> json) {
    return OutdoorActivityDTO(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      date: json['date'] as String? ?? '',
      durationMin: json['duration_min'] as int? ?? 0,
      isTargetMet: json['is_target_met'] as bool? ?? false,
      targetProgress: (json['target_progress'] as num?)?.toDouble() ?? 0,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'user_id': userId,
        'date': date,
        'duration_min': durationMin,
        'is_target_met': isTargetMet,
        'target_progress': targetProgress,
        'created_at': createdAt,
      };
}

/// 用眼提醒DTO
class EyeReminderDTO {
  /// 提醒ID
  final String id;

  /// 用户ID
  final String userId;

  /// 提醒类型
  final String type;

  /// 触发时间
  final String triggeredAt;

  /// 是否已确认
  final bool acknowledged;

  /// 创建时间
  final String createdAt;

  const EyeReminderDTO({
    required this.id,
    required this.userId,
    required this.type,
    required this.triggeredAt,
    this.acknowledged = false,
    required this.createdAt,
  });

  /// 从JSON创建EyeReminderDTO
  factory EyeReminderDTO.fromJson(Map<String, dynamic> json) {
    return EyeReminderDTO(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      type: json['type'] as String? ?? '',
      triggeredAt: json['triggered_at'] as String? ?? '',
      acknowledged: json['acknowledged'] as bool? ?? false,
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'user_id': userId,
        'type': type,
        'triggered_at': triggeredAt,
        'acknowledged': acknowledged,
        'created_at': createdAt,
      };
}

/// 视力趋势DTO
class VisionTrendDTO {
  /// 视力记录列表
  final List<VisionRecordDTO> records;

  const VisionTrendDTO({
    this.records = const [],
  });

  /// 从JSON创建VisionTrendDTO
  factory VisionTrendDTO.fromJson(Map<String, dynamic> json) {
    final recordsRaw = json['records'] as List<dynamic>? ?? [];
    return VisionTrendDTO(
      records: recordsRaw
          .map((e) => VisionRecordDTO.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'records': records.map((e) => e.toJson()).toList(),
      };
}

// ==================== 视力远程数据源 ====================

/// 视力远程数据源
/// 负责与后端视力相关API通信
class VisionRemoteDataSource {
  /// API客户端实例
  final ApiClient _apiClient;

  VisionRemoteDataSource({ApiClient? apiClient})
      : _apiClient = apiClient ?? ApiClient.instance;

  /// 创建视力记录
  /// POST /api/v1/vision/records
  Future<ApiResponse<VisionRecordDTO>> createVisionRecord(
      CreateVisionRecordRequest req) async {
    return _apiClient.post(
      AppConstants.visionRecords,
      data: req.toJson(),
      fromJsonT: (data) =>
          VisionRecordDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取视力记录列表（分页）
  /// GET /api/v1/vision/records
  Future<ApiResponse<PaginatedData<VisionRecordDTO>>> getVisionRecords({
    String? childId,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (childId != null) {
      queryParams['child_id'] = childId;
    }

    return _apiClient.getPaged(
      AppConstants.visionRecords,
      queryParameters: queryParams,
      fromJsonT: (data) =>
          VisionRecordDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取视力趋势
  /// GET /api/v1/vision/trend
  Future<ApiResponse<VisionTrendDTO>> getVisionTrend({
    String? childId,
  }) async {
    final queryParams = <String, dynamic>{};
    if (childId != null) {
      queryParams['child_id'] = childId;
    }

    return _apiClient.get(
      AppConstants.visionTrend,
      queryParameters: queryParams,
      fromJsonT: (data) =>
          VisionTrendDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取今日户外活动
  /// GET /api/v1/vision/outdoor/today
  Future<ApiResponse<OutdoorActivityDTO>> getTodayOutdoor() async {
    return _apiClient.get(
      AppConstants.todayOutdoor,
      fromJsonT: (data) =>
          OutdoorActivityDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 同步户外活动数据
  /// POST /api/v1/vision/outdoor/sync
  Future<ApiResponse<OutdoorActivityDTO>> syncOutdoorData(
      int durationMin) async {
    return _apiClient.post(
      AppConstants.syncOutdoor,
      data: {'duration_min': durationMin},
      fromJsonT: (data) =>
          OutdoorActivityDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取用眼提醒列表（分页）
  /// GET /api/v1/vision/reminders
  Future<ApiResponse<PaginatedData<EyeReminderDTO>>> getReminders({
    int page = 1,
    int pageSize = 20,
  }) async {
    return _apiClient.getPaged(
      AppConstants.eyeReminders,
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
      fromJsonT: (data) =>
          EyeReminderDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 确认提醒
  /// PUT /api/v1/vision/reminders/{reminderId}/ack
  Future<ApiResponse<void>> ackReminder(String reminderId) async {
    return _apiClient.putVoid(
      '${AppConstants.eyeReminders}/$reminderId/ack',
    );
  }
}
