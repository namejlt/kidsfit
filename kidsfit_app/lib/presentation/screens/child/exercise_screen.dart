import 'dart:async';
import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_mlkit_pose_detection/google_mlkit_pose_detection.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../domain/entities/exercise_record.dart';
import '../../../services/exercise_controller.dart';
import '../../../services/action_recognizer.dart';
import '../../../services/pose_detection_service.dart';
import '../../providers/exercise_provider.dart';

/// 运动页面
/// AI骨骼识别实时运动训练
class ExerciseScreen extends ConsumerStatefulWidget {
  /// 运动类型
  final String exerciseType;

  const ExerciseScreen({
    super.key,
    required this.exerciseType,
  });

  @override
  ConsumerState<ExerciseScreen> createState() => _ExerciseScreenState();
}

class _ExerciseScreenState extends ConsumerState<ExerciseScreen>
    with TickerProviderStateMixin {
  late ExerciseType _type;
  ExerciseController? _exerciseController;
  CameraController? _cameraController;
  PoseDetector? _poseDetector;
  bool _isInitialized = false;
  bool _isCountingDown = true;
  int _countdown = 3;
  bool _isPaused = false;

  StreamSubscription? _scoreSubscription;
  ActionScore? _currentScore;

  @override
  void initState() {
    super.initState();
    _type = ExerciseType.values.firstWhere(
      (e) => e.value == widget.exerciseType,
      orElse: () => ExerciseType.jumpRope,
    );
    _initializeCamera();
  }

  /// 初始化摄像头和骨骼检测
  Future<void> _initializeCamera() async {
    try {
      final cameras = await availableCameras();
      if (cameras.isEmpty) {
        if (mounted) {
          _showError('没有可用的摄像头');
        }
        return;
      }

      final camera = cameras.firstWhere(
        (c) => c.lensDirection == CameraLensDirection.back,
        orElse: () => cameras.first,
      );

      _cameraController = CameraController(
        camera,
        ResolutionPreset.medium,
        enableAudio: false,
        imageFormatGroup: ImageFormatGroup.yuv420,
      );

      await _cameraController!.initialize();

      final options = PoseDetectorOptions(
        mode: PoseDetectionMode.stream,
        model: PoseDetectionModel.base,
      );
      _poseDetector = PoseDetector(options: options);

      _cameraController!.startImageStream(_processImage);

      if (mounted) {
        setState(() {
          _isInitialized = true;
        });
        _startCountdown();
      }
    } catch (e) {
      _showError('初始化失败: $e');
    }
  }

  /// 处理摄像头图像
  Future<void> _processImage(CameraImage image) async {
    if (_poseDetector == null || _isPaused) return;

    try {
      final inputImage = _convertCameraImage(image);
      if (inputImage == null) return;

      final poses = await _poseDetector!.processImage(inputImage);
      if (poses.isEmpty || !mounted) return;

      final pose = poses.first;

      // 更新运动评分
      if (_exerciseController?.isExercising ?? false) {
        _updateScoreFromPose(pose);
      }
    } catch (e) {
      debugPrint('处理图像失败: $e');
    }
  }

  /// 从骨骼姿态更新评分
  void _updateScoreFromPose(Pose pose) {
    final landmarks = Map.fromEntries(pose.landmarks.entries.map((entry) {
      return MapEntry(
        entry.key,
        SkeletonPoint(
          type: entry.key,
          position: Offset(entry.value.x, entry.value.y),
          confidence: entry.value.likelihood,
        ),
      );
    }));

    final skeleton = SkeletonModel(
      landmarks: landmarks.values.toList(),
      timestamp: DateTime.now(),
      rotation: InputImageRotation.rotation0deg,
    );

    final recognizer = ActionRecognizer(currentSkeleton: skeleton);

    ActionScore score;
    switch (_type) {
      case ExerciseType.jumpRope:
        score = recognizer.analyzeJumpRope();
        break;
      case ExerciseType.jumpingJack:
        score = recognizer.analyzeJumpingJack();
        break;
      case ExerciseType.squat:
        score = recognizer.analyzeSquat();
        break;
      case ExerciseType.sitUp:
        score = recognizer.analyzeSitUp();
        break;
      case ExerciseType.highKnee:
        score = recognizer.analyzeHighKnee();
        break;
      case ExerciseType.pushUp:
        score = recognizer.analyzePushUp();
        break;
    }

    if (mounted) {
      setState(() {
        _currentScore = score;
      });
    }
  }

  /// 转换摄像头图像
  InputImage? _convertCameraImage(CameraImage image) {
    try {
      final camera = _cameraController!.description;
      final rotation = InputImageRotationValue.fromRawValue(camera.sensorOrientation);
      if (rotation == null) return null;

      final format = InputImageFormatValue.fromRawValue(image.format.raw);
      if (format == null) return null;

      final plane = image.planes.first;

      return InputImage.fromBytes(
        bytes: plane.bytes,
        metadata: InputImageMetadata(
          size: Size(image.width.toDouble(), image.height.toDouble()),
          rotation: rotation,
          format: format,
          bytesPerRow: plane.bytesPerRow,
        ),
      );
    } catch (e) {
      debugPrint('图像转换失败: $e');
      return null;
    }
  }

  /// 显示错误信息
  void _showError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  /// 开始倒计时
  void _startCountdown() {
    Future.doWhile(() async {
      if (!mounted) return false;
      await Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          _countdown--;
          if (_countdown <= 0) {
            _isCountingDown = false;
            _startExercise();
          }
        });
      }
      return _countdown > 0;
    });
  }

  /// 开始运动
  void _startExercise() {
    _exerciseController = ExerciseControllerFactory.create(_type);
    _exerciseController?.initialize().then((success) {
      if (success) {
        _exerciseController?.startExercise();
      }
    });

    ref.read(exerciseProvider.notifier).startExercise(_type);
  }

  @override
  Widget build(BuildContext context) {
    final exerciseState = ref.watch(exerciseProvider);

    return Scaffold(
      backgroundColor: AppColors.black,
      body: SafeArea(
        child: Stack(
          children: [
            // 摄像头预览区域
            _buildCameraPreview(),

            // 骨骼检测覆盖层
            if (_isInitialized) _buildSkeletonOverlay(),

            // 顶部信息栏
            _buildTopBar(),

            // 底部控制面板
            _buildBottomPanel(exerciseState),

            // 倒计时覆盖层
            if (_isCountingDown) _buildCountdownOverlay(),
          ],
        ),
      ),
    );
  }

  /// 构建摄像头预览
  Widget _buildCameraPreview() {
    if (!_isInitialized || _cameraController == null) {
      return Container(
        color: Colors.black87,
        child: const Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              CircularProgressIndicator(color: AppColors.primary),
              SizedBox(height: 16),
              Text(
                '正在初始化摄像头...',
                style: TextStyle(color: Colors.white54, fontSize: 16),
              ),
            ],
          ),
        ),
      );
    }

    return SizedBox.expand(
      child: FittedBox(
        fit: BoxFit.cover,
        child: SizedBox(
          width: _cameraController!.value.previewSize!.height,
          height: _cameraController!.value.previewSize!.width,
          child: CameraPreview(_cameraController!),
        ),
      ),
    );
  }

  /// 构建骨骼检测覆盖层
  Widget _buildSkeletonOverlay() {
    return Container(
      color: Colors.transparent,
      child: CustomPaint(
        painter: _SkeletonPainter(
          score: _currentScore,
        ),
        size: Size.infinite,
      ),
    );
  }

  /// 构建顶部信息栏
  Widget _buildTopBar() {
    return Positioned(
      top: 0,
      left: 0,
      right: 0,
      child: Container(
        padding: const EdgeInsets.all(AppDimensions.pagePadding),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Colors.black.withOpacity(0.7),
              Colors.transparent,
            ],
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              onPressed: () => _showExitConfirmation(),
              icon: const Icon(Icons.close, color: Colors.white),
            ),
            Text(
              _type.displayName,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            IconButton(
              onPressed: () {
                setState(() => _isPaused = !_isPaused);
                if (_isPaused) {
                  _exerciseController?.pauseExercise();
                } else {
                  _exerciseController?.resumeExercise();
                }
              },
              icon: Icon(_isPaused ? Icons.play_arrow : Icons.pause, color: Colors.white),
            ),
          ],
        ),
      ),
    );
  }

  /// 显示退出确认对话框
  void _showExitConfirmation() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认退出'),
        content: const Text('确定要退出运动吗？本次运动数据将不会被保存。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              _exerciseController?.stopExercise();
              Navigator.pop(context);
            },
            child: const Text('退出', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  /// 构建底部控制面板
  Widget _buildBottomPanel(ExerciseState state) {
    final score = _currentScore ?? ActionScore(
      rhythmScore: 0,
      amplitudeScore: 0,
      symmetryScore: 0,
      continuityScore: 0,
      totalScore: 0,
      count: 0,
    );

    return Positioned(
      bottom: 0,
      left: 0,
      right: 0,
      child: Container(
        padding: const EdgeInsets.all(AppDimensions.pagePadding),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.bottomCenter,
            end: Alignment.topCenter,
            colors: [
              Colors.black.withOpacity(0.9),
              Colors.transparent,
            ],
          ),
        ),
        child: Column(
          children: [
            _buildScoreDisplay(score),
            const SizedBox(height: AppDimensions.spacingLG),
            _buildCountDisplay(score),
            const SizedBox(height: AppDimensions.spacingLG),
            if (score.correction != null) _buildCorrectionTip(score.correction!),
            const SizedBox(height: AppDimensions.spacingMD),
            _buildControlButtons(),
          ],
        ),
      ),
    );
  }

  /// 构建评分显示
  Widget _buildScoreDisplay(ActionScore score) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildScoreItem('节奏', score.rhythmScore, AppColors.primary),
        const SizedBox(width: 24),
        _buildScoreItem('幅度', score.amplitudeScore, AppColors.secondary),
        const SizedBox(width: 24),
        _buildScoreItem('综合', score.totalScore, AppColors.accent),
      ],
    );
  }

  /// 构建评分项
  Widget _buildScoreItem(String label, int score, Color color) {
    return Column(
      children: [
        Text(
          label,
          style: const TextStyle(color: Colors.white70, fontSize: 12),
        ),
        const SizedBox(height: 4),
        Container(
          width: 60,
          height: 60,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            border: Border.all(color: color, width: 3),
          ),
          child: Center(
            child: Text(
              '$score',
              style: TextStyle(
                color: color,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
      ],
    );
  }

  /// 构建计数显示
  Widget _buildCountDisplay(ActionScore score) {
    return Column(
      children: [
        const Text(
          '完成次数',
          style: TextStyle(color: Colors.white70, fontSize: 14),
        ),
        const SizedBox(height: 8),
        Text(
          '${score.count}',
          style: const TextStyle(
            color: Colors.white,
            fontSize: 72,
            fontWeight: FontWeight.bold,
          ),
        ),
        const Text(
          '目标: 50次',
          style: TextStyle(color: Colors.white54, fontSize: 16),
        ),
      ],
    );
  }

  /// 构建纠正提示
  Widget _buildCorrectionTip(String tip) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.2),
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: AppColors.primary),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.record_voice_over, color: AppColors.primary, size: 20),
          const SizedBox(width: 8),
          Text(tip, style: const TextStyle(color: Colors.white, fontSize: 14)),
        ],
      ),
    );
  }

  /// 构建控制按钮
  Widget _buildControlButtons() {
    return Row(
      children: [
        Expanded(
          child: ElevatedButton.icon(
            onPressed: () => _showExitConfirmation(),
            icon: const Icon(Icons.stop),
            label: const Text('结束'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red.withOpacity(0.8),
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
              ),
            ),
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: ElevatedButton.icon(
            onPressed: () {
              setState(() => _isPaused = !_isPaused);
              if (_isPaused) {
                _exerciseController?.pauseExercise();
              } else {
                _exerciseController?.resumeExercise();
              }
            },
            icon: Icon(_isPaused ? Icons.play_arrow : Icons.pause),
            label: Text(_isPaused ? '继续' : '暂停'),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(AppDimensions.radiusMD),
              ),
            ),
          ),
        ),
      ],
    );
  }

  /// 构建倒计时覆盖层
  Widget _buildCountdownOverlay() {
    return Container(
      color: Colors.black.withOpacity(0.8),
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              '准备开始',
              style: TextStyle(color: Colors.white70, fontSize: 24),
            ),
            const SizedBox(height: 24),
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(color: AppColors.primary, width: 4),
              ),
              child: Center(
                child: Text(
                  '$_countdown',
                  style: const TextStyle(
                    color: AppColors.primary,
                    fontSize: 72,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _scoreSubscription?.cancel();
    _exerciseController?.dispose();
    _cameraController?.dispose();
    _poseDetector?.close();
    super.dispose();
  }
}

/// 骨骼绘制器
class _SkeletonPainter extends CustomPainter {
  final ActionScore? score;

  _SkeletonPainter({this.score});

  @override
  void paint(Canvas canvas, Size size) {
    if (score == null) return;

    final paint = Paint()
      ..color = AppColors.primary
      ..strokeWidth = 3
      ..style = PaintingStyle.stroke;

    final pointPaint = Paint()
      ..color = AppColors.accent
      ..style = PaintingStyle.fill;

    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width * 0.15;

    canvas.drawCircle(center, radius, paint);
    canvas.drawCircle(center, 8, pointPaint);

    final armStart = Offset(center.dx - radius, center.dy);
    final armEnd = Offset(center.dx + radius, center.dy);
    canvas.drawLine(armStart, armEnd, paint);

    final leftLegStart = Offset(center.dx, center.dy);
    final leftLegEnd = Offset(center.dx - radius * 0.8, center.dy + radius * 1.2);
    canvas.drawLine(leftLegStart, leftLegEnd, paint);

    final rightLegStart = Offset(center.dx, center.dy);
    final rightLegEnd = Offset(center.dx + radius * 0.8, center.dy + radius * 1.2);
    canvas.drawLine(rightLegStart, rightLegEnd, paint);
  }

  @override
  bool shouldRepaint(covariant _SkeletonPainter oldDelegate) {
    return oldDelegate.score != score;
  }
}
