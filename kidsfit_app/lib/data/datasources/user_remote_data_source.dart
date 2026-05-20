import '../../core/network/api_client.dart';
import '../../core/network/api_response.dart';
import '../../core/constants/app_constants.dart';

// ==================== DTO类定义 ====================

/// 登录响应DTO
/// 匹配后端登录响应：{"access_token":"xxx", "refresh_token":"xxx", "expires_in":7200, "user":{...}}
class LoginResponseDTO {
  /// 访问Token
  final String accessToken;

  /// 刷新Token
  final String refreshToken;

  /// 过期时间（秒）
  final int expiresIn;

  /// 用户信息
  final UserDTO user;

  const LoginResponseDTO({
    required this.accessToken,
    required this.refreshToken,
    required this.expiresIn,
    required this.user,
  });

  /// 从JSON创建LoginResponseDTO
  factory LoginResponseDTO.fromJson(Map<String, dynamic> json) {
    return LoginResponseDTO(
      accessToken: json['access_token'] as String? ?? '',
      refreshToken: json['refresh_token'] as String? ?? '',
      expiresIn: json['expires_in'] as int? ?? 0,
      user: UserDTO.fromJson(json['user'] as Map<String, dynamic>? ?? {}),
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'access_token': accessToken,
        'refresh_token': refreshToken,
        'expires_in': expiresIn,
        'user': user.toJson(),
      };
}

/// 用户DTO
/// 匹配后端用户数据结构
class UserDTO {
  /// 用户ID
  final String id;

  /// 用户类型（parent/child）
  final String type;

  /// 昵称
  final String nickname;

  /// 头像URL
  final String? avatar;

  /// 手机号
  final String? phone;

  /// 用户状态
  final String status;

  /// 创建时间
  final String createdAt;

  const UserDTO({
    required this.id,
    required this.type,
    required this.nickname,
    this.avatar,
    this.phone,
    this.status = 'active',
    required this.createdAt,
  });

  /// 从JSON创建UserDTO
  factory UserDTO.fromJson(Map<String, dynamic> json) {
    return UserDTO(
      id: json['id'] as String? ?? '',
      type: json['type'] as String? ?? 'parent',
      nickname: json['nickname'] as String? ?? '',
      avatar: json['avatar'] as String?,
      phone: json['phone'] as String?,
      status: json['status'] as String? ?? 'active',
      createdAt: json['created_at'] as String? ?? '',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'type': type,
        'nickname': nickname,
        'avatar': avatar,
        'phone': phone,
        'status': status,
        'created_at': createdAt,
      };
}

/// 儿童DTO
/// 匹配后端儿童数据结构
class ChildDTO {
  /// 儿童ID
  final String id;

  /// 昵称
  final String nickname;

  /// 头像URL
  final String? avatar;

  /// 年龄
  final int age;

  /// 年龄组
  final String? ageGroup;

  /// 状态
  final String status;

  const ChildDTO({
    required this.id,
    required this.nickname,
    this.avatar,
    required this.age,
    this.ageGroup,
    this.status = 'active',
  });

  /// 从JSON创建ChildDTO
  factory ChildDTO.fromJson(Map<String, dynamic> json) {
    return ChildDTO(
      id: json['id'] as String? ?? '',
      nickname: json['nickname'] as String? ?? '',
      avatar: json['avatar'] as String?,
      age: json['age'] as int? ?? 0,
      ageGroup: json['age_group'] as String?,
      status: json['status'] as String? ?? 'active',
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'id': id,
        'nickname': nickname,
        'avatar': avatar,
        'age': age,
        'age_group': ageGroup,
        'status': status,
      };
}

/// 家长设置DTO
/// 匹配后端家长设置数据结构
class ParentSettingsDTO {
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

  const ParentSettingsDTO({
    this.dailyLimitMin = 30,
    this.availableFrom = '08:00',
    this.availableTo = '21:00',
    this.cameraAllowed = true,
    this.locationAllowed = true,
    this.dataUploadCloud = false,
  });

