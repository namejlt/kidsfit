import 'dart:math';
import 'package:google_mlkit_pose_detection/google_mlkit_pose_detection.dart';
import 'pose_detection_service.dart';

/// 动作评分结果
class ActionScore {
  final int rhythmScore;
  final int amplitudeScore;
  final int symmetryScore;
  final int continuityScore;
  final int totalScore;
  final String? correction;
  final int count;

  ActionScore({
    required this.rhythmScore,
    required this.amplitudeScore,
    required this.symmetryScore,
    required this.continuityScore,
    required this.totalScore,
    this.correction,
    this.count = 0,
  });
}

/// 动作识别器
/// 基于骨骼关键点检测结果进行动作识别和评分
class ActionRecognizer {
  SkeletonModel? _previousSkeleton;
  DateTime? _lastActionTime;
  int _actionCount = 0;
  List<int> _recentScores = [];

  final SkeletonModel? currentSkeleton;

  ActionRecognizer({this.currentSkeleton});

  /// 分析跳绳动作
  ActionScore analyzeJumpRope() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;

    // 获取关键骨骼点
    final leftWrist = skeleton.getPoint(PoseLandmarkType.leftWrist);
    final rightWrist = skeleton.getPoint(PoseLandmarkType.rightWrist);
    final leftShoulder = skeleton.getPoint(PoseLandmarkType.leftShoulder);
    final rightShoulder = skeleton.getPoint(PoseLandmarkType.rightShoulder);
    final leftAnkle = skeleton.getPoint(PoseLandmarkType.leftAnkle);
    final rightAnkle = skeleton.getPoint(PoseLandmarkType.rightAnkle);

    if (leftWrist == null || rightWrist == null) {
      return _createCurrentScore(null);
    }

    // 检测跳跃动作
    bool isJumping = false;
    if (leftAnkle != null && rightAnkle != null) {
      final jumpHeight = (leftAnkle.position.dy + rightAnkle.position.dy) / 2;
      final shoulderHeight = ((leftShoulder?.position.dy ?? 0) + (rightShoulder?.position.dy ?? 0)) / 2;
      isJumping = jumpHeight < shoulderHeight - 100;
    }

    // 检测手臂动作（绳子位置）
    final wristAvgY = (leftWrist.position.dy + rightWrist.position.dy) / 2;
    final shoulderAvgY = ((leftShoulder?.position.dy ?? 0) + (rightShoulder?.position.dy ?? 0)) / 2;
    final wristAboveShoulder = wristAvgY < shoulderAvgY - 50;

    // 计数逻辑：脚离地 + 手在腰部以上
    if (isJumping && wristAboveShoulder) {
      final now = DateTime.now();
      if (_lastActionTime == null || now.difference(_lastActionTime!).inMilliseconds > 300) {
        _actionCount++;
        _lastActionTime = now;
      }
    }

    // 计算各项评分
    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = _calculateAmplitudeScore(isJumping, leftAnkle, rightAnkle);
    final symmetryScore = _calculateSymmetryScore(leftWrist, rightWrist, leftShoulder, rightShoulder);
    final continuityScore = _calculateContinuityScore();

