import 'dart:convert';
import 'package:hive_flutter/hive_flutter.dart';

/// 存储服务
/// 负责本地数据的安全存储，支持加密存储敏感数据
class StorageService {
  StorageService._();

  /// Box名称
  static const String _userBoxName = 'user_box';
  static const String _settingsBoxName = 'settings_box';
  static const String _exerciseBoxName = 'exercise_box';
  static const String _visionBoxName = 'vision_box';

  /// 存储键名
  static const String _userDataKey = 'user_data';
  static const String _tokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _childrenKey = 'children_data';

  /// Boxes
  static Box? _userBox;
  static Box? _settingsBox;
  static Box? _exerciseBox;
  static Box? _visionBox;

  /// 初始化存储服务
  static Future<void> initialize() async {
    // 打开加密的Boxes
    _userBox = await Hive.openBox(_userBoxName);
    _settingsBox = await Hive.openBox(_settingsBoxName);
    _exerciseBox = await Hive.openBox(_exerciseBoxName);
    _visionBox = await Hive.openBox(_visionBoxName);
  }

  // ==================== 用户数据操作 ====================

  /// 保存用户数据
  static Future<void> saveUserData(Map<String, dynamic> userData) async {
    await _userBox?.put(_userDataKey, jsonEncode(userData));
  }

  /// 获取用户数据
  static Future<Map<String, dynamic>?> getUserData() async {
    final data = _userBox?.get(_userDataKey);
    if (data == null) return null;
    return jsonDecode(data as String) as Map<String, dynamic>;
  }

  /// 保存Token
  static Future<void> saveToken(String token) async {
    await _userBox?.put(_tokenKey, token);
  }

  /// 获取Token
  static Future<String?> getToken() async {
    return _userBox?.get(_tokenKey) as String?;
  }

  /// 保存刷新Token
  static Future<void> saveRefreshToken(String token) async {
    await _userBox?.put(_refreshTokenKey, token);
  }

  /// 获取刷新Token
  static Future<String?> getRefreshToken() async {
    return _userBox?.get(_refreshTokenKey) as String?;
  }

  /// 保存儿童列表数据
  static Future<void> saveChildrenData(List<Map<String, dynamic>> children) async {
    await _userBox?.put(_childrenKey, jsonEncode(children));
  }

  /// 获取儿童列表数据
  static Future<List<Map<String, dynamic>>> getChildrenData() async {
    final data = _userBox?.get(_childrenKey);
    if (data == null) return [];
    final list = jsonDecode(data as String) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  // ==================== 运动数据操作 ====================

  /// 保存运动记录
  static Future<void> saveExerciseRecord(String id, Map<String, dynamic> record) async {
    await _exerciseBox?.put(id, jsonEncode(record));
  }

  /// 获取运动记录
  static Future<Map<String, dynamic>?> getExerciseRecord(String id) async {
    final data = _exerciseBox?.get(id);
    if (data == null) return null;
    return jsonDecode(data as String) as Map<String, dynamic>;
  }

  /// 获取所有运动记录
  static Future<List<Map<String, dynamic>>> getAllExerciseRecords() async {
    final records = <Map<String, dynamic>>[];
    for (final key in _exerciseBox?.keys ?? []) {
      final data = _exerciseBox?.get(key);
      if (data != null) {
        records.add(jsonDecode(data as String) as Map<String, dynamic>);
      }
    }
    return records;
  }

  /// 获取用户运动记录
  static Future<List<Map<String, dynamic>>> getExerciseRecordsByUser(String userId) async {
    final allRecords = await getAllExerciseRecords();
    return allRecords.where((r) => r['user_id'] == userId).toList();
  }

  /// 删除运动记录
  static Future<void> deleteExerciseRecord(String id) async {
    await _exerciseBox?.delete(id);
  }

  // ==================== 视力数据操作 ====================

  /// 保存视力记录
  static Future<void> saveVisionRecord(String id, Map<String, dynamic> record) async {
    await _visionBox?.put(id, jsonEncode(record));
  }

  /// 获取视力记录
  static Future<Map<String, dynamic>?> getVisionRecord(String id) async {
    final data = _visionBox?.get(id);
    if (data == null) return null;
    return jsonDecode(data as String) as Map<String, dynamic>;
  }

  /// 获取所有视力记录
  static Future<List<Map<String, dynamic>>> getAllVisionRecords() async {
    final records = <Map<String, dynamic>>[];
    for (final key in _visionBox?.keys ?? []) {
      final data = _visionBox?.get(key);
      if (data != null) {
        records.add(jsonDecode(data as String) as Map<String, dynamic>);
      }
    }
    return records;
  }

  /// 获取用户视力记录
  static Future<List<Map<String, dynamic>>> getVisionRecordsByChild(String childId) async {
    final allRecords = await getAllVisionRecords();
    return allRecords.where((r) => r['child_id'] == childId).toList();
  }

  // ==================== 设置操作 ====================

  /// 保存设置
  static Future<void> saveSetting(String key, dynamic value) async {
    await _settingsBox?.put(key, value);
  }

  /// 获取设置
  static Future<T?> getSetting<T>(String key) async {
    return _settingsBox?.get(key) as T?;
  }

  /// 删除设置
  static Future<void> deleteSetting(String key) async {
    await _settingsBox?.delete(key);
  }

  // ==================== 通用操作 ====================

  /// 清除所有数据
  static Future<void> clearAll() async {
    await _userBox?.clear();
    await _exerciseBox?.clear();
    await _visionBox?.clear();
  }

  /// 清除用户数据
  static Future<void> clearUserData() async {
    await _userBox?.clear();
  }

  /// 关闭所有Boxes
  static Future<void> close() async {
    await _userBox?.close();
    await _settingsBox?.close();
    await _exerciseBox?.close();
    await _visionBox?.close();
  }
}