  /// 从JSON创建ParentSettingsDTO
  factory ParentSettingsDTO.fromJson(Map<String, dynamic> json) {
    return ParentSettingsDTO(
      dailyLimitMin: json['daily_limit_min'] as int? ?? 30,
      availableFrom: json['available_from'] as String? ?? '08:00',
      availableTo: json['available_to'] as String? ?? '21:00',
      cameraAllowed: json['camera_allowed'] as bool? ?? true,
      locationAllowed: json['location_allowed'] as bool? ?? true,
      dataUploadCloud: json['data_upload_cloud'] as bool? ?? false,
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'daily_limit_min': dailyLimitMin,
        'available_from': availableFrom,
        'available_to': availableTo,
        'camera_allowed': cameraAllowed,
        'location_allowed': locationAllowed,
        'data_upload_cloud': dataUploadCloud,
      };
}

// ==================== 用户远程数据源 ====================

/// 用户远程数据源
/// 负责与后端用户相关API通信
class UserRemoteDataSource {
  /// API客户端实例
  final ApiClient _apiClient;

  UserRemoteDataSource({ApiClient? apiClient})
      : _apiClient = apiClient ?? ApiClient.instance;

  /// 用户注册
  /// POST /api/v1/auth/register
  Future<ApiResponse<LoginResponseDTO>> register({
    required String phone,
    required String password,
    required String nickname,
  }) async {
    return _apiClient.post(
      AppConstants.authRegister,
      data: {
        'phone': phone,
        'password': password,
        'nickname': nickname,
      },
      fromJsonT: (data) =>
          LoginResponseDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 用户登录
  /// POST /api/v1/auth/login
  Future<ApiResponse<LoginResponseDTO>> login({
    required String phone,
    required String password,
  }) async {
    return _apiClient.post(
      AppConstants.authLogin,
      data: {
        'phone': phone,
        'password': password,
      },
      fromJsonT: (data) =>
          LoginResponseDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 刷新Token
  /// POST /api/v1/auth/refresh
  Future<ApiResponse<LoginResponseDTO>> refreshToken(
      String refreshToken) async {
    return _apiClient.post(
      AppConstants.authRefresh,
      data: {'refresh_token': refreshToken},
      fromJsonT: (data) =>
          LoginResponseDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 退出登录
  /// POST /api/v1/auth/logout
  Future<ApiResponse<void>> logout() async {
    return _apiClient.postVoid(AppConstants.authLogout);
  }

  /// 获取当前用户信息
  /// GET /api/v1/users/me
  Future<ApiResponse<UserDTO>> getCurrentUser() async {
    return _apiClient.get(
      AppConstants.userMe,
      fromJsonT: (data) => UserDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 更新用户信息
  /// PUT /api/v1/users/me
  Future<ApiResponse<UserDTO>> updateUser({
    String? nickname,
    String? avatar,
  }) async {
    final data = <String, dynamic>{};
    if (nickname != null) {
      data['nickname'] = nickname;
    }
    if (avatar != null) {
      data['avatar'] = avatar;
    }

    return _apiClient.put(
      AppConstants.userUpdate,
      data: data,
      fromJsonT: (data) => UserDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 添加儿童
  /// POST /api/v1/users/children
  Future<ApiResponse<ChildDTO>> addChild({
    required String nickname,
    required int age,
    String? avatar,
  }) async {
    return _apiClient.post(
      AppConstants.userAddChild,
      data: {
        'nickname': nickname,
        'age': age,
        if (avatar != null) 'avatar': avatar,
      },
      fromJsonT: (data) => ChildDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 获取儿童列表
  /// GET /api/v1/users/children
  Future<ApiResponse<List<ChildDTO>>> getChildren() async {
    return _apiClient.get(
      AppConstants.userChildren,
      fromJsonT: (data) {
        final list = data as List<dynamic>;
        return list
            .map((e) => ChildDTO.fromJson(e as Map<String, dynamic>))
            .toList();
      },
    );
  }

  /// 获取家长设置
  /// GET /api/v1/users/parent-settings
  Future<ApiResponse<ParentSettingsDTO>> getParentSettings() async {
    return _apiClient.get(
      AppConstants.parentSettings,
      fromJsonT: (data) =>
          ParentSettingsDTO.fromJson(data as Map<String, dynamic>),
    );
  }

  /// 更新家长设置
  /// PUT /api/v1/users/parent-settings
  Future<ApiResponse<ParentSettingsDTO>> updateParentSettings(
      ParentSettingsDTO dto) async {
    return _apiClient.put(
      AppConstants.parentSettings,
      data: dto.toJson(),
      fromJsonT: (data) =>
          ParentSettingsDTO.fromJson(data as Map<String, dynamic>),
    );
  }
}
