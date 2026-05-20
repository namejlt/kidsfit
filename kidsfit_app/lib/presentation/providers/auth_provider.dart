import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/user.dart';
import '../../services/storage_service.dart';
import '../../data/datasources/user_remote_data_source.dart';
import '../../core/constants/app_constants.dart';

/// 认证状态枚举
enum AuthStatus {
  /// 未知
  unknown,

  /// 已认证
  authenticated,

  /// 未认证
  unauthenticated,
}

/// 认证状态数据类
class AuthState {
  /// 当前认证状态
  final AuthStatus status;

  /// 当前用户
  final User? user;

  /// Token
  final String? token;

  /// 刷新Token
  final String? refreshToken;

  /// 错误消息
  final String? error;

  /// 是否加载中
  final bool isLoading;

  const AuthState({
    this.status = AuthStatus.unknown,
    this.user,
    this.token,
    this.refreshToken,
    this.error,
    this.isLoading = false,
  });

  /// 是否已认证
  bool get isAuthenticated => status == AuthStatus.authenticated;

  /// 是否为家长
  bool get isParent => user?.isParent ?? false;

  /// 是否为儿童
  bool get isChild => user?.isChild ?? false;

  /// 创建认证成功状态
  factory AuthState.authenticated({
    required User user,
    required String token,
    String? refreshToken,
  }) {
    return AuthState(
      status: AuthStatus.authenticated,
      user: user,
      token: token,
      refreshToken: refreshToken,
    );
  }

  /// 创建未认证状态
  factory AuthState.unauthenticated({String? error}) {
    return AuthState(
      status: AuthStatus.unauthenticated,
      error: error,
    );
  }

  /// 创建加载状态
  factory AuthState.loading() {
    return const AuthState(status: AuthStatus.unknown, isLoading: true);
  }

  AuthState copyWith({
    AuthStatus? status,
    User? user,
    String? token,
    String? refreshToken,
    String? error,
    bool? isLoading,
  }) {
    return AuthState(
      status: status ?? this.status,
      user: user ?? this.user,
      token: token ?? this.token,
      refreshToken: refreshToken ?? this.refreshToken,
      error: error,
      isLoading: isLoading ?? this.isLoading,
    );
  }
}

/// 用户远程数据源Provider
final userRemoteDataSourceProvider = Provider<UserRemoteDataSource>((ref) {
  return UserRemoteDataSource();
});

/// 认证状态Provider
final authStateProvider =
    StateNotifierProvider<AuthStateNotifier, AuthState>((ref) {
  final dataSource = ref.watch(userRemoteDataSourceProvider);
  return AuthStateNotifier(dataSource);
});

/// 认证状态Notifier
class AuthStateNotifier extends StateNotifier<AuthState> {
  /// 用户远程数据源
  final UserRemoteDataSource _dataSource;

  AuthStateNotifier(this._dataSource) : super(AuthState.loading()) {
    _checkAuthStatus();
  }

  /// 检查本地认证状态
  /// 检查本地token是否存在且未过期，如果有效则恢复登录态
  Future<void> _checkAuthStatus() async {
    try {
      final token = await StorageService.getToken();
      if (token == null || token.isEmpty) {
        state = AuthState.unauthenticated();
        return;
      }

      // 检查token是否过期
      final expiryMs = await StorageService.getSetting<int>(
          AppConstants.tokenExpiryKey);
      if (expiryMs != null) {
        final expiryTime = DateTime.fromMillisecondsSinceEpoch(expiryMs);
        if (DateTime.now().isAfter(expiryTime)) {
          // Token已过期，尝试刷新
          final refreshed = await refreshToken();
          if (!refreshed) {
            state = AuthState.unauthenticated();
            return;
          }
          return; // refreshToken内部会更新state
        }
      }

      // Token有效，获取当前用户信息
      final response = await _dataSource.getCurrentUser();
      if (response.isSuccess && response.data != null) {
        final user = _convertUserDTO(response.data!);
        final refreshTokenValue =
            await StorageService.getRefreshToken();
        state = AuthState.authenticated(
          user: user,
          token: token,
          refreshToken: refreshTokenValue,
        );
      } else {
        // Token无效或服务端错误，尝试用本地缓存
        final userData = await StorageService.getUserData();
        if (userData != null) {
          final user = User(
            id: userData['id'] as String,
            type: userData['type'] == 'parent'
                ? UserType.parent
                : UserType.child,
            nickname: userData['nickname'] as String,
            age: userData['age'] as int?,
            avatar: userData['avatar'] as String?,
            phone: userData['phone'] as String?,
            createdAt: DateTime.parse(userData['created_at'] as String),
          );
          state = AuthState.authenticated(user: user, token: token);
        } else {
          state = AuthState.unauthenticated();
        }
      }
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
    }
  }

