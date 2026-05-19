import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../presentation/screens/splash/splash_screen.dart';
import '../../presentation/screens/auth/login_screen.dart';
import '../../presentation/screens/auth/register_screen.dart';
import '../../presentation/screens/child/child_home_screen.dart';
import '../../presentation/screens/child/exercise_screen.dart';
import '../../presentation/screens/child/challenge_screen.dart';
import '../../presentation/screens/child/achievement_screen.dart';
import '../../presentation/screens/parent/parent_home_screen.dart';
import '../../presentation/screens/parent/vision_screen.dart';
import '../../presentation/screens/parent/training_plan_screen.dart';
import '../../presentation/providers/auth_provider.dart';

/// KidsFit应用路由配置
/// 使用go_router进行路由管理
class AppRouter {
  AppRouter._();

  /// 路由名称常量
  static const String splash = '/';
  static const String login = '/login';
  static const String register = '/register';
  static const String childHome = '/child/home';
  static const String parentHome = '/parent/home';
  static const String exercise = '/child/exercise';
  static const String challenge = '/child/challenge';
  static const String achievement = '/child/achievement';
  static const String vision = '/parent/vision';
  static const String trainingPlan = '/parent/training-plan';

  /// 路由配置
  static final GoRouter router = GoRouter(
    initialLocation: splash,
    debugLogDiagnostics: true,
    routes: [
      // 启动页
      GoRoute(
        path: splash,
        name: 'splash',
        builder: (context, state) => const SplashScreen(),
      ),

      // 登录页
      GoRoute(
        path: login,
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),

      // 注册页
      GoRoute(
        path: register,
        name: 'register',
        builder: (context, state) => const RegisterScreen(),
      ),

      // 儿童端首页
      GoRoute(
        path: childHome,
        name: 'childHome',
        builder: (context, state) => const ChildHomeScreen(),
      ),

      // 儿童端运动页
      GoRoute(
        path: exercise,
        name: 'exercise',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>?;
          return ExerciseScreen(
            exerciseType: extra?['exerciseType'] ?? 'jump_rope',
          );
        },
      ),

      // 儿童端挑战页
      GoRoute(
        path: challenge,
        name: 'challenge',
        builder: (context, state) => const ChallengeScreen(),
      ),

      // 儿童端成就页
      GoRoute(
        path: achievement,
        name: 'achievement',
        builder: (context, state) => const AchievementScreen(),
      ),

      // 家长端首页
      GoRoute(
        path: parentHome,
        name: 'parentHome',
        builder: (context, state) => const ParentHomeScreen(),
      ),

      // 家长端视力页
      GoRoute(
        path: vision,
        name: 'vision',
        builder: (context, state) => const VisionScreen(),
      ),

      // 家长端训练计划页
      GoRoute(
        path: trainingPlan,
        name: 'trainingPlan',
        builder: (context, state) => const TrainingPlanScreen(),
      ),
    ],

    // 错误处理
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Text('路由错误: ${state.error}'),
      ),
    ),

    // 路由重定向
    redirect: (context, state) {
      // 在这里可以实现路由守卫逻辑
      // 例如：未登录时重定向到登录页
      return null;
    },
  );
}

/// 路由跳转助手类
class AppNavigator {
  AppNavigator._();

  /// 跳转到登录页
  static void goToLogin(BuildContext context) {
    context.pushReplacement(AppRouter.login);
  }

  /// 跳转到注册页
  static void goToRegister(BuildContext context) {
    context.push(AppRouter.register);
  }

  /// 跳转到儿童端首页
  static void goToChildHome(BuildContext context) {
    context.pushReplacement(AppRouter.childHome);
  }

  /// 跳转到家长端首页
  static void goToParentHome(BuildContext context) {
    context.pushReplacement(AppRouter.parentHome);
  }

  /// 跳转到运动页面
  static void goToExercise(BuildContext context, String exerciseType) {
    context.push(
      AppRouter.exercise,
      extra: {'exerciseType': exerciseType},
    );
  }

  /// 跳转到挑战页面
  static void goToChallenge(BuildContext context) {
    context.push(AppRouter.challenge);
  }

  /// 跳转到成就页面
  static void goToAchievement(BuildContext context) {
    context.push(AppRouter.achievement);
  }

  /// 跳转到视力页面
  static void goToVision(BuildContext context) {
    context.push(AppRouter.vision);
  }

  /// 跳转到训练计划页面
  static void goToTrainingPlan(BuildContext context) {
    context.push(AppRouter.trainingPlan);
  }

  /// 返回上一页
  static void goBack(BuildContext context) {
    if (context.canPop()) {
      context.pop();
    }
  }
}
