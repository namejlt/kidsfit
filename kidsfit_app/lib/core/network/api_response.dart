/// 统一API响应模型
/// 匹配后端契约：{"code": 0, "message": "success", "data": {...}}
class ApiResponse<T> {
  /// 响应码，0表示成功
  final int code;

  /// 响应消息
  final String message;

  /// 响应数据
  final T? data;

  const ApiResponse({
    required this.code,
    required this.message,
    this.data,
  });

  /// 请求是否成功
  bool get isSuccess => code == 0;

  /// 从JSON创建ApiResponse，需要传入fromJsonT解析data字段
  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic) fromJsonT,
  ) {
    return ApiResponse<T>(
      code: json['code'] as int? ?? -1,
      message: json['message'] as String? ?? '',
      data: json['data'] != null ? fromJsonT(json['data']) : null,
    );
  }

  @override
  String toString() => 'ApiResponse(code: $code, message: $message, data: $data)';
}

/// 分页信息模型
/// 匹配后端契约：{"page":1, "page_size":20, "total":100, "total_pages":5}
class Pagination {
  /// 当前页码
  final int page;

  /// 每页条数
  final int pageSize;

  /// 总记录数
  final int total;

  /// 总页数
  final int totalPages;

  const Pagination({
    required this.page,
    required this.pageSize,
    required this.total,
    required this.totalPages,
  });

  /// 从JSON创建Pagination
  factory Pagination.fromJson(Map<String, dynamic> json) {
    return Pagination(
      page: json['page'] as int? ?? 1,
      pageSize: json['page_size'] as int? ?? 20,
      total: json['total'] as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 0,
    );
  }

  /// 转换为JSON
  Map<String, dynamic> toJson() => {
        'page': page,
        'page_size': pageSize,
        'total': total,
        'total_pages': totalPages,
      };

  /// 是否有下一页
  bool get hasNextPage => page < totalPages;

  @override
  String toString() =>
      'Pagination(page: $page, pageSize: $pageSize, total: $total, totalPages: $totalPages)';
}

/// 分页数据模型
/// 匹配后端契约：{"list": [...], "pagination": {...}}
class PaginatedData<T> {
  /// 数据列表
  final List<T> list;

  /// 分页信息
  final Pagination pagination;

  const PaginatedData({
    required this.list,
    required this.pagination,
  });

  /// 是否为空
  bool get isEmpty => list.isEmpty;

  /// 是否非空
  bool get isNotEmpty => list.isNotEmpty;

  /// 从JSON创建PaginatedData，需要传入fromJsonT解析列表项
  factory PaginatedData.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic) fromJsonT,
  ) {
    final listRaw = json['list'] as List<dynamic>? ?? [];
    return PaginatedData<T>(
      list: listRaw.map((e) => fromJsonT(e)).toList(),
      pagination: json['pagination'] != null
          ? Pagination.fromJson(json['pagination'] as Map<String, dynamic>)
          : const Pagination(page: 1, pageSize: 20, total: 0, totalPages: 0),
    );
  }

  @override
  String toString() =>
      'PaginatedData(list: ${list.length} items, pagination: $pagination)';
}
