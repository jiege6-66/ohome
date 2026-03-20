import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:get/get.dart';

import '../../../data/models/user_model.dart';
import '../../../theme/app_theme.dart';
import '../controllers/plugin_controller.dart';

class PluginView extends GetView<PluginController> {
  const PluginView({super.key});

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.of(context).padding.bottom + 108.h;

    return SafeArea(
      child: Stack(
        children: [
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            height: 220.h,
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    AppThemeColors.primary.withValues(alpha: 0.32),
                    AppThemeColors.secondary.withValues(alpha: 0.18),
                    Colors.transparent,
                  ],
                ),
              ),
            ),
          ),
          ListView(
            padding: EdgeInsets.fromLTRB(20.w, 28.h, 20.w, bottomInset),
            children: [
              Text(
                '设置',
                style: TextStyle(
                  fontSize: 22.sp,
                  fontWeight: FontWeight.w700,
                  color: Colors.white,
                ),
              ),
              SizedBox(height: 20.h),
              _buildProfileCard(),
              SizedBox(height: 24.h),
              Text(
                '基础管理',
                style: TextStyle(
                  fontSize: 14.sp,
                  fontWeight: FontWeight.w600,
                  color: Colors.white70,
                ),
              ),
              SizedBox(height: 12.h),
              Obx(() {
                final isSuperAdmin = controller.isSuperAdmin;
                return Column(
                  children: [
                    if (isSuperAdmin) ...[
                      _SettingsMenuCard(
                        icon: Icons.group_outlined,
                        iconColor: const Color(0xFF80CBC4),
                        title: '用户管理',
                        subtitle: '新增、编辑、删除、重置密码',
                        onTap: controller.openUserManagement,
                      ),
                      SizedBox(height: 12.h),
                      _SettingsMenuCard(
                        icon: Icons.cookie_outlined,
                        iconColor: const Color(0xFFFFB74D),
                        title: '夸克登录',
                        subtitle: '设置 quark_cookies 参数',
                        onTap: controller.openQuarkLogin,
                      ),
                      SizedBox(height: 12.h),
                      _SettingsMenuCard(
                        icon: Icons.travel_explore_rounded,
                        iconColor: const Color(0xFF64B5F6),
                        title: '夸克搜索',
                        subtitle: '配置 HTTP 代理、HTTPS 代理、TG 频道和启用插件',
                        onTap: controller.openQuarkSearchSettings,
                      ),
                      SizedBox(height: 12.h),
                    ],
                    _SettingsMenuCard(
                      icon: Icons.sync_alt_rounded,
                      iconColor: const Color(0xFF81C784),
                      title: '夸克同步',
                      subtitle: '管理夸克自动转存同步任务',
                      onTap: controller.openQuarkSync,
                    ),
                  ],
                );
              }),
              SizedBox(height: 24.h),
              Text(
                '系统',
                style: TextStyle(
                  fontSize: 14.sp,
                  fontWeight: FontWeight.w600,
                  color: Colors.white70,
                ),
              ),
              SizedBox(height: 12.h),
              Obx(
                () => _SettingsMenuCard(
                  icon: Icons.system_update_alt_rounded,
                  iconColor: const Color(0xFF64B5F6),
                  title: '检查更新',
                  subtitle: controller.appUpdateChecking.value
                      ? '正在检查版本信息'
                      : '检查当前是否有可用新版本',
                  onTap: controller.appUpdateChecking.value
                      ? null
                      : controller.checkAppUpdate,
                  trailing: controller.appUpdateChecking.value
                      ? SizedBox(
                          width: 18.w,
                          height: 18.w,
                          child: const CircularProgressIndicator(
                            strokeWidth: 2,
                          ),
                        )
                      : null,
                ),
              ),
              SizedBox(height: 12.h),
              _SettingsMenuCard(
                icon: Icons.logout_rounded,
                iconColor: const Color(0xFFE57373),
                title: '退出登录',
                subtitle: '退出当前账号并返回登录页',
                onTap: controller.logout,
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildProfileCard() {
    return Obx(() {
      final user = controller.user.value;
      final uploading = controller.avatarUploading.value;
      final name = user?.realName.isNotEmpty == true
          ? user!.realName
          : user?.name ?? '未登录';
      final subTitle = user?.roleName.isNotEmpty == true
          ? user!.roleName
          : '移动端设置中心';
      final avatarUrl = user?.avatarUrl ?? '';
      final avatarText = user?.name.isNotEmpty == true
          ? user!.name.substring(0, 1).toUpperCase()
          : 'SZ';

      return Container(
        padding: EdgeInsets.all(18.w),
        decoration: BoxDecoration(
          color: const Color(0xFF1A1A1A),
          borderRadius: BorderRadius.circular(22.r),
          border: Border.all(color: Colors.white.withValues(alpha: 0.06)),
        ),
        child: Row(
          children: [
            GestureDetector(
              onTap: uploading ? null : controller.uploadAvatar,
              child: Stack(
                clipBehavior: Clip.none,
                children: [
                  CircleAvatar(
                    radius: 28.r,
                    backgroundColor: AppThemeColors.primary.withValues(
                      alpha: 0.16,
                    ),
                    backgroundImage: avatarUrl.isNotEmpty
                        ? NetworkImage(avatarUrl)
                        : null,
                    child: avatarUrl.isEmpty
                        ? Text(
                            avatarText,
                            style: TextStyle(
                              fontSize: 16.sp,
                              fontWeight: FontWeight.w700,
                              color: Colors.white,
                            ),
                          )
                        : null,
                  ),
                  Positioned(
                    right: -2.w,
                    bottom: -2.h,
                    child: Container(
                      width: 22.w,
                      height: 22.w,
                      decoration: BoxDecoration(
                        color: const Color(0xFF2A2A2A),
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: Colors.white.withValues(alpha: 0.12),
                        ),
                      ),
                      child: Center(
                        child: uploading
                            ? SizedBox(
                                width: 10.w,
                                height: 10.w,
                                child: const CircularProgressIndicator(
                                  strokeWidth: 1.6,
                                ),
                              )
                            : Icon(
                                Icons.camera_alt_rounded,
                                size: 12.w,
                                color: Colors.white70,
                              ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(width: 14.w),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      fontSize: 18.sp,
                      fontWeight: FontWeight.w600,
                      color: Colors.white,
                    ),
                  ),
                  SizedBox(height: 4.h),
                  Text(
                    subTitle,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(fontSize: 12.sp, color: Colors.white54),
                  ),
                ],
              ),
            ),
            if (user != null) ...[
              SizedBox(width: 12.w),
              _ProfileEditButton(
                loading: controller.profileUpdating.value,
                onTap: controller.profileUpdating.value
                    ? null
                    : () => Get.dialog<void>(
                        _EditProfileDialog(
                          controller: controller,
                          initialUser: user,
                        ),
                      ),
              ),
            ],
          ],
        ),
      );
    });
  }
}

class _ProfileEditButton extends StatelessWidget {
  const _ProfileEditButton({required this.loading, required this.onTap});

  final bool loading;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(999.r),
      child: SizedBox(
        width: 28.w,
        height: 28.w,
        child: Center(
          child: loading
              ? SizedBox(
                  width: 16.w,
                  height: 16.w,
                  child: const CircularProgressIndicator(strokeWidth: 2),
                )
              : Icon(
                  Icons.chevron_right_rounded,
                  size: 24.w,
                  color: Colors.white38,
                ),
        ),
      ),
    );
  }
}

class _EditProfileDialog extends StatefulWidget {
  const _EditProfileDialog({
    required this.controller,
    required this.initialUser,
  });

  final PluginController controller;
  final UserModel initialUser;

  @override
  State<_EditProfileDialog> createState() => _EditProfileDialogState();
}

class _EditProfileDialogState extends State<_EditProfileDialog> {
  final _formKey = GlobalKey<FormState>();

  late final TextEditingController _nameController;
  late final TextEditingController _realNameController;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.initialUser.name);
    _realNameController = TextEditingController(
      text: widget.initialUser.realName,
    );
  }

  @override
  void dispose() {
    _nameController.dispose();
    _realNameController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    FocusScope.of(context).unfocus();
    if (!(_formKey.currentState?.validate() ?? false)) return;

    final changed = await widget.controller.updateProfile(
      name: _nameController.text,
      realName: _realNameController.text,
    );
    if (changed && mounted) {
      Get.back<void>();
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      scrollable: true,
      backgroundColor: const Color(0xFF1A1A1A),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24.r)),
      titlePadding: EdgeInsets.fromLTRB(20.w, 20.h, 20.w, 8.h),
      contentPadding: EdgeInsets.fromLTRB(20.w, 0, 20.w, 8.h),
      actionsPadding: EdgeInsets.fromLTRB(12.w, 0, 12.w, 12.h),
      title: Text(
        '编辑资料',
        style: TextStyle(
          fontSize: 18.sp,
          fontWeight: FontWeight.w700,
          color: Colors.white,
        ),
      ),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '修改当前账号的用户名和昵称',
              style: TextStyle(fontSize: 12.sp, color: Colors.white54),
            ),
            SizedBox(height: 16.h),
            _ProfileTextField(
              label: '用户名',
              controller: _nameController,
              hintText: '请输入用户名',
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return '请输入用户名';
                }
                return null;
              },
            ),
            SizedBox(height: 14.h),
            _ProfileTextField(
              label: '昵称',
              controller: _realNameController,
              hintText: '请输入昵称',
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return '请输入昵称';
                }
                return null;
              },
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Get.back<void>(),
          child: Text(
            '取消',
            style: TextStyle(fontSize: 13.sp, color: Colors.white60),
          ),
        ),
        Obx(
          () => FilledButton(
            onPressed: widget.controller.profileUpdating.value ? null : _submit,
            style: FilledButton.styleFrom(
              backgroundColor: AppThemeColors.primary,
            ),
            child: widget.controller.profileUpdating.value
                ? SizedBox(
                    width: 16.w,
                    height: 16.w,
                    child: const CircularProgressIndicator(strokeWidth: 2),
                  )
                : Text(
                    '保存',
                    style: TextStyle(
                      fontSize: 13.sp,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
          ),
        ),
      ],
    );
  }
}