  /// 用户注册
  /// 调用API注册新用户，成功后保存token和用户信息
  Future<bool> register({
    required String phone,
    required String password,
    required String nickname,
  }) async {
    try {
      state = AuthState.loading();

      final response = await _dataSource.register(
        phone: phone,
        password: password,
        nickname: nickname,
      );

      if (response.isSuccess && response.data != null) {
        final loginData = response.data!;
        final user = _convertUserDTO(loginData.user);

        // 保存token和用户数据到本地
        await _saveAuthData(
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
          expiresIn: loginData.expiresIn,
          user: user,
        );

        state = AuthState.authenticated(
          user: user,
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
        );
        return true;
      } else {
        state = AuthState.unauthenticated(error: response.message);
        return false;
      }
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
      return false;
    }
  }

  /// 用户登录
  /// 调用API登录，成功后保存token和用户信息
  Future<bool> login({
    required String phone,
    required String password,
  }) async {
    try {
      state = AuthState.loading();

      final response = await _dataSource.login(
        phone: phone,
        password: password,
      );

      if (response.isSuccess && response.data != null) {
        final loginData = response.data!;
        final user = _convertUserDTO(loginData.user);

        // 保存token和用户数据到本地
        await _saveAuthData(
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
          expiresIn: loginData.expiresIn,
          user: user,
        );

        state = AuthState.authenticated(
          user: user,
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
        );
        return true;
      } else {
        state = AuthState.unauthenticated(error: response.message);
        return false;
      }
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
      return false;
    }
  }

  /// 用户登出
  /// 调用API登出，清除本地存储
  Future<void> logout() async {
    try {
      await _dataSource.logout();
    } catch (_) {
      // 即使API调用失败也清除本地数据
    }
    await StorageService.clearAll();
    state = AuthState.unauthenticated();
  }

