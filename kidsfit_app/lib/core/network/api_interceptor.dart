import 'dart:async';
import 'package:dio/dio.dart';
import '../../services/storage_service.dart';
import '../constants/app_constants.dart';

/// JWT认证拦截器
/// 自动注入Authorization Bearer Token，401时自动刷新Token并重试原请求
class AuthInterceptor extends Interceptor {
  /// 是否正在刷新Token
  bool _isRefreshing = false;

  /// 等待刷新完成的Completer
  Completer<void>? _refreshCompleter;

  @override
  void onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    /// 自动注入Authorization Bearer Token
    final token = await StorageService.getToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    /// 仅处理401未授权错误
    if (err.response?.statusCode != 401) {
      handler.next(err);
      return;
    }

    /// 如果是刷新Token接口本身返回401，直接失败
    if (err.requestOptions.path == AppConstants.authRefresh) {
      await _handleRefreshFailure();
      handler.next(err);
      return;
    }

    /// 尝试刷新Token
    try {
      final newToken = await _refreshToken();
      if (newToken != null) {
        /// 刷新成功，重试原请求
        final options = err.requestOptions;
        options.headers['Authorization'] = 'Bearer $newToken';
        final response = await Dio().fetch(options);
        handler.resolve(response);
      } else {
        /// 刷新失败，清除本地Token
        await _handleRefreshFailure();
        handler.next(err);
      }
    } catch (e) {
      /// 刷新异常，清除本地Token
      await _handleRefreshFailure();
      handler.next(err);
    }
  }

  /// 刷新Token
  /// 使用refresh_token调用/api/v1/auth/refresh接口
  Future<String?> _refreshToken() async {
    /// 如果已经在刷新中，等待刷新完成
    if (_isRefreshing) {
      await _refreshCompleter?.future;
      return StorageService.getToken();
    }

    _isRefreshing = true;
    _refreshCompleter = Completer<void>();

    try {
      final refreshToken = await StorageService.getRefreshToken();
      if (refreshToken == null || refreshToken.isEmpty) {
        return null;
      }

      final dio = Dio();
      final response = await dio.post(
        '${AppConstants.currentApiBaseUrl}${AppConstants.authRefresh}',
        data: {'refresh_token': refreshToken},
        options: Options(
          headers: {'Content-Type': 'application/json'},
        ),
      );

      final data = response.data as Map<String, dynamic>;
      final code = data['code'] as int? ?? -1;
      if (code == 0 && data['data'] != null) {
        final tokenData = data['data'] as Map<String, dynamic>;
        final newAccessToken = tokenData['access_token'] as String?;
        final newRefreshToken = tokenData['refresh_token'] as String?;
        final expiresIn = tokenData['expires_in'] as int?;

        if (newAccessToken != null) {
          /// 保存新Token
          await StorageService.saveToken(newAccessToken);
          if (newRefreshToken != null) {
            await StorageService.saveRefreshToken(newRefreshToken);
          }
          if (expiresIn != null) {
            final expiryTime = DateTime.now()
                .add(Duration(seconds: expiresIn))
                .millisecondsSinceEpoch;
            await StorageService.saveSetting(
                AppConstants.tokenExpiryKey, expiryTime);
          }
          return newAccessToken;
        }
      }

      return null;
    } catch (e) {
      return null;
    } finally {
      _isRefreshing = false;
      _refreshCompleter?.complete();
      _refreshCompleter = null;
    }
  }

  /// 处理Token刷新失败
  /// 清除本地Token，通知退出登录
  Future<void> _handleRefreshFailure() async {
    await StorageService.clearUserData();
  }
}
