import 'package:equatable/equatable.dart';

/// 数据来源枚举
enum VisionDataSource {
  /// OCR识别
  ocr('ocr'),

  /// 手动录入
  manual('manual');

  const VisionDataSource(this.value);
  final String value;
}

/// 提醒类型枚举
enum ReminderType {
  /// 20-20-20法则提醒
  rule202020('20_20_20', '20-20-20法则'),

  /// 户外运动提醒
  outdoor('outdoor', '户外运动提醒'),

  /// 休息提醒
  breakReminder('break', '休息提醒');

  const ReminderType(this.value, this.displayName);
  final String value;
  final String displayName;
}

/// 单眼视力数据
class EyeData extends Equatable {
  /// 球镜度数
  final double sph;

  /// 柱镜度数
  final double cyl;

  /// 轴位
  final int axis;

  /// 矫正视力
  final double va;

  const EyeData({
    this.sph = 0.0,
    this.cyl = 0.0,
    this.axis = 0,
    this.va = 1.0,
  });

  /// 是否近视
  bool get isMyopic => sph < -0.25;

  /// 是否远视
  bool get isHyperopic => sph > 0.25;

  /// 是否散光
  bool get hasAstigmatism => cyl.abs() > 0.5;

  /// 获取等效球镜
  double get sphericalEquivalent => sph + cyl / 2;

  @override
  List<Object?> get props => [sph, cyl, axis, va];
}

/// 视力档案记录实体
class VisionRecord extends Equatable {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 儿童ID
  final String childId;

  /// 记录日期
  final DateTime date;

  /// 右眼数据
  final EyeData rightEye;

  /// 左眼数据
  final EyeData leftEye;

  /// 右眼眼轴长度
  final double? axialLengthRight;

  /// 左眼眼轴长度
  final double? axialLengthLeft;

  /// 远视储备
  final double? hyperopiaReserve;

  /// 数据来源
  final VisionDataSource source;

  /// 验光单照片URL
  final String? imageUrl;

  /// 创建时间
  final DateTime createdAt;

  const VisionRecord({
    required this.id,
    required this.userId,
    required this.childId,
    required this.date,
    this.rightEye = const EyeData(),
    this.leftEye = const EyeData(),
    this.axialLengthRight,
    this.axialLengthLeft,
    this.hyperopiaReserve,
    this.source = VisionDataSource.manual,
    this.imageUrl,
    required this.createdAt,
  });

  /// 获取双眼平均球镜
  double get averageSph => (rightEye.sph + leftEye.sph) / 2;

  /// 获取双眼平均视力
  double get averageVa => (rightEye.va + leftEye.va) / 2;

  /// 获取屈光状态描述
  String get refractiveStatus {
    if (averageSph < -0.5) return '近视';
    if (averageSph > 0.5) return '远视';
    return '正视';
  }

  /// 获取视力状态
  VisionStatus get visionStatus {
    if (averageVa >= 1.0) return VisionStatus.good;
    if (averageVa >= 0.8) return VisionStatus.medium;
    return VisionStatus.concern;
  }

  VisionRecord copyWith({
    String? id,
    String? userId,
    String? childId,
    DateTime? date,
    EyeData? rightEye,
    EyeData? leftEye,
    double? axialLengthRight,
    double? axialLengthLeft,
    double? hyperopiaReserve,
    VisionDataSource? source,
    String? imageUrl,
    DateTime? createdAt,
  }) {
    return VisionRecord(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      childId: childId ?? this.childId,
      date: date ?? this.date,
      rightEye: rightEye ?? this.rightEye,
      leftEye: leftEye ?? this.leftEye,
      axialLengthRight: axialLengthRight ?? this.axialLengthRight,
      axialLengthLeft: axialLengthLeft ?? this.axialLengthLeft,
      hyperopiaReserve: hyperopiaReserve ?? this.hyperopiaReserve,
      source: source ?? this.source,
      imageUrl: imageUrl ?? this.imageUrl,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        childId,
        date,
        rightEye,
        leftEye,
        axialLengthRight,
        axialLengthLeft,
        hyperopiaReserve,
        source,
        imageUrl,
        createdAt,
      ];
}

/// 视力状态枚举
enum VisionStatus {
  /// 良好
  good('良好', '😀'),

  /// 中等
  medium('关注', '😐'),

  /// 需要关注
  concern('需关注', '😟');

  const VisionStatus(this.displayName, this.emoji);
  final String displayName;
  final String emoji;
}

/// 户外运动记录实体
class OutdoorActivity extends Equatable {
  /// 记录ID
  final String id;

