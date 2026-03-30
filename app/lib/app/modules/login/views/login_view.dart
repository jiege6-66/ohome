import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';

import '../../../theme/app_theme.dart';
import '../../../utils/app_env.dart';
import '../../../widgets/backend_address_badge.dart';
import '../../../widgets/server_settings_sheet.dart';
import '../controllers/login_controller.dart';

class LoginView extends GetView<LoginController> {
  const LoginView({super.key});

  @override
  Widget build(BuildContext context) {
    final viewInsets = MediaQuery.viewInsetsOf(context);
    final isDesktop = _isDesktopLayout(context);

    return Scaffold(
      body: Stack(
        children: [
          Container(color: const Color(0xFF0F172A)),
          Positioned(
            top: -_rh(context, 140, 120),
            right: -_rw(context, 50, 40),
            child: Container(
              width: _rw(context, 300, 260),
              height: _rw(context, 300, 260),
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: const Color(0xFF3B82F6).withValues(alpha: 0.4),
              ),
            ),
          ),
          Positioned(
            bottom: -_rh(context, 50, 40),
            left: -_rw(context, 100, 80),
            child: Container(
              width: _rw(context, 250, 220),
              height: _rw(context, 250, 220),
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: const Color(0xFF8B5CF6).withValues(alpha: 0.4),
              ),
            ),
          ),
          BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 80, sigmaY: 80),
            child: Container(color: Colors.transparent),
          ),
          SafeArea(
            child: SingleChildScrollView(
              padding: EdgeInsets.fromLTRB(
                _rw(context, 20, 32),
                _rh(context, 16, 24),
                _rw(context, 20, 32),
                viewInsets.bottom + _rh(context, 24, 32),
              ),
              child: Center(
                child: ConstrainedBox(
                  constraints: BoxConstraints(maxWidth: isDesktop ? 1040 : 420),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Align(
                        alignment: Alignment.centerRight,
                        child: ValueListenableBuilder<TextEditingValue>(
                          valueListenable: controller.apiBaseUrlController,
                          builder: (context, value, _) {
                            final address = value.text.trim().isEmpty
                                ? AppEnv.instance.apiBaseUrlInputValue
                                : value.text.trim();
                            return Wrap(
                              alignment: WrapAlignment.end,
                              crossAxisAlignment: WrapCrossAlignment.center,
                              spacing: _rw(context, 12, 12),
                              runSpacing: _rh(context, 12, 10),
                              children: [
                                BackendAddressBadge(address: address),
                                const _SettingsIconButton(),
                              ],
                            );
                          },
                        ),
                      ),
                      SizedBox(height: _rh(context, 64, 40)),
                      Center(
                        child: ClipRRect(
                          borderRadius: BorderRadius.circular(
                            _rr(context, 32, 24),
                          ),
                          child: BackdropFilter(
                            filter: ImageFilter.blur(sigmaX: 16, sigmaY: 16),
                            child: Container(
                              width: isDesktop ? 420 : _rw(context, 340, 420),
                              padding: EdgeInsets.fromLTRB(
                                _rw(context, 32, 32),
                                _rh(context, 40, 32),
                                _rw(context, 32, 32),
                                _rh(context, 32, 28),
                              ),
                              decoration: BoxDecoration(
                                color: Colors.white.withValues(alpha: 0.04),
                                borderRadius: BorderRadius.circular(
                                  _rr(context, 32, 24),
                                ),
                                border: Border.all(
                                  color: Colors.white.withValues(alpha: 0.1),
                                  width: 1,
                                ),
                              ),
                              child: Obx(
                                () => Form(
                                  key: controller.loginFormKey,
                                  autovalidateMode:
                                      controller.autoValidateMode.value,
                                  child: Column(
                                    children: [
                                      const _Header(),
                                      SizedBox(height: _rh(context, 28, 24)),
                                      const _CredentialsFields(),
                                      SizedBox(height: _rh(context, 28, 24)),
                                      const _SubmitButton(),
                                      SizedBox(height: _rh(context, 14, 12)),
                                      const _RegisterEntryButton(),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _Header extends StatelessWidget {
  const _Header();

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '欢迎回来',
                style: TextStyle(
                  fontSize: _rs(context, 15, 14),
                  fontWeight: FontWeight.w400,
                  color: Colors.white70,
                  letterSpacing: 1.1,
                ),
              ),
              SizedBox(height: _rh(context, 4, 4)),
              Text(
                '登录',
                style: TextStyle(
                  fontSize: _rs(context, 34, 28),
                  fontWeight: FontWeight.w800,
                  color: Colors.white,
                  letterSpacing: 1.2,
                ),
              ),
            ],
          ),
        ),
        Container(
          padding: EdgeInsets.all(_rw(context, 4, 4)),
          decoration: BoxDecoration(
            color: Colors.white.withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(_rr(context, 16, 14)),
            border: Border.all(color: Colors.white.withValues(alpha: 0.15)),
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(_rr(context, 12, 10)),
            child: Image.asset(
              'assets/images/logo.png',
              width: _rw(context, 56, 48),
              height: _rw(context, 56, 48),
              fit: BoxFit.cover,
            ),
          ),
        ),
      ],
    );
  }
}

class _SettingsIconButton extends GetView<LoginController> {
  const _SettingsIconButton();

  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final hasFoundServer = controller.hasFoundServer;

      return GestureDetector(
        onTap: () => _openServerSettingsSheet(context),
        child: Stack(
          clipBehavior: Clip.none,
          children: [
            Container(
              width: _rw(context, 44, 40),
              height: _rw(context, 44, 40),
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.08),
                borderRadius: BorderRadius.circular(_rr(context, 14, 12)),
                border: Border.all(color: Colors.white12),
              ),
              child: Icon(
                Icons.settings_rounded,
                size: _rs(context, 22, 20),
                color: Colors.white,
              ),
            ),
            if (hasFoundServer)
              Positioned(
                top: -_rh(context, 4, 4),
                right: -_rw(context, 4, 4),
                child: Container(
                  width: _rw(context, 18, 16),
                  height: _rw(context, 18, 16),
                  decoration: BoxDecoration(
                    color: const Color(0xFF21C47B),
                    borderRadius: BorderRadius.circular(999),
                    border: Border.all(
                      color: const Color(0xFF1E1E1E),
                      width: 2,
                    ),
                  ),
                  child: Icon(
                    Icons.check_rounded,
                    size: _rs(context, 10, 9),
                    color: Colors.white,
                  ),
                ),
              ),
          ],
        ),
      );
    });
  }

  Future<void> _openServerSettingsSheet(BuildContext context) {
    return Get.bottomSheet<void>(
      Obx(
        () => ServerSettingsSheet(
          apiBaseUrlController: controller.apiBaseUrlController,
          servers: controller.discoveredServers,
          errorText: controller.discoveryErrorMessage.value,
          isDiscovering: controller.isDiscovering.value,
          isManualEntryMode: controller.isManualEntryMode.value,
          selectedServer: controller.selectedServer.value,
          onToggleManualEntryMode: controller.toggleManualEntryMode,
          onRefreshDiscovery: controller.refreshDiscovery,
          onSelectServer: (server) {
            controller.selectDiscoveredServer(server);
            Get.back<void>();
          },
        ),
      ),
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
    );
  }
}

class _CredentialsFields extends GetView<LoginController> {
  const _CredentialsFields();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        TextFormField(
          controller: controller.nameController,
          style: TextStyle(color: Colors.white, fontSize: _rs(context, 16, 15)),
          decoration: _fieldDecoration(
            context,
            hintText: '输入用户名',
            icon: Icons.person_outline_rounded,
          ),
          validator: controller.validateName,
        ),
        SizedBox(height: _rh(context, 18, 16)),
        TextFormField(
          controller: controller.passwordController,
          obscureText: true,
          style: TextStyle(color: Colors.white, fontSize: _rs(context, 16, 15)),
          decoration: _fieldDecoration(
            context,
            hintText: '输入密码',
            icon: Icons.lock_outline_rounded,
          ),
          validator: controller.validatePassword,
        ),
      ],
    );
  }

  InputDecoration _fieldDecoration(
    BuildContext context, {
    required String hintText,
    required IconData icon,
  }) {
    return InputDecoration(
      hintText: hintText,
      hintStyle: TextStyle(color: Colors.white54, fontSize: _rs(context, 15, 14)),
      prefixIcon: Icon(icon, color: Colors.white54, size: _rs(context, 22, 20)),
      filled: true,
      fillColor: Colors.white.withValues(alpha: 0.05),
      contentPadding: EdgeInsets.symmetric(
        vertical: _rh(context, 18, 16),
        horizontal: _rw(context, 20, 18),
      ),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(_rr(context, 16, 14)),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(_rr(context, 16, 14)),
        borderSide: const BorderSide(color: AppThemeColors.primary, width: 1.5),
      ),
      errorStyle: TextStyle(
        color: Colors.redAccent.shade100,
        fontSize: _rs(context, 12, 12),
      ),
    );
  }
}

