import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_dimensions.dart';
import '../../../core/constants/app_strings.dart';
import '../../../core/router/app_router.dart';
import '../../providers/auth_provider.dart';

/// 注册页面
/// 家长用户注册入口，支持添加儿童用户
class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  final _nicknameController = TextEditingController();
  final _childAgeController = TextEditingController();

  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;
  bool _hasChild = false;
  int _currentStep = 0;

  @override
  void dispose() {
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    _nicknameController.dispose();
    _childAgeController.dispose();
    super.dispose();
  }

  /// 处理注册
  Future<void> _handleRegister() async {
    if (!_formKey.currentState!.validate()) return;

    final success = await ref.read(authStateProvider.notifier).register(
          phone: _phoneController.text.trim(),
          password: _passwordController.text,
          nickname: _nicknameController.text.trim(),
          isParent: true,
          childAge: _hasChild ? int.tryParse(_childAgeController.text) : null,
        );

    if (success && mounted) {
      AppNavigator.goToParentHome(context);
    } else if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(AppStrings.registerFailed),
          backgroundColor: AppColors.error,
        ),
      );
    }
  }

  /// 验证手机号
  String? _validatePhone(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入手机号';
    }
    if (!RegExp(r'^1[3-9]\d{9}$').hasMatch(value)) {
      return AppStrings.phoneFormatError;
    }
    return null;
  }

  /// 验证密码
  String? _validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入密码';
    }
    if (value.length < 6) {
      return AppStrings.passwordFormatError;
    }
    return null;
  }

  /// 验证确认密码
  String? _validateConfirmPassword(String? value) {
    if (value == null || value.isEmpty) {
      return '请确认密码';
    }
    if (value != _passwordController.text) {
      return '两次密码不一致';
    }
    return null;
  }

  /// 验证昵称
  String? _validateNickname(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入昵称';
    }
    if (value.length < 2) {
      return '昵称至少2个字符';
    }
    return null;
  }

  /// 验证儿童年龄
  String? _validateChildAge(String? value) {
    if (!_hasChild) return null;
    if (value == null || value.isEmpty) {
      return '请输入孩子年龄';
    }
    final age = int.tryParse(value);
    if (age == null || age < 3 || age > 12) {
      return '年龄需在3-12岁之间';
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.register),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => AppNavigator.goBack(context),
        ),
      ),
      body: Stepper(
        currentStep: _currentStep,
        onStepContinue: () {
          if (_currentStep < 2) {
            setState(() => _currentStep++);
          } else {
            _handleRegister();
          }
        },
        onStepCancel: () {
          if (_currentStep > 0) {
            setState(() => _currentStep--);
          }
        },
        steps: [
          // 步骤1: 账号信息
          Step(
            title: const Text('账号信息'),
            content: _buildAccountStep(),
            isActive: _currentStep >= 0,
            state: _currentStep > 0 ? StepState.complete : StepState.indexed,
          ),
          // 步骤2: 基本信息
          Step(
            title: const Text('基本信息'),
            content: _buildBasicInfoStep(),
            isActive: _currentStep >= 1,
            state: _currentStep > 1 ? StepState.complete : StepState.indexed,
          ),
          // 步骤3: 儿童信息
          Step(
            title: const Text('儿童信息'),
            content: _buildChildInfoStep(),
            isActive: _currentStep >= 2,
            state: StepState.indexed,
          ),
        ],
      ),
    );
  }

  /// 构建账号信息步骤
  Widget _buildAccountStep() {
    return Form(
      key: _formKey,
      child: Column(
        children: [
          TextFormField(
            controller: _phoneController,
            keyboardType: TextInputType.phone,
            validator: _validatePhone,
            decoration: const InputDecoration(
              labelText: AppStrings.phone,
              prefixIcon: Icon(Icons.phone_outlined),
            ),
          ),
          const SizedBox(height: AppDimensions.spacingMD),
          TextFormField(
            controller: _passwordController,
            obscureText: _obscurePassword,
            validator: _validatePassword,
            decoration: InputDecoration(
              labelText: AppStrings.password,
              prefixIcon: const Icon(Icons.lock_outlined),
              suffixIcon: IconButton(
                icon: Icon(
                  _obscurePassword ? Icons.visibility_off : Icons.visibility,
                ),
                onPressed: () {
                  setState(() => _obscurePassword = !_obscurePassword);
                },
              ),
            ),
          ),
          const SizedBox(height: AppDimensions.spacingMD),
          TextFormField(
            controller: _confirmPasswordController,
            obscureText: _obscureConfirmPassword,
            validator: _validateConfirmPassword,
            decoration: InputDecoration(
              labelText: '确认密码',
              prefixIcon: const Icon(Icons.lock_outlined),
              suffixIcon: IconButton(
                icon: Icon(
                  _obscureConfirmPassword
                      ? Icons.visibility_off
                      : Icons.visibility,
                ),
                onPressed: () {
                  setState(() =>
                      _obscureConfirmPassword = !_obscureConfirmPassword);
                },
              ),
            ),
          ),
        ],
      ),
    );
  }

  /// 构建基本信息步骤
  Widget _buildBasicInfoStep() {
    return Column(
      children: [
        TextFormField(
          controller: _nicknameController,
          validator: _validateNickname,
          decoration: const InputDecoration(
            labelText: '您的昵称',
            prefixIcon: Icon(Icons.person_outlined),
            hintText: '家长昵称',
          ),
        ),
        const SizedBox(height: AppDimensions.spacingMD),
        const Text(
          '请设置您的账号信息',
          style: TextStyle(
            color: AppColors.textSecondary,
            fontSize: 14,
          ),
        ),
      ],
    );
  }

  /// 构建儿童信息步骤
  Widget _buildChildInfoStep() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // 是否添加儿童开关
        SwitchListTile(
          title: const Text('添加儿童账号'),
          subtitle: const Text('为孩子创建专属的运动账号'),
          value: _hasChild,
          onChanged: (value) {
            setState(() => _hasChild = value);
          },
          activeColor: AppColors.primary,
          contentPadding: EdgeInsets.zero,
        ),

        if (_hasChild) ...[
          const SizedBox(height: AppDimensions.spacingMD),
          TextFormField(
            controller: _childAgeController,
            keyboardType: TextInputType.number,
            validator: _validateChildAge,
            decoration: const InputDecoration(
              labelText: '孩子年龄',
              prefixIcon: Icon(Icons.child_care),
              hintText: '3-12岁',
            ),
          ),
          const SizedBox(height: AppDimensions.spacingMD),
          const Text(
            '您之后也可以在设置中添加更多儿童账号',
            style: TextStyle(
              color: AppColors.textSecondary,
              fontSize: 12,
            ),
          ),
        ],
      ],
    );
  }
}