    // 生成纠正建议
    String? correction = _generateJumpRopeCorrection(rhythmScore, amplitudeScore, symmetryScore);

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: correction,
      count: _actionCount,
    );
  }

  /// 分析开合跳
  ActionScore analyzeJumpingJack() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;
    final leftShoulder = skeleton.getPoint(PoseLandmarkType.leftShoulder);
    final rightShoulder = skeleton.getPoint(PoseLandmarkType.rightShoulder);
    final leftWrist = skeleton.getPoint(PoseLandmarkType.leftWrist);
    final rightWrist = skeleton.getPoint(PoseLandmarkType.rightWrist);
    final leftHip = skeleton.getPoint(PoseLandmarkType.leftHip);
    final rightHip = skeleton.getPoint(PoseLandmarkType.rightHip);
    final leftAnkle = skeleton.getPoint(PoseLandmarkType.leftAnkle);
    final rightAnkle = skeleton.getPoint(PoseLandmarkType.rightAnkle);

    // 计算手臂张开角度
    double armAngle = 0;
    if (leftShoulder != null && rightShoulder != null && leftWrist != null && rightWrist != null) {
      final shoulderWidth = (rightShoulder.position - leftShoulder.position).distance;
      final armLength = (leftShoulder.position - leftWrist.position).distance +
                        (rightShoulder.position - rightWrist.position).distance / 2;
      armAngle = acos(1 - (shoulderWidth / (armLength * 2)).clamp(0, 1)) * 180 / pi;
    }

    // 计算腿部张开角度
    double legAngle = 0;
    if (leftHip != null && rightHip != null && leftAnkle != null && rightAnkle != null) {
      final hipWidth = (rightHip.position - leftHip.position).distance;
      final legLength = (leftHip.position - leftAnkle.position).distance +
                        (rightHip.position - rightAnkle.position).distance / 2;
      legAngle = acos(1 - (hipWidth / (legLength * 2)).clamp(0, 1)) * 180 / pi;
    }

    // 检测动作完成状态
    bool isFullyOpen = armAngle > 120 && legAngle > 90;
    bool isFullyClosed = armAngle < 60 && legAngle < 60;

    // 计数：完成一次开合
    if (isFullyOpen && _previousSkeleton != null) {
      final prevLeftAnkle = _previousSkeleton!.getPoint(PoseLandmarkType.leftAnkle);
      final prevRightAnkle = _previousSkeleton!.getPoint(PoseLandmarkType.rightAnkle);
      if (prevLeftAnkle != null && prevRightAnkle != null &&
          (prevLeftAnkle.position - leftAnkle!.position).distance > 30) {
        _actionCount++;
      }
    }

    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = ((armAngle / 180) * 50 + (legAngle / 180) * 50).round().clamp(0, 100);
    final symmetryScore = _calculateSymmetryScore(leftWrist, rightWrist, leftShoulder, rightShoulder);
    final continuityScore = _calculateContinuityScore();

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: _generateJumpingJackCorrection(amplitudeScore, symmetryScore),
      count: _actionCount,
    );
  }

  /// 分析深蹲
  ActionScore analyzeSquat() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;
    final leftHip = skeleton.getPoint(PoseLandmarkType.leftHip);
    final rightHip = skeleton.getPoint(PoseLandmarkType.rightHip);
    final leftKnee = skeleton.getPoint(PoseLandmarkType.leftKnee);
    final rightKnee = skeleton.getPoint(PoseLandmarkType.rightKnee);
    final leftAnkle = skeleton.getPoint(PoseLandmarkType.leftAnkle);
    final rightAnkle = skeleton.getPoint(PoseLandmarkType.rightAnkle);
    final leftShoulder = skeleton.getPoint(PoseLandmarkType.leftShoulder);
    final rightShoulder = skeleton.getPoint(PoseLandmarkType.rightShoulder);

    // 计算膝盖角度
    double leftKneeAngle = 0;
    double rightKneeAngle = 0;

    if (leftHip != null && leftKnee != null && leftAnkle != null) {
      leftKneeAngle = skeleton.calculateAngle(
        PoseLandmarkType.leftHip,
        PoseLandmarkType.leftKnee,
        PoseLandmarkType.leftAnkle,
      );
    }

    if (rightHip != null && rightKnee != null && rightAnkle != null) {
      rightKneeAngle = skeleton.calculateAngle(
        PoseLandmarkType.rightHip,
        PoseLandmarkType.rightKnee,
        PoseLandmarkType.rightAnkle,
      );
    }

    final avgKneeAngle = (leftKneeAngle + rightKneeAngle) / 2;

    // 检测下蹲完成
    bool isSquatDown = avgKneeAngle < 100;
    bool isStanding = avgKneeAngle > 160;

    // 计数：完成一次深蹲
    if (isSquatDown && _previousSkeleton != null) {
      final prevKneeAngle = _previousSkeleton!.calculateAngle(
        PoseLandmarkType.leftHip,
        PoseLandmarkType.leftKnee,
        PoseLandmarkType.leftAnkle,
      );
      if (prevKneeAngle > 150 && avgKneeAngle < 100) {
        _actionCount++;
      }
    }

    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = _calculateSquatAmplitudeScore(avgKneeAngle);
    final symmetryScore = _calculateLegSymmetryScore(leftKneeAngle, rightKneeAngle);
    final continuityScore = _calculateContinuityScore();

    // 检查膝盖是否内扣
    String? correction;
    if (symmetryScore < 70) {
      correction = '注意膝盖不要内扣';
    } else if (amplitudeScore < 70) {
      correction = '下蹲幅度再大一些';
    }

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: correction,
      count: _actionCount,
    );
  }

  /// 分析仰卧起坐
  ActionScore analyzeSitUp() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;
    final nose = skeleton.getPoint(PoseLandmarkType.nose);
    final leftShoulder = skeleton.getPoint(PoseLandmarkType.leftShoulder);
    final rightShoulder = skeleton.getPoint(PoseLandmarkType.rightShoulder);
    final leftHip = skeleton.getPoint(PoseLandmarkType.leftHip);
    final rightHip = skeleton.getPoint(PoseLandmarkType.rightHip);

    // 计算躯干抬起角度
    double trunkAngle = 0;
    if (nose != null && leftShoulder != null && rightShoulder != null &&
        leftHip != null && rightHip != null) {
      final shoulderMidY = (leftShoulder.position.dy + rightShoulder.position.dy) / 2;
      final hipMidY = (leftHip.position.dy + rightHip.position.dy) / 2;
      final shoulderMidX = (leftShoulder.position.dx + rightShoulder.position.dx) / 2;
      final hipMidX = (leftHip.position.dx + rightHip.position.dx) / 2;

      trunkAngle = atan2(hipMidY - shoulderMidY, hipMidX - shoulderMidX) * 180 / pi;
      trunkAngle = trunkAngle.abs();
    }

    // 检测完成状态
    bool isUp = trunkAngle > 60;
    bool isDown = trunkAngle < 30;

    // 计数
    if (isUp && _previousSkeleton != null) {
      final prevNose = _previousSkeleton!.getPoint(PoseLandmarkType.nose);
      if (prevNose != null && nose != null) {
        final rise = nose.position.dy - prevNose.position.dy;
        if (rise.abs() > 20) {
          _actionCount++;
        }
      }
    }

    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = (trunkAngle / 90 * 100).round().clamp(0, 100);
    final symmetryScore = _calculateSymmetryScore(
      leftShoulder, rightShoulder, leftHip, rightHip
    );
    final continuityScore = _calculateContinuityScore();

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: amplitudeScore < 70 ? '请将身体抬得更高' : null,
      count: _actionCount,
    );
  }

  /// 分析高抬腿
  ActionScore analyzeHighKnee() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;
    final leftHip = skeleton.getPoint(PoseLandmarkType.leftHip);
    final rightHip = skeleton.getPoint(PoseLandmarkType.rightHip);
    final leftKnee = skeleton.getPoint(PoseLandmarkType.leftKnee);
    final rightKnee = skeleton.getPoint(PoseLandmarkType.rightKnee);
    final leftAnkle = skeleton.getPoint(PoseLandmarkType.leftAnkle);
    final rightAnkle = skeleton.getPoint(PoseLandmarkType.rightAnkle);

    // 计算抬腿高度
    double leftKneeHeight = 0;
    double rightKneeHeight = 0;

    if (leftHip != null && leftKnee != null) {
      leftKneeHeight = leftHip!.position.dy - leftKnee!.position.dy;
    }
    if (rightHip != null && rightKnee != null) {
      rightKneeHeight = rightHip!.position.dy - rightKnee!.position.dy;
    }

    final avgKneeHeight = (leftKneeHeight + rightKneeHeight) / 2;

    // 计数：交替抬腿
    bool isLeftUp = leftKneeHeight > 50;
    bool isRightUp = rightKneeHeight > 50;

    if ((isLeftUp || isRightUp) && _previousSkeleton != null) {
      final prevLeftKnee = _previousSkeleton!.getPoint(PoseLandmarkType.leftKnee);
      final prevRightKnee = _previousSkeleton!.getPoint(PoseLandmarkType.rightKnee);

      if (prevLeftKnee != null && prevRightKnee != null) {
        if ((leftKneeHeight > 50 && prevLeftKnee.position.dy > leftKnee!.position.dy) ||
            (rightKneeHeight > 50 && prevRightKnee.position.dy > rightKnee!.position.dy)) {
          _actionCount++;
        }
      }
    }

    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = (avgKneeHeight / 100 * 100).round().clamp(0, 100);
    final symmetryScore = _calculateLegSymmetryScore(leftKneeHeight, rightKneeHeight);
    final continuityScore = _calculateContinuityScore();

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: amplitudeScore < 60 ? '请将膝盖抬得更高' : null,
      count: _actionCount,
    );
  }

  /// 分析俯卧撑
  ActionScore analyzePushUp() {
    if (currentSkeleton == null) {
      return ActionScore(
        rhythmScore: 0,
        amplitudeScore: 0,
        symmetryScore: 0,
        continuityScore: 0,
        totalScore: 0,
      );
    }

    final skeleton = currentSkeleton!;
    final leftShoulder = skeleton.getPoint(PoseLandmarkType.leftShoulder);
    final rightShoulder = skeleton.getPoint(PoseLandmarkType.rightShoulder);
    final leftElbow = skeleton.getPoint(PoseLandmarkType.leftElbow);
    final rightElbow = skeleton.getPoint(PoseLandmarkType.rightElbow);
    final leftWrist = skeleton.getPoint(PoseLandmarkType.leftWrist);
    final rightWrist = skeleton.getPoint(PoseLandmarkType.rightWrist);
    final leftHip = skeleton.getPoint(PoseLandmarkType.leftHip);
    final rightHip = skeleton.getPoint(PoseLandmarkType.rightHip);

    // 计算手臂角度
    double leftArmAngle = 0;
    double rightArmAngle = 0;

    if (leftShoulder != null && leftElbow != null && leftWrist != null) {
      leftArmAngle = skeleton.calculateAngle(
        PoseLandmarkType.leftShoulder,
        PoseLandmarkType.leftElbow,
        PoseLandmarkType.leftWrist,
      );
    }

    if (rightShoulder != null && rightElbow != null && rightWrist != null) {
      rightArmAngle = skeleton.calculateAngle(
        PoseLandmarkType.rightShoulder,
        PoseLandmarkType.rightElbow,
        PoseLandmarkType.rightWrist,
      );
    }

    final avgArmAngle = (leftArmAngle + rightArmAngle) / 2;

    // 检测身体是否平直
    double bodyAlignment = 100;
    if (leftShoulder != null && rightShoulder != null &&
        leftHip != null && rightHip != null) {
      final shoulderMidY = (leftShoulder.position.dy + rightShoulder.position.dy) / 2;
      final hipMidY = (leftHip.position.dy + rightHip.position.dy) / 2;
      bodyAlignment = 100 - ((shoulderMidY - hipMidY).abs() / 5);
      bodyAlignment = bodyAlignment.clamp(0, 100);
    }

    // 计数
    bool isDown = avgArmAngle < 100;
    if (isDown && _previousSkeleton != null) {
      final prevAngle = _previousSkeleton!.calculateAngle(
        PoseLandmarkType.leftShoulder,
        PoseLandmarkType.leftElbow,
        PoseLandmarkType.leftWrist,
      );
      if (prevAngle > 150 && avgArmAngle < 100) {
        _actionCount++;
      }
    }

    final rhythmScore = _calculateRhythmScore();
    final amplitudeScore = avgArmAngle < 90 ? 100 : ((180 - avgArmAngle) / 90 * 100).round();
    final symmetryScore = _calculateArmSymmetryScore(leftArmAngle, rightArmAngle);
    final continuityScore = (bodyAlignment * 0.5 + _calculateContinuityScore() * 0.5).round();

    String? correction;
    if (bodyAlignment < 80) {
      correction = '请保持身体平直';
    } else if (symmetryScore < 70) {
      correction = '请保持双臂用力均匀';
    }

    return ActionScore(
      rhythmScore: rhythmScore,
      amplitudeScore: amplitudeScore,
      symmetryScore: symmetryScore,
      continuityScore: continuityScore,
      totalScore: _calculateTotalScore(rhythmScore, amplitudeScore, symmetryScore, continuityScore),
      correction: correction,
      count: _actionCount,
    );
  }

  /// 计算节奏评分（基于动作频率稳定性）
  int _calculateRhythmScore() {
    if (_lastActionTime == null || _recentScores.isEmpty) {
      return 75;
    }

    final now = DateTime.now();
    final timeDiff = now.difference(_lastActionTime!).inMilliseconds;

    if (timeDiff > 2000) {
      return 60;
    }

    return 100 - ((timeDiff - 1000).abs() / 20).round().clamp(0, 40);
  }

  /// 计算幅度评分（通用）
  int _calculateAmplitudeScore(bool isJumping, SkeletonPoint? leftAnkle, SkeletonPoint? rightAnkle) {
    if (!isJumping || leftAnkle == null || rightAnkle == null) {
      return 70;
    }
    return 85;
  }

  /// 计算深蹲幅度评分
  int _calculateSquatAmplitudeScore(double kneeAngle) {
    if (kneeAngle < 90) return 100;
    if (kneeAngle < 120) return 80;
    if (kneeAngle < 150) return 60;
    return 40;
  }

  /// 计算上肢对称性
  int _calculateSymmetryScore(
    SkeletonPoint? leftWrist, SkeletonPoint? rightWrist,
    SkeletonPoint? leftShoulder, SkeletonPoint? rightShoulder,
  ) {
    if (leftWrist == null || rightWrist == null ||
        leftShoulder == null || rightShoulder == null) {
      return 70;
    }

    final leftArmLength = (leftWrist.position - leftShoulder.position).distance;
    final rightArmLength = (rightWrist.position - rightShoulder.position).distance;

    final diff = (leftArmLength - rightArmLength).abs() / ((leftArmLength + rightArmLength) / 2);
    final score = (1 - diff) * 100;

    return score.round().clamp(0, 100);
  }

  /// 计算腿部对称性
  int _calculateLegSymmetryScore(double leftValue, double rightValue) {
    if (leftValue == 0 && rightValue == 0) return 70;

    final avg = (leftValue + rightValue) / 2;
    if (avg == 0) return 70;

    final diff = (leftValue - rightValue).abs() / avg;
    final score = (1 - diff) * 100;

    return score.round().clamp(0, 100);
  }

  /// 计算手臂对称性
  int _calculateArmSymmetryScore(double leftAngle, double rightAngle) {
    if (leftAngle == 0 && rightAngle == 0) return 70;

    final avg = (leftAngle + rightAngle) / 2;
    if (avg == 0) return 70;

    final diff = (leftAngle - rightAngle).abs() / avg;
    final score = (1 - diff) * 100;

    return score.round().clamp(0, 100);
  }

  /// 计算连贯性评分（基于最近得分的稳定性）
  int _calculateContinuityScore() {
    if (_recentScores.length < 3) {
      return 75;
    }

    final avg = _recentScores.reduce((a, b) => a + b) / _recentScores.length;
    final variance = _recentScores.map((s) => pow(s - avg, 2)).reduce((a, b) => a + b) / _recentScores.length;
    final stability = 100 - (variance / 100).round().clamp(0, 50);

    return stability;
  }

  /// 计算总分
  int _calculateTotalScore(int rhythm, int amplitude, int symmetry, int continuity) {
    // 权重分配
    const rhythmWeight = 0.25;
    const amplitudeWeight = 0.30;
    const symmetryWeight = 0.25;
    const continuityWeight = 0.20;

    final total = rhythm * rhythmWeight +
                  amplitude * amplitudeWeight +
                  symmetry * symmetryWeight +
                  continuity * continuityWeight;

    return total.round().clamp(0, 100);
  }

  /// 生成跳绳纠正建议
  String? _generateJumpRopeCorrection(int rhythmScore, int amplitudeScore, int symmetryScore) {
    if (rhythmScore < 60) return '请保持更稳定的节奏';
    if (amplitudeScore < 60) return '请跳得更高一些';
    if (symmetryScore < 60) return '请保持身体平衡';
    return null;
  }

  /// 生成开合跳纠正建议
  String? _generateJumpingJackCorrection(int amplitudeScore, int symmetryScore) {
    if (amplitudeScore < 70) return '请将手脚张开更大';
    if (symmetryScore < 70) return '请保持左右动作一致';
    return null;
  }

  /// 创建当前评分（当骨骼数据不可用时）
  ActionScore _createCurrentScore(String? correction) {
    return ActionScore(
      rhythmScore: _calculateRhythmScore(),
      amplitudeScore: 70,
      symmetryScore: 70,
      continuityScore: _calculateContinuityScore(),
      totalScore: 75,
      correction: correction,
      count: _actionCount,
    );
  }

  /// 更新历史骨骼数据
  void updatePreviousSkeleton(SkeletonModel skeleton) {
    _previousSkeleton = skeleton;
  }

  /// 重置计数
  void resetCount() {
    _actionCount = 0;
    _recentScores.clear();
  }

  /// 添加评分到历史
  void addScoreToHistory(int score) {
    _recentScores.add(score);
    if (_recentScores.length > 10) {
      _recentScores.removeAt(0);
    }
  }
}

/// 动作分析器工厂
class ActionRecognizerFactory {
  /// 根据动作类型创建分析器
  static ActionRecognizer create(String actionType, SkeletonModel? skeleton) {
    return ActionRecognizer(currentSkeleton: skeleton);
  }
}