class _ProfileTextField extends StatelessWidget {
  const _ProfileTextField({
    required this.label,
    required this.controller,
    required this.hintText,
    required this.validator,
  });

  final String label;
  final TextEditingController controller;
  final String hintText;
  final String? Function(String?) validator;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 13.sp,
            fontWeight: FontWeight.w600,
            color: Colors.white70,
          ),
        ),
        SizedBox(height: 8.h),
        TextFormField(
          controller: controller,
          validator: validator,
          style: TextStyle(fontSize: 14.sp, color: Colors.white),
          decoration: InputDecoration(
            hintText: hintText,
            hintStyle: TextStyle(fontSize: 13.sp, color: Colors.white38),
            filled: true,
            fillColor: const Color(0xFF101010),
            contentPadding: EdgeInsets.symmetric(
              horizontal: 14.w,
              vertical: 14.h,
            ),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(14.r),
              borderSide: BorderSide(
                color: Colors.white.withValues(alpha: 0.08),
              ),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(14.r),
              borderSide: BorderSide(
                color: Colors.white.withValues(alpha: 0.08),
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(14.r),
              borderSide: const BorderSide(color: AppThemeColors.primary),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(14.r),
              borderSide: const BorderSide(color: Colors.redAccent),
            ),
            focusedErrorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(14.r),
              borderSide: const BorderSide(color: Colors.redAccent),
            ),
          ),
        ),
      ],
    );
  }
}