  /// 用户ID
  final String userId;

  /// 活动日期
  final DateTime date;

  /// 总时长（分钟）
  final int durationMin;

  /// 户外时段列表
  final List<OutdoorSegment> segments;

  /// 创建时间
  final DateTime createdAt;

  const OutdoorActivity({
    required this.id,
    required this.userId,
    required this.date,
    this.durationMin = 0,
    this.segments = const [],
    required this.createdAt,
  });

  /// 目标户外时长（分钟）
  static const int targetMinutes = 120;

  /// 最低有效户外时长（分钟）
  static const int minimumMinutes = 60;

  /// 是否达标
  bool get isTargetMet => durationMin >= targetMinutes;

  /// 是否达到最低标准
  bool get isMinimumMet => durationMin >= minimumMinutes;

  /// 目标完成百分比
  double get targetProgress =>
      (durationMin / targetMinutes * 100).clamp(0, 100).toDouble();

  /// 距离目标还差多少分钟
  int get remainingMinutes =>
      (targetMinutes - durationMin).clamp(0, targetMinutes);

  /// 格式化时长显示
  String get formattedDuration {
    if (durationMin < 60) {
      return '$durationMin分钟';
    }
    final hours = durationMin ~/ 60;
    final minutes = durationMin % 60;
    if (minutes == 0) {
      return '${hours}小时';
    }
    return '${hours}小时${minutes}分钟';
  }

  OutdoorActivity copyWith({
    String? id,
    String? userId,
    DateTime? date,
    int? durationMin,
    List<OutdoorSegment>? segments,
    DateTime? createdAt,
  }) {
    return OutdoorActivity(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      date: date ?? this.date,
      durationMin: durationMin ?? this.durationMin,
      segments: segments ?? this.segments,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        date,
        durationMin,
        segments,
        createdAt,
      ];
}

/// 户外时段实体
class OutdoorSegment extends Equatable {
  /// 时段ID
  final String id;

  /// 所属活动ID
  final String activityId;

  /// 开始时间
  final DateTime startTime;

  /// 结束时间
  final DateTime endTime;

  /// 时长（分钟）
  final int durationMin;

  /// 地点描述
  final String? location;

  const OutdoorSegment({
    required this.id,
    required this.activityId,
    required this.startTime,
    required this.endTime,
    required this.durationMin,
    this.location,
  });

  @override
  List<Object?> get props => [
        id,
        activityId,
        startTime,
        endTime,
        durationMin,
        location,
      ];
}

/// 用眼提醒实体
class EyeReminder extends Equatable {
  /// 提醒ID
  final String id;

  /// 用户ID
  final String userId;

  /// 提醒类型
  final ReminderType type;

  /// 触发时间
  final DateTime triggeredAt;

  /// 是否已确认
  final bool acknowledged;

  /// 确认时间
  final DateTime? acknowledgedAt;

  /// 创建时间
  final DateTime createdAt;

  const EyeReminder({
    required this.id,
    required this.userId,
    required this.type,
    required this.triggeredAt,
    this.acknowledged = false,
    this.acknowledgedAt,
    required this.createdAt,
  });

  /// 获取提醒标题
  String get title {
    switch (type) {
      case ReminderType.rule202020:
        return '远眺休息';
      case ReminderType.outdoor:
        return '户外运动提醒';
      case ReminderType.breakReminder:
        return '休息提醒';
    }
  }

  /// 获取提醒内容
  String get content {
    switch (type) {
      case ReminderType.rule202020:
        return '已经看了20分钟啦，快看看20英尺（6米）外的东西，休息20秒吧！';
      case ReminderType.outdoor:
        return '已经在室内很久啦，去窗边看看远方吧！';
      case ReminderType.breakReminder:
        return '该休息一下了，做做眼保健操吧！';
    }
  }

  EyeReminder copyWith({
    String? id,
    String? userId,
    ReminderType? type,
    DateTime? triggeredAt,
    bool? acknowledged,
    DateTime? acknowledgedAt,
    DateTime? createdAt,
  }) {
    return EyeReminder(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      type: type ?? this.type,
      triggeredAt: triggeredAt ?? this.triggeredAt,
      acknowledged: acknowledged ?? this.acknowledged,
      acknowledgedAt: acknowledgedAt ?? this.acknowledgedAt,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        type,
        triggeredAt,
        acknowledged,
        acknowledgedAt,
        createdAt,
      ];
}
