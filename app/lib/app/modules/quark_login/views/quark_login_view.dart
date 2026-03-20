import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';

import '../../../theme/app_theme.dart';
import '../controllers/quark_login_controller.dart';

class QuarkLoginView extends GetView<QuarkLoginController> {
  const QuarkLoginView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('夸克登录')),
      body: Obx(() {
        if (controller.loading.value) {
          return const Center(child: CircularProgressIndicator());
        }

        return ListView(
          padding: EdgeInsets.fromLTRB(20.w, 20.h, 20.w, 120.h),
          children: [_buildStatusCard()],
        );
      }),
      bottomNavigationBar: SafeArea(
        minimum: EdgeInsets.fromLTRB(20.w, 0, 20.w, 20.h),
        child: Obx(() {
          final label = controller.saving.value
              ? '保存中'
              : controller.configured
              ? '重新登录'
              : '网页登录';
          return FilledButton.icon(
            onPressed: controller.saving.value ? null : controller.openWebLogin,
            style: FilledButton.styleFrom(
              minimumSize: Size(double.infinity, 50.h),
              backgroundColor: AppThemeColors.primary,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16.r),
              ),
            ),
            icon: controller.saving.value
                ? SizedBox(
                    width: 18.w,
                    height: 18.w,
                    child: const CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.public_rounded),
            label: Text(label),
          );
        }),
      ),
    );
  }

  Widget _buildStatusCard() {
    final updatedAt = controller.updatedAt;
    final updatedText = updatedAt == null ? '尚未保存' : _formatDate(updatedAt);
    final statusText = controller.configured ? '已保存' : '未保存';

    return Container(
      padding: EdgeInsets.all(18.w),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A1A),
        borderRadius: BorderRadius.circular(22.r),
        border: Border.all(color: Colors.white.withValues(alpha: 0.06)),
      ),
      child: Row(
        children: [
          Container(
            width: 42.w,
            height: 42.w,
            decoration: BoxDecoration(
              color: const Color(0xFF80CBC4).withValues(alpha: 0.12),
              borderRadius: BorderRadius.circular(14.r),
            ),
            child: const Icon(Icons.cookie_outlined, color: Color(0xFF80CBC4)),
          ),
          SizedBox(width: 12.w),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '夸克 Cookie',
                  style: TextStyle(
                    fontSize: 16.sp,
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                  ),
                ),
                SizedBox(height: 4.h),
                Text(
                  updatedText,
                  style: TextStyle(fontSize: 12.sp, color: Colors.white54),
                ),
                SizedBox(height: 6.h),
                Text(
                  statusText,
                  style: TextStyle(fontSize: 13.sp, color: Colors.white70),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  static String _formatDate(DateTime value) {
    String two(int number) => number.toString().padLeft(2, '0');
    return '${value.year}-${two(value.month)}-${two(value.day)} '
        '${two(value.hour)}:${two(value.minute)}';
  }
}
