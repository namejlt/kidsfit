import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/user.dart';
import '../../services/storage_service.dart';
import '../../core/router/app_router.dart';

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

  const AuthState({
    this.status = AuthStatus.unknown,
    this.user,
    this.token,
    this.refreshToken,
    this.error,
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
    return const AuthState(status: AuthStatus.unknown);
  }

  AuthState copyWith({
    AuthStatus? status,
    User? user,
    String? token,
    String? refreshToken,
    String? error,
  }) {
    return AuthState(
      status: status ?? this.status,
      user: user ?? this.user,
      token: token ?? this.token,
      refreshToken: refreshToken ?? this.refreshToken,
      error: error ?? this.error,
    );
  }
}

/// 认证状态Provider
final authStateProvider =
    StateNotifierProvider<AuthStateNotifier, AuthState>((ref) {
  return AuthStateNotifier();
});

/// 认证状态Notifier
class AuthStateNotifier extends StateNotifier<AuthState> {
  AuthStateNotifier() : super(AuthState.loading()) {
    _checkAuthStatus();
  }

  /// 检查认证状态
  Future<void> _checkAuthStatus() async {
    try {
      final userData = await StorageService.getUserData();
      if (userData != null) {
        final user = User(
          id: userData['id'] as String,
          type: userData['type'] == 'parent' ? UserType.parent : UserType.child,
          nickname: userData['nickname'] as String,
          age: userData['age'] as int?,
          avatar: userData['avatar'] as String?,
          phone: userData['phone'] as String?,
          createdAt: DateTime.parse(userData['created_at'] as String),
        );
        final token = await StorageService.getToken();
        if (token != null) {
          state = AuthState.authenticated(user: user, token: token);
          return;
        }
      }
      state = AuthState.unauthenticated();
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
    }
  }

  /// 登录
  Future<bool> login({
    required String phone,
    required String password,
  }) async {
    try {
      state = AuthState.loading();

      // TODO: 调用API进行登录
      // 模拟登录成功
      await Future.delayed(const Duration(seconds: 1));

      final user = User(
        id: 'user_${DateTime.now().millisecondsSinceEpoch}',
        type: UserType.parent,
        nickname: '家长用户',
        phone: phone,
        createdAt: DateTime.now(),
      );

      const token = 'mock_jwt_token';

      // 保存用户数据
      await StorageService.saveUserData({
        'id': user.id,
        'type': user.type.value,
        'nickname': user.nickname,
        'phone': user.phone,
        'created_at': user.createdAt.toIso8601String(),
      });
      await StorageService.saveToken(token);

      state = AuthState.authenticated(user: user, token: token);
      return true;
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
      return false;
    }
  }

  /// 注册
  Future<bool> register({
    required String phone,
    required String password,
    required String nickname,
    required bool isParent,
    int? childAge,
  }) async {
    try {
      state = AuthState.loading();

      // TODO: 调用API进行注册
      await Future.delayed(const Duration(seconds: 1));

      final user = User(
        id: 'user_${DateTime.now().millisecondsSinceEpoch}',
        type: isParent ? UserType.parent : UserType.child,
        nickname: nickname,
        age: childAge,
        phone: isParent ? phone : null,
        createdAt: DateTime.now(),
      );

      const token = 'mock_jwt_token';

      await StorageService.saveUserData({
        'id': user.id,
        'type': user.type.value,
        'nickname': user.nickname,
        'age': user.age?.toString() ?? '',
        'phone': user.phone ?? '',
        'created_at': user.createdAt.toIso8601String(),
      });
      await StorageService.saveToken(token);

      state = AuthState.authenticated(user: user, token: token);
      return true;
    } catch (e) {
      state = AuthState.unauthenticated(error: e.toString());
      return false;
    }
  }

  /// 登出
  Future<void> logout() async {
    await StorageService.clearAll();
    state = AuthState.unauthenticated();
  }

  /// 刷新Token
  Future<bool> refreshToken() async {
    try {
      // TODO: 调用API刷新Token
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 更新用户信息
  Future<void> updateUser(User user) async {
    await StorageService.saveUserData({
      'id': user.id,
      'type': user.type.value,
      'nickname': user.nickname,
      'age': user.age?.toString() ?? '',
      'avatar': user.avatar ?? '',
      'phone': user.phone ?? '',
      'created_at': user.createdAt.toIso8601String(),
    });
    state = state.copyWith(user: user);
  }

  /// 添加儿童用户
  Future<bool> addChild({
    required String nickname,
    required int age,
  }) async {
    try {
      // TODO: 调用API添加儿童
      await Future.delayed(const Duration(milliseconds: 500));

      final child = User(
        id: 'child_${DateTime.now().millisecondsSinceEpoch}',
        type: UserType.child,
        parentId: state.user?.id,
        nickname: nickname,
        age: age,
        createdAt: DateTime.now(),
      );

      // 保存儿童数据到本地
      final children = await StorageService.getChildrenData();
      final updatedChildren = [
        ...children,
        {
          'id': child.id,
          'type': child.type.value,
          'nickname': child.nickname,
          'age': child.age.toString(),
          'parent_id': child.parentId,
          'created_at': child.createdAt.toIso8601String(),
        }
      ];
      await StorageService.saveChildrenData(updatedChildren);

      return true;
    } catch (e) {
      return false;
    }
  }

  /// 获取儿童列表
  Future<List<User>> getChildren() async {
    try {
      final childrenData = await StorageService.getChildrenData();
      return childrenData.map((data) {
        return User(
          id: data['id'] as String,
          type: UserType.child,
          parentId: data['parent_id'] as String?,
          nickname: data['nickname'] as String,
          age: int.tryParse(data['age']?.toString() ?? ''),
          createdAt: DateTime.parse(data['created_at'] as String),
        );
      }).toList();
    } catch (e) {
      return [];
    }
  }

  /// 切换当前儿童用户
  Future<void> switchToChild(User child) async {
    final token = await StorageService.getToken();
    if (token != null) {
      await StorageService.saveUserData({
        'id': child.id,
        'type': child.type.value,
        'nickname': child.nickname,
        'age': child.age?.toString() ?? '',
        'avatar': child.avatar ?? '',
        'parent_id': child.parentId,
        'created_at': child.createdAt.toIso8601String(),
      });
      state = AuthState.authenticated(user: child, token: token);
    }
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

  void setPhone(String phone) {
    state = state.copyWith(phone: phone);
  }

  void setPassword(String password) {
    state = state.copyWith(password: password);
  }

  void setLoading(bool isLoading) {
    state = state.copyWith(isLoading: isLoading);
  }

  void setError(String? error) {
    state = state.copyWith(error: error);
  }

  void clear() {
    state = const LoginFormState();
  }
}
