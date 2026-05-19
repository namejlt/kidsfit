import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:permission_handler/permission_handler.dart';
import '../domain/entities/exercise_record.dart';
import 'pose_detection_service.dart';
import 'action_recognizer.dart';

/// 运动控制器状态
enum ExerciseControllerState {
  idle,
  initializing,
  ready,
  exercising,
  paused,
  completed,
  error,
}

/// 运动控制器
/// 负责管理运动训练的整个生命周期
class ExerciseController {
  final ExerciseType exerciseType;
  PoseDetectionService? _poseService;
  StreamSubscription? _skeletonSubscription;
  ActionRecognizer? _actionRecognizer;
  SkeletonModel? _currentSkeleton;

  final _stateController = StreamController<ExerciseControllerState>.broadcast();
  final _scoreController = StreamController<ActionScore>.broadcast();
  final _errorController = StreamController<String>.broadcast();

  ExerciseControllerState _state = ExerciseControllerState.idle;
  DateTime? _startTime;
  Timer? _scoreUpdateTimer;

  /// 状态流
  Stream<ExerciseControllerState> get stateStream => _stateController.stream;

  /// 评分流
  Stream<ActionScore> get scoreStream => _scoreController.stream;

  /// 错误流
  Stream<String> get errorStream => _errorController.stream;

  /// 当前状态
  ExerciseControllerState get state => _state;

  /// 是否正在运动
  bool get isExercising => _state == ExerciseControllerState.exercising;

  /// 构造函数
  ExerciseController({required this.exerciseType});

  /// 初始化控制器
  Future<bool> initialize() async {
    try {
      _updateState(ExerciseControllerState.initializing);

      // 请求摄像头权限
      final cameraStatus = await Permission.camera.request();
      if (!cameraStatus.isGranted) {
        _errorController.add('需要摄像头权限才能进行运动检测');
        _updateState(ExerciseControllerState.error);
        return false;
      }

      // 初始化骨骼检测服务
      _poseService = PoseDetectionService();
      final initialized = await _poseService!.initialize();
      if (!initialized) {
        _errorController.add('骨骼检测初始化失败');
        _updateState(ExerciseControllerState.error);
        return false;
      }

      // 订阅骨骼数据
      _skeletonSubscription = _poseService!.skeletonStream.listen(_onSkeletonDetected);

      _updateState(ExerciseControllerState.ready);
      return true;
    } catch (e) {
      _errorController.add('初始化失败: $e');
      _updateState(ExerciseControllerState.error);
      return false;
    }
  }

  /// 开始运动
  Future<void> startExercise() async {
    if (_state != ExerciseControllerState.ready && _state != ExerciseControllerState.paused) {
      return;
    }

    try {
      _startTime = DateTime.now();
      _actionRecognizer = ActionRecognizer();

      // 启动骨骼检测
      await _poseService?.startDetection();

      // 启动评分更新定时器
      _scoreUpdateTimer = Timer.periodic(
        const Duration(milliseconds: 100),
        (_) => _updateScore(),
      );

      _updateState(ExerciseControllerState.exercising);
    } catch (e) {
      _errorController.add('启动运动失败: $e');
    }
  }

  /// 暂停运动
  Future<void> pauseExercise() async {
    if (_state != ExerciseControllerState.exercising) return;

    await _poseService?.stopDetection();
    _scoreUpdateTimer?.cancel();
    _updateState(ExerciseControllerState.paused);
  }

  /// 恢复运动
  Future<void> resumeExercise() async {
    if (_state != ExerciseControllerState.paused) return;

    await _poseService?.startDetection();
    _scoreUpdateTimer = Timer.periodic(
      const Duration(milliseconds: 100),
      (_) => _updateScore(),
    );
    _updateState(ExerciseControllerState.exercising);
  }

  /// 停止运动
  Future<void> stopExercise() async {
    _scoreUpdateTimer?.cancel();
    await _poseService?.stopDetection();
    _updateState(ExerciseControllerState.completed);
  }

  /// 处理骨骼数据
  void _onSkeletonDetected(SkeletonModel skeleton) {
    _currentSkeleton = skeleton;
    _actionRecognizer?.updatePreviousSkeleton(skeleton);
  }

  /// 更新评分
  void _updateScore() {
    if (_actionRecognizer == null) return;

    ActionScore score;
    switch (exerciseType) {
      case ExerciseType.jumpRope:
        score = _actionRecognizer!.analyzeJumpRope();
        break;
      case ExerciseType.jumpingJack:
        score = _actionRecognizer!.analyzeJumpingJack();
        break;
      case ExerciseType.squat:
        score = _actionRecognizer!.analyzeSquat();
        break;
      case ExerciseType.sitUp:
        score = _actionRecognizer!.analyzeSitUp();
        break;
      case ExerciseType.highKnee:
        score = _actionRecognizer!.analyzeHighKnee();
        break;
      case ExerciseType.pushUp:
        score = _actionRecognizer!.analyzePushUp();
        break;
    }

    _scoreController.add(score);
    _actionRecognizer?.addScoreToHistory(score.totalScore);
  }

  /// 创建运动记录
  ExerciseRecord createExerciseRecord() {
    final duration = DateTime.now().difference(_startTime!).inSeconds;

    return ExerciseRecord(
      id: 'record_${DateTime.now().millisecondsSinceEpoch}',
      userId: 'current_user',
      type: exerciseType,
      durationSeconds: duration,
      count: _actionRecognizer?.analyzeJumpRope().count ?? 0,
      score: _actionRecognizer?.analyzeJumpRope().totalScore ?? 0,
      rhythmScore: _actionRecognizer?.analyzeJumpRope().rhythmScore ?? 0,
      amplitudeScore: _actionRecognizer?.analyzeJumpRope().amplitudeScore ?? 0,
      symmetryScore: _actionRecognizer?.analyzeJumpRope().symmetryScore ?? 0,
      continuityScore: _actionRecognizer?.analyzeJumpRope().continuityScore ?? 0,
      corrections: [],
      isOffline: false,
      startedAt: _startTime!,
      completedAt: DateTime.now(),
      createdAt: DateTime.now(),
    );
  }

  /// 更新状态
  void _updateState(ExerciseControllerState newState) {
    _state = newState;
    _stateController.add(newState);
  }

  /// 释放资源
  Future<void> dispose() async {
    _scoreUpdateTimer?.cancel();
    await _skeletonSubscription?.cancel();
    await _poseService?.dispose();
    await _stateController.close();
    await _scoreController.close();
    await _errorController.close();
  }
}

/// 运动控制器工厂
class ExerciseControllerFactory {
  /// 创建运动控制器
  static ExerciseController create(ExerciseType type) {
    return ExerciseController(exerciseType: type);
  }
}
