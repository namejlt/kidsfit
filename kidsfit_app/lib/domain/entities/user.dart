import 'package:equatable/equatable.dart';

/// 用户类型枚举
enum UserType {
  /// 家长
  parent('parent'),

  /// 儿童
  child('child');

  const UserType(this.value);
  final String value;
}

/// 用户状态枚举
enum UserStatus {
  /// 激活
  active('active'),

  /// 未激活
  inactive('inactive'),

  /// 已删除
  deleted('deleted');

  const UserStatus(this.value);
  final String value;
}

/// 用户实体
/// 代表KidsFit应用中的用户（家长或儿童）
class User extends Equatable {
  /// 用户ID
  final String id;

  /// 用户类型
  final UserType type;

  /// 家长ID（仅儿童用户有）
  final String? parentId;

  /// 年龄（仅儿童用户有）
  /// 有效范围：3-12岁
  final int? age;

  /// 昵称
  final String nickname;

  /// 头像URL
  final String? avatar;

  /// 手机号（仅家长用户有）
  final String? phone;

  /// 用户状态
  final UserStatus status;

  /// 创建时间
  final DateTime createdAt;

  /// 更新时间
  final DateTime? updatedAt;

  const User({
    required this.id,
    required this.type,
    this.parentId,
    this.age,
    required this.nickname,
    this.avatar,
    this.phone,
    this.status = UserStatus.active,
    required this.createdAt,
    this.updatedAt,
  });

  /// 判断是否为家长用户
  bool get isParent => type == UserType.parent;

  /// 判断是否为儿童用户
  bool get isChild => type == UserType.child;

  /// 判断用户是否激活
  bool get isActive => status == UserStatus.active;

  /// 年龄组分类
  /// 根据年龄返回对应的年龄组
  String? get ageGroup {
    if (age == null) return null;
    if (age! >= 3 && age! <= 6) return '3-6岁';
    if (age! >= 7 && age! <= 9) return '7-9岁';
    if (age! >= 10 && age! <= 12) return '10-12岁';
    return null;
  }

  /// 创建新用户
  User copyWith({
    String? id,
    UserType? type,
    String? parentId,
    int? age,
    String? nickname,
    String? avatar,
    String? phone,
    UserStatus? status,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return User(
      id: id ?? this.id,
      type: type ?? this.type,
      parentId: parentId ?? this.parentId,
      age: age ?? this.age,
      nickname: nickname ?? this.nickname,
      avatar: avatar ?? this.avatar,
      phone: phone ?? this.phone,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        type,
        parentId,
        age,
        nickname,
        avatar,
        phone,
        status,
        createdAt,
        updatedAt,
      ];

  @override
  String toString() {
    return 'User(id: $id, type: $type, nickname: $nickname)';
  }
}

/// 家长设置实体
class ParentSettings extends Equatable {
  /// 设置ID
  final String id;

  /// 家长ID
  final String parentId;

  /// 每日使用时长限制（分钟）
  final int dailyLimitMin;

  /// 可用时段开始时间
  final String availableFrom;

  /// 可用时段结束时间
  final String availableTo;

  /// 是否允许使用摄像头
  final bool cameraAllowed;

  /// 是否允许获取位置
  final bool locationAllowed;

  /// 是否允许数据上传云端
  final bool dataUploadCloud;

  /// 创建时间
  final DateTime createdAt;

  /// 更新时间
  final DateTime? updatedAt;

  const ParentSettings({
    required this.id,
    required this.parentId,
    this.dailyLimitMin = 30,
    this.availableFrom = '08:00',
    this.availableTo = '21:00',
    this.cameraAllowed = true,
    this.locationAllowed = true,
    this.dataUploadCloud = false,
    required this.createdAt,
    this.updatedAt,
  });

  /// 验证当前时间是否在可用时段内
  bool isWithinAvailableTime() {
    final now = DateTime.now();
    final currentTime =
        '${now.hour.toString().padLeft(2, '0')}:${now.minute.toString().padLeft(2, '0')}';
    return currentTime.compareTo(availableFrom) >= 0 &&
        currentTime.compareTo(availableTo) <= 0;
  }

  ParentSettings copyWith({
    String? id,
    String? parentId,
    int? dailyLimitMin,
    String? availableFrom,
    String? availableTo,
    bool? cameraAllowed,
    bool? locationAllowed,
    bool? dataUploadCloud,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return ParentSettings(
      id: id ?? this.id,
      parentId: parentId ?? this.parentId,
      dailyLimitMin: dailyLimitMin ?? this.dailyLimitMin,
      availableFrom: availableFrom ?? this.availableFrom,
      availableTo: availableTo ?? this.availableTo,
      cameraAllowed: cameraAllowed ?? this.cameraAllowed,
      locationAllowed: locationAllowed ?? this.locationAllowed,
      dataUploadCloud: dataUploadCloud ?? this.dataUploadCloud,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        parentId,
        dailyLimitMin,
        availableFrom,
        availableTo,
        cameraAllowed,
        locationAllowed,
        dataUploadCloud,
        createdAt,
        updatedAt,
      ];
}

/// 家庭关系实体
class Family extends Equatable {
  /// 家庭ID
  final String id;

  /// 家长ID
  final String parentId;

  /// 儿童ID
  final String childId;

  /// 关系类型
  final String relation;

  /// 创建时间
  final DateTime createdAt;

  const Family({
    required this.id,
    required this.parentId,
    required this.childId,
    required this.relation,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [id, parentId, childId, relation, createdAt];
}
