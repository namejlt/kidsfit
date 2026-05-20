import 'package:dio/dio.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';
import '../constants/app_constants.dart';
import 'api_interceptor.dart';
import 'api_response.dart';

/// API客户端
/// 基于Dio的单例封装，统一处理请求、响应和错误
class ApiClient {
  /// 单例实例
  static final ApiClient _instance = ApiClient._internal();

  /// 获取单例实例
  static ApiClient get instance => _instance;

  /// Dio实例
  late final Dio _dio;

  /// 私有构造函数
  ApiClient._internal() {
    _dio = Dio(
      BaseOptions(
        baseUrl: AppConstants.currentApiBaseUrl,
        connectTimeout:
            const Duration(milliseconds: AppConstants.connectionTimeout),
        receiveTimeout:
            const Duration(milliseconds: AppConstants.receiveTimeout),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    /// 添加拦截器：认证 + 日志
    _dio.interceptors.add(AuthInterceptor());
    _dio.interceptors.add(PrettyDioLogger(
      requestHeader: true,
      requestBody: true,
      responseHeader: false,
      responseBody: true,
      error: true,
      compact: true,
    ));
  }

  /// GET请求
  /// 返回统一响应ApiResponse<T>，需要传入fromJsonT解析data字段
  Future<ApiResponse<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJsonT,
    Options? options,
  }) async {
    try {
      final response = await _dio.get(
        path,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(
        response.data as Map<String, dynamic>,
        fromJsonT,
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<T>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// POST请求
  /// 返回统一响应ApiResponse<T>，需要传入fromJsonT解析data字段
  Future<ApiResponse<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJsonT,
    Options? options,
  }) async {
    try {
      final response = await _dio.post(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(
        response.data as Map<String, dynamic>,
        fromJsonT,
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<T>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// PUT请求
  /// 返回统一响应ApiResponse<T>，需要传入fromJsonT解析data字段
  Future<ApiResponse<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJsonT,
    Options? options,
  }) async {
    try {
      final response = await _dio.put(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(
        response.data as Map<String, dynamic>,
        fromJsonT,
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<T>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// DELETE请求
  /// 返回统一响应ApiResponse<T>，需要传入fromJsonT解析data字段
  Future<ApiResponse<T>> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJsonT,
    Options? options,
  }) async {
    try {
      final response = await _dio.delete(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(
        response.data as Map<String, dynamic>,
        fromJsonT,
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<T>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// GET请求（分页）
  /// 返回分页数据PaginatedData<T>
  Future<ApiResponse<PaginatedData<T>>> getPaged<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJsonT,
    Options? options,
  }) async {
    try {
      final response = await _dio.get(
        path,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(
        response.data as Map<String, dynamic>,
        (dynamic data) => PaginatedData.fromJson(
          data as Map<String, dynamic>,
          fromJsonT,
        ),
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<PaginatedData<T>>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// POST请求（无返回数据解析）
  /// 用于不需要解析data字段的场景，如logout
  Future<ApiResponse<void>> postVoid(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) async {
    try {
      final response = await _dio.post(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      final json = response.data as Map<String, dynamic>;
      return ApiResponse<void>(
        code: json['code'] as int? ?? -1,
        message: json['message'] as String? ?? '',
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<void>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// PUT请求（无返回数据解析）
  Future<ApiResponse<void>> putVoid(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) async {
    try {
      final response = await _dio.put(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      final json = response.data as Map<String, dynamic>;
      return ApiResponse<void>(
        code: json['code'] as int? ?? -1,
        message: json['message'] as String? ?? '',
      );
    } on DioException catch (e) {
      return _handleDioError(e);
    } catch (e) {
      return ApiResponse<void>(
        code: -1,
        message: '未知错误: $e',
      );
    }
  }

  /// 处理Dio错误
  /// 统一处理网络错误、超时、服务端错误码
  ApiResponse<T> _handleDioError<T>(DioException e) {
    String message;
    int code;

    switch (e.type) {
      case DioExceptionType.connectionTimeout:
        message = '连接超时，请检查网络';
        code = -2;
        break;
      case DioExceptionType.sendTimeout:
        message = '请求发送超时，请检查网络';
        code = -3;
        break;
      case DioExceptionType.receiveTimeout:
        message = '响应超时，请检查网络';
        code = -4;
        break;
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        if (statusCode == 401) {
          message = '认证失败，请重新登录';
          code = 401;
        } else if (statusCode == 403) {
          message = '无权限访问';
          code = 403;
        } else if (statusCode == 404) {
          message = '请求的资源不存在';
          code = 404;
        } else if (statusCode != null && statusCode >= 500) {
          message = '服务器错误，请稍后重试';
          code = statusCode;
        } else {
          /// 尝试从响应体中获取后端返回的错误信息
          try {
            final data = e.response?.data as Map<String, dynamic>?;
            message = data?['message'] as String? ?? '请求失败';
            code = data?['code'] as int? ?? (statusCode ?? -5);
          } catch (_) {
            message = '请求失败';
            code = statusCode ?? -5;
          }
        }
        break;
      case DioExceptionType.cancel:
        message = '请求已取消';
        code = -6;
        break;
      case DioExceptionType.connectionError:
        message = '网络连接失败，请检查网络设置';
        code = -7;
        break;
      default:
        message = '网络异常: ${e.message}';
        code = -8;
    }

    return ApiResponse<T>(code: code, message: message);
  }
}