class _SubmitButton extends GetView<LoginController> {
  const _SubmitButton();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: _rh(context, 52, 48),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(_rr(context, 16, 14)),
        gradient: const LinearGradient(
          colors: [Color(0xFF3B82F6), Color(0xFF2563EB)],
          begin: Alignment.centerLeft,
          end: Alignment.centerRight,
        ),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF3B82F6).withValues(alpha: 0.3),
            blurRadius: _rr(context, 12, 10),
            offset: Offset(0, _rh(context, 4, 4)),
          ),
        ],
      ),
      child: ElevatedButton(
        onPressed: controller.login,
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_rr(context, 16, 14)),
          ),
        ),
        child: Obx(
          () => controller.isLoading.value
              ? SizedBox(
                  width: _rs(context, 24, 20),
                  height: _rs(context, 24, 20),
                  child: const CircularProgressIndicator(
                    strokeWidth: 2.5,
                    valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                  ),
                )
              : Text(
                  '立即登录',
                  style: TextStyle(
                    fontSize: _rs(context, 16, 15),
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                    letterSpacing: 1.1,
                  ),
                ),
        ),
      ),
    );
  }
}

class _RegisterEntryButton extends GetView<LoginController> {
  const _RegisterEntryButton();

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: TextButton(
        onPressed: controller.openRegister,
        style: TextButton.styleFrom(
          alignment: Alignment.centerRight,
          padding: EdgeInsets.zero,
          minimumSize: Size.zero,
          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        ),
        child: Text.rich(
          TextSpan(
            style: TextStyle(
              fontSize: _rs(context, 13, 13),
              color: Colors.white70,
              fontWeight: FontWeight.w500,
            ),
            children: const [
              TextSpan(text: '没有账号？ '),
              TextSpan(
                text: '立即注册',
                style: TextStyle(
                  color: AppThemeColors.primary,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
          textAlign: TextAlign.right,
        ),
      ),
    );
  }
}

bool _isDesktopLayout(BuildContext context) =>
    MediaQuery.sizeOf(context).width >= 900;

double _rw(BuildContext context, double mobile, double desktop) =>
    _isDesktopLayout(context) ? desktop : mobile.w;

double _rh(BuildContext context, double mobile, double desktop) =>
    _isDesktopLayout(context) ? desktop : mobile.h;

double _rs(BuildContext context, double mobile, double desktop) =>
    _isDesktopLayout(context) ? desktop : mobile.sp;

double _rr(BuildContext context, double mobile, double desktop) =>
    _isDesktopLayout(context) ? desktop : mobile.r;