class _SettingsMenuCard extends StatelessWidget {
  const _SettingsMenuCard({
    required this.icon,
    required this.iconColor,
    required this.title,
    required this.subtitle,
    required this.onTap,
    this.trailing,
  });

  final IconData icon;
  final Color iconColor;
  final String title;
  final String subtitle;
  final VoidCallback? onTap;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: const Color(0xFF1A1A1A),
      borderRadius: BorderRadius.circular(22.r),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(22.r),
        child: Ink(
          padding: EdgeInsets.symmetric(horizontal: 14.w, vertical: 14.h),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(22.r),
            border: Border.all(color: Colors.white.withValues(alpha: 0.06)),
          ),
          child: Row(
            children: [
              Container(
                width: 40.w,
                height: 40.w,
                decoration: BoxDecoration(
                  color: iconColor.withValues(alpha: 0.12),
                  borderRadius: BorderRadius.circular(14.r),
                ),
                child: Icon(icon, color: iconColor, size: 22.w),
              ),
              SizedBox(width: 14.w),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: TextStyle(
                        fontSize: 15.sp,
                        fontWeight: FontWeight.w600,
                        color: Colors.white,
                      ),
                    ),
                    SizedBox(height: 4.h),
                    Text(
                      subtitle,
                      style: TextStyle(fontSize: 12.sp, color: Colors.white54),
                    ),
                  ],
                ),
              ),
              trailing ??
                  Icon(
                    Icons.chevron_right_rounded,
                    size: 20.w,
                    color: Colors.white38,
                  ),
            ],
          ),
        ),
      ),
    );
  }
}
