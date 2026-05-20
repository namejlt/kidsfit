import 'dart:async';
import 'dart:math';
import 'package:camera/camera.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:google_mlkit_pose_detection/google_mlkit_pose_detection.dart';

/// 骨骼关键点
class SkeletonPoint {
  final PoseLandmarkType type;
  final Offset position;
  final double confidence;

  SkeletonPoint({
    required this.type,
    required this.position,
    required this.confidence,
  });
}

/// 骨骼模型
class SkeletonModel {
  final List<SkeletonPoint> landmarks;
  final DateTime timestamp;
  final InputImageRotation rotation;

  SkeletonModel({
    required this.landmarks,
    required this.timestamp,
    required this.rotation,
  });

  /// 获取特定关键点
  SkeletonPoint? getPoint(PoseLandmarkType type) {
    try {
      return landmarks.firstWhere((p) => p.type == type);
    } catch (_) {
      return null;
    }
  }

  /// 计算两点之间的角度（弧度转度）
  double calculateAngle(
    PoseLandmarkType point1Type,
    PoseLandmarkType vertexType,
    PoseLandmarkType point2Type,
  ) {
    final p1 = getPoint(point1Type);
    final p2 = getPoint(vertexType);
    final p3 = getPoint(point2Type);

    if (p1 == null || p2 == null || p3 == null) return 0;

    final angle = atan2(p3.position.dy - p2.position.dy, p3.position.dx - p2.position.dx) -
        atan2(p1.position.dy - p2.position.dy, p1.position.dx - p2.position.dx);

    var degrees = angle * 180 / pi;
    if (degrees < 0) degrees += 360;
    if (degrees > 180) degrees = 360 - degrees;

    return degrees;
  }
}

/// 骨骼检测服务
/// 负责初始化摄像头、处理视频流、进行骨骼关键点检测
class PoseDetectionService {
  CameraController? _cameraController;
  PoseDetector? _poseDetector;
  bool _isInitialized = false;
  bool _isProcessing = false;

  final _skeletonController = StreamController<SkeletonModel>.broadcast();
  final _errorController = StreamController<String>.broadcast();

  /// 骨骼数据流
  Stream<SkeletonModel> get skeletonStream => _skeletonController.stream;

  /// 错误流
  Stream<String> get errorStream => _errorController.stream;

  /// 是否已初始化
  bool get isInitialized => _isInitialized;

  /// 初始化服务
  Future<bool> initialize() async {
    if (_isInitialized) return true;

    try {
      // 获取可用摄像头列表
      final cameras = await availableCameras();
      if (cameras.isEmpty) {
        _errorController.add('没有可用的摄像头');
        return false;
      }

      // 选择后置摄像头
      final camera = cameras.firstWhere(
        (c) => c.lensDirection == CameraLensDirection.back,
        orElse: () => cameras.first,
      );

      // 初始化摄像头控制器
      _cameraController = CameraController(
        camera,
        ResolutionPreset.medium,
        enableAudio: false,
        imageFormatGroup: ImageFormatGroup.yuv420,
      );

      await _cameraController!.initialize();

      // 初始化骨骼检测器
      final options = PoseDetectorOptions(
        mode: PoseDetectionMode.stream,
        model: PoseDetectionModel.base,
      );
      _poseDetector = PoseDetector(options: options);

      _isInitialized = true;
      return true;
    } catch (e) {
      _errorController.add('初始化失败: $e');
      return false;
    }
  }

  /// 启动骨骼检测
  Future<void> startDetection() async {
    if (!_isInitialized || _isProcessing) return;

    _isProcessing = true;
    _cameraController!.startImageStream(_processImage);
  }

  /// 停止骨骼检测
  Future<void> stopDetection() async {
    if (!_isInitialized || !_isProcessing) return;

    _isProcessing = false;
    await _cameraController?.stopImageStream();
  }

  /// 处理摄像头图像
  Future<void> _processImage(CameraImage image) async {
    if (_poseDetector == null) return;

    try {
      // 转换图像格式
      final inputImage = _convertCameraImage(image);
      if (inputImage == null) return;

      // 进行骨骼检测
      final poses = await _poseDetector!.processImage(inputImage);

      if (poses.isNotEmpty) {
        final pose = poses.first;
        final skeleton = _convertToSkeleton(pose, InputImageRotation.rotation0deg);
        _skeletonController.add(skeleton);
      }
    } catch (e) {
      _errorController.add('检测失败: $e');
    }
  }

  /// 转换摄像头图像为ML Kit输入图像
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

  /// 转换检测结果为骨骼模型
  SkeletonModel _convertToSkeleton(Pose pose, InputImageRotation rotation) {
    final landmarks = pose.landmarks.entries.map((entry) {
      final landmark = entry.value;
      return SkeletonPoint(
        type: entry.key,
        position: Offset(landmark.x, landmark.y),
        confidence: landmark.likelihood,
      );
    }).toList();

    return SkeletonModel(
      landmarks: landmarks,
      timestamp: DateTime.now(),
      rotation: rotation,
    );
  }

  /// 释放资源
  Future<void> dispose() async {
    await stopDetection();
    await _cameraController?.dispose();
    _poseDetector?.close();
    _skeletonController.close();
    _errorController.close();
    _isInitialized = false;
  }
}