  /// 刷新Token
  /// 调用API刷新token，成功后更新本地存储
  Future<bool> refreshToken() async {
    try {
      final storedRefreshToken = await StorageService.getRefreshToken();
      if (storedRefreshToken == null || storedRefreshToken.isEmpty) {
        return false;
      }

      final response =
          await _dataSource.refreshToken(storedRefreshToken);
      if (response.isSuccess && response.data != null) {
        final loginData = response.data!;
        final user = _convertUserDTO(loginData.user);

        // 保存新token
        await _saveAuthData(
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
          expiresIn: loginData.expiresIn,
          user: user,
        );

        state = AuthState.authenticated(
          user: user,
          token: loginData.accessToken,
          refreshToken: loginData.refreshToken,
        );
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  /// 获取当前用户信息
  /// 调用API获取当前用户，更新本地状态
  Future<User?> getCurrentUser() async {
    try {
      final response = await _dataSource.getCurrentUser();
      if (response.isSuccess && response.data != null) {
        final user = _convertUserDTO(response.data!);
        // 更新本地缓存
        await _saveUserData(user);
        state = state.copyWith(user: user);
        return user;
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// 更新用户信息
  /// 调用API更新用户昵称或头像
  Future<bool> updateUser({
    String? nickname,
    String? avatar,
  }) async {
    try {
      final response = await _dataSource.updateUser(
        nickname: nickname,
        avatar: avatar,
      );
      if (response.isSuccess && response.data != null) {
        final user = _convertUserDTO(response.data!);
        await _saveUserData(user);
        state = state.copyWith(user: user);
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  /// 添加儿童
  /// 调用API添加儿童用户
  Future<bool> addChild({
    required String nickname,
    required int age,
    String? avatar,
  }) async {
    try {
      final response = await _dataSource.addChild(
        nickname: nickname,
        age: age,
        avatar: avatar,
      );
      if (response.isSuccess && response.data != null) {
        // 添加成功，更新本地儿童列表缓存
        final children = await getChildren();
        final childDTO = response.data!;
        final child = _convertChildDTO(childDTO);
        final updatedChildren = [...children, child];
        await StorageService.saveChildrenData(
          updatedChildren.map((c) => _userToMap(c)).toList(),
        );
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  /// 获取儿童列表
  /// 调用API获取当前用户的儿童列表
  Future<List<User>> getChildren() async {
    try {
      final response = await _dataSource.getChildren();
      if (response.isSuccess && response.data != null) {
        final children = response.data!
            .map((dto) => _convertChildDTO(dto))
            .toList();
        // 更新本地缓存
        await StorageService.saveChildrenData(
          children.map((c) => _userToMap(c)).toList(),
        );
        return children;
      }
      return [];
    } catch (e) {
      // API失败时尝试从本地缓存获取
      try {
        final childrenData = await StorageService.getChildrenData();
        return childrenData.map((data) {
          return User(
            id: data['id'] as String,
            type: UserType.child,
            parentId: data['parent_id'] as String?,
            nickname: data['nickname'] as String,
            age: int.tryParse(data['age']?.toString() ?? ''),
            avatar: data['avatar'] as String?,
            createdAt: DateTime.parse(data['created_at'] as String),
          );
        }).toList();
      } catch (_) {
        return [];
      }
    }
  }

  /// 获取家长设置
  /// 调用API获取家长控制设置
  Future<ParentSettings?> getParentSettings() async {
    try {
      final response = await _dataSource.getParentSettings();
      if (response.isSuccess && response.data != null) {
        final dto = response.data!;
        return ParentSettings(
          id: '',
          parentId: state.user?.id ?? '',
          dailyLimitMin: dto.dailyLimitMin,
          availableFrom: dto.availableFrom,
          availableTo: dto.availableTo,
          cameraAllowed: dto.cameraAllowed,
          locationAllowed: dto.locationAllowed,
          dataUploadCloud: dto.dataUploadCloud,
          createdAt: DateTime.now(),
        );
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// 更新家长设置
  /// 调用API更新家长控制设置
  Future<ParentSettings?> updateParentSettings(
      ParentSettingsDTO dto) async {
    try {
      final response = await _dataSource.updateParentSettings(dto);
      if (response.isSuccess && response.data != null) {
        final result = response.data!;
        return ParentSettings(
          id: '',
          parentId: state.user?.id ?? '',
          dailyLimitMin: result.dailyLimitMin,
          availableFrom: result.availableFrom,
          availableTo: result.availableTo,
          cameraAllowed: result.cameraAllowed,
          locationAllowed: result.locationAllowed,
          dataUploadCloud: result.dataUploadCloud,
          createdAt: DateTime.now(),
        );
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// 切换到儿童模式
  /// 本地切换当前用户为指定儿童，不调用API
  Future<void> switchToChild(String childId) async {
    try {
      final children = await getChildren();
      final child = children.where((c) => c.id == childId).firstOrNull;
      if (child == null) return;

      final token = await StorageService.getToken();
      if (token != null) {
        await _saveUserData(child);
        state = AuthState.authenticated(user: child, token: token);
      }
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  // ==================== 私有辅助方法 ====================

  /// 将UserDTO转换为本地User实体
  User _convertUserDTO(UserDTO dto) {
    return User(
      id: dto.id,
      type: dto.type == 'parent' ? UserType.parent : UserType.child,
      nickname: dto.nickname,
      avatar: dto.avatar,
      phone: dto.phone,
      status: dto.status == 'active'
          ? UserStatus.active
          : dto.status == 'inactive'
              ? UserStatus.inactive
              : UserStatus.deleted,
      createdAt: dto.createdAt.isNotEmpty
          ? DateTime.parse(dto.createdAt)
          : DateTime.now(),
    );
  }

  /// 将ChildDTO转换为本地User实体
  User _convertChildDTO(ChildDTO dto) {
    return User(
      id: dto.id,
      type: UserType.child,
      nickname: dto.nickname,
      avatar: dto.avatar,
      age: dto.age,
      status: dto.status == 'active'
          ? UserStatus.active
          : dto.status == 'inactive'
              ? UserStatus.inactive
              : UserStatus.deleted,
      createdAt: DateTime.now(),
    );
  }

  /// 保存认证数据到本地存储
  Future<void> _saveAuthData({
    required String token,
    required String refreshToken,
    required int expiresIn,
    required User user,
  }) async {
    await StorageService.saveToken(token);
    await StorageService.saveRefreshToken(refreshToken);
    await _saveUserData(user);
    // 保存token过期时间
    final expiryTime = DateTime.now()
        .add(Duration(seconds: expiresIn))
        .millisecondsSinceEpoch;
    await StorageService.saveSetting(
        AppConstants.tokenExpiryKey, expiryTime);
  }

  /// 保存用户数据到本地存储
  Future<void> _saveUserData(User user) async {
    await StorageService.saveUserData({
      'id': user.id,
      'type': user.type.value,
      'nickname': user.nickname,
      'age': user.age?.toString() ?? '',
      'avatar': user.avatar ?? '',
      'phone': user.phone ?? '',
      'status': user.status.value,
      'created_at': user.createdAt.toIso8601String(),
    });
  }

  /// 将User实体转换为Map
  Map<String, dynamic> _userToMap(User user) {
    return {
      'id': user.id,
      'type': user.type.value,
      'nickname': user.nickname,
      'age': user.age?.toString() ?? '',
      'avatar': user.avatar ?? '',
      'parent_id': user.parentId,
      'created_at': user.createdAt.toIso8601String(),
    };
  }
}

/// 登录表单状态
class LoginFormState {
  final String phone;
  final String password;
  final bool isLoading;
  final String? error;

  const LoginFormState({
    this.phone = '',
    this.password = '',
    this.isLoading = false,
    this.error,
  });

  LoginFormState copyWith({
    String? phone,
    String? password,
    bool? isLoading,
    String? error,
  }) {
    return LoginFormState(
      phone: phone ?? this.phone,
      password: password ?? this.password,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// 登录表单Provider
final loginFormProvider =
    StateNotifierProvider<LoginFormNotifier, LoginFormState>((ref) {
  return LoginFormNotifier();
});

/// 登录表单Notifier
class LoginFormNotifier extends StateNotifier<LoginFormState> {
  LoginFormNotifier() : super(const LoginFormState());

  /// 设置手机号
  void setPhone(String phone) {
    state = state.copyWith(phone: phone);
  }

  /// 设置密码
  void setPassword(String password) {
    state = state.copyWith(password: password);
  }

  /// 设置加载状态
  void setLoading(bool isLoading) {
    state = state.copyWith(isLoading: isLoading);
  }

  /// 设置错误信息
  void setError(String? error) {
    state = state.copyWith(error: error);
  }

  /// 清空表单
  void clear() {
    state = const LoginFormState();
  }
}
