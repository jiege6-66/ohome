import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';

import '../../../theme/app_theme.dart';
import '../controllers/messages_controller.dart';

class MessagesView extends GetView<MessagesController> {
  const MessagesView({super.key});

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.of(context).padding.bottom + 24.h;
    return Scaffold(
      backgroundColor: AppThemeColors.pageBackground,
      appBar: AppBar(
        title: const Text('消息'),
        actions: [
          TextButton(
            onPressed: controller.markAllRead,
            child: const Text('全部已读'),
          ),
          SizedBox(width: 8.w),
        ],
      ),
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: EdgeInsets.fromLTRB(20.w, 16.h, 20.w, 12.h),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        '全局通知与提醒消息',
                        style: TextStyle(
                          fontSize: 13.sp,
                          color: Colors.white60,
                        ),
                      ),
                      GestureDetector(
                        onTap: () => controller.toggleUnreadOnly(
                          !controller.unreadOnly.value,
                        ),
                        child: Obx(
                          () => Container(
                            padding: EdgeInsets.symmetric(
                              horizontal: 10.w,
                              vertical: 4.h,
                            ),
                            decoration: BoxDecoration(
                              color: controller.unreadOnly.value
                                  ? AppThemeColors.primary.withValues(
                                      alpha: 0.15,
                                    )
                                  : Colors.white.withValues(alpha: 0.05),
                              borderRadius: BorderRadius.circular(999.r),
                              border: Border.all(
                                color: controller.unreadOnly.value
                                    ? AppThemeColors.primary.withValues(
                                        alpha: 0.5,
                                      )
                                    : Colors.white.withValues(alpha: 0.1),
                              ),
                            ),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(
                                  controller.unreadOnly.value
                                      ? Icons.check_circle
                                      : Icons.radio_button_unchecked,
                                  size: 14.sp,
                                  color: controller.unreadOnly.value
                                      ? AppThemeColors.primary
                                      : Colors.white54,
                                ),
                                SizedBox(width: 4.w),
                                Text(
                                  '未读 (${controller.unreadCount.value})',
                                  style: TextStyle(
                                    fontSize: 12.sp,
                                    color: controller.unreadOnly.value
                                        ? AppThemeColors.primary
                                        : Colors.white54,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 12.h),
                  _buildSourceFilter(),
                ],
              ),
            ),
            Expanded(child: _buildList(bottomInset)),
          ],
        ),
      ),
    );
  }

  Widget _buildList(double bottomInset) {
    return Obx(() {
      if (controller.loading.value && controller.messages.isEmpty) {
        return const Center(child: CircularProgressIndicator());
      }

      if (controller.messages.isEmpty) {
        return RefreshIndicator(
          onRefresh: () => controller.loadMessages(refresh: true),
          child: ListView(
            padding: EdgeInsets.only(bottom: bottomInset),
            children: [
              const SizedBox(height: 120),
              Center(
                child: Text(
                  controller.emptyStateText,
                  style: const TextStyle(color: Colors.white54),
                ),
              ),
            ],
          ),
        );
      }

      return RefreshIndicator(
        onRefresh: () => controller.loadMessages(refresh: true),
        child: ListView.separated(
          controller: controller.scrollController,
          padding: EdgeInsets.fromLTRB(20.w, 0, 20.w, bottomInset),
          itemCount:
              controller.messages.length + (controller.hasMore.value ? 1 : 0),
          separatorBuilder: (_, _) => SizedBox(height: 12.h),
          itemBuilder: (_, index) {
            if (index >= controller.messages.length) {
              return const Padding(
                padding: EdgeInsets.symmetric(vertical: 12),
                child: Center(child: CircularProgressIndicator()),
              );
            }
            final message = controller.messages[index];
            final sourceColor = controller.sourceColorOf(message);
            return Material(
              color: const Color(0xFF1A1A1A),
              borderRadius: BorderRadius.circular(20.r),
              child: InkWell(
                borderRadius: BorderRadius.circular(20.r),
                onTap: () => controller.openMessage(message),
                child: Ink(
                  padding: EdgeInsets.all(16.w),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20.r),
                    border: Border.all(
                      color: Colors.white.withValues(alpha: 0.06),
                    ),
                  ),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                        width: 10.w,
                        height: 10.w,
                        margin: EdgeInsets.only(top: 6.h),
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: message.read
                              ? Colors.white24
                              : const Color(0xFFE53935),
                        ),
                      ),
                      SizedBox(width: 12.w),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Container(
                                  padding: EdgeInsets.symmetric(
                                    horizontal: 8.w,
                                    vertical: 4.h,
                                  ),
                                  decoration: BoxDecoration(
                                    color: sourceColor.withValues(alpha: 0.14),
                                    borderRadius: BorderRadius.circular(999.r),
                                  ),
                                  child: Text(
                                    controller.sourceLabelOf(message),
                                    style: TextStyle(
                                      fontSize: 10.sp,
                                      fontWeight: FontWeight.w600,
                                      color: sourceColor,
                                    ),
                                  ),
                                ),
                                const Spacer(),
                                Obx(() {
                                  final deleting = controller.isDeleting(
                                    message.id,
                                  );
                                  return IconButton(
                                    onPressed: deleting
                                        ? null
                                        : () =>
                                              controller.deleteMessage(message),
                                    tooltip: '删除消息',
                                    splashRadius: 18.r,
                                    visualDensity: VisualDensity.compact,
                                    icon: deleting
                                        ? SizedBox(
                                            width: 16.w,
                                            height: 16.w,
                                            child:
                                                const CircularProgressIndicator(
                                                  strokeWidth: 2,
                                                ),
                                          )
                                        : Icon(
                                            Icons.delete_outline_rounded,
                                            size: 18.sp,
                                            color: Colors.white54,
                                          ),
                                  );
                                }),
                              ],
                            ),
                            SizedBox(height: 10.h),
                            Text(
                              message.title,
                              style: TextStyle(
                                fontSize: 14.sp,
                                fontWeight: FontWeight.w600,
                                color: Colors.white,
                              ),
                            ),
                            SizedBox(height: 6.h),
                            Text(
                              controller.summaryOf(message),
                              style: TextStyle(
                                fontSize: 12.sp,
                                color: Colors.white60,
                              ),
                            ),
                            SizedBox(height: 8.h),
                            Text(
                              message.triggerDate == null
                                  ? '提醒时间未知'
                                  : '提醒日期 ${_formatDate(message.triggerDate!)}',
                              style: TextStyle(
                                fontSize: 11.sp,
                                color: Colors.white38,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        ),
      );
    });
  }

  Widget _buildSourceFilter() {
    const filters = <String>['', 'drops', 'quark', 'system'];
    return Obx(
      () => Wrap(
        spacing: 8.w,
        runSpacing: 8.h,
        children: filters
            .map(
              (source) => _SourceFilterChip(
                label: controller.sourceFilterLabel(source),
                selected: controller.sourceFilter.value == source,
                color: source.isEmpty
                    ? AppThemeColors.primary
                    : controller.sourceColor(source),
                onTap: () => controller.changeSourceFilter(source),
              ),
            )
            .toList(growable: false),
      ),
    );
  }

  static String _formatDate(DateTime value) {
    String two(int number) => number.toString().padLeft(2, '0');
    return '${value.year}-${two(value.month)}-${two(value.day)}';
  }
}

class _SourceFilterChip extends StatelessWidget {
  const _SourceFilterChip({
    required this.label,
    required this.selected,
    required this.color,
    required this.onTap,
  });

  final String label;
  final bool selected;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(999.r),
        child: Ink(
          padding: EdgeInsets.symmetric(horizontal: 12.w, vertical: 8.h),
          decoration: BoxDecoration(
            color: selected
                ? color.withValues(alpha: 0.18)
                : const Color(0xFF101010),
            borderRadius: BorderRadius.circular(999.r),
            border: Border.all(
              color: selected
                  ? color.withValues(alpha: 0.7)
                  : Colors.white.withValues(alpha: 0.08),
            ),
          ),
          child: Text(
            label,
            style: TextStyle(
              fontSize: 12.sp,
              fontWeight: selected ? FontWeight.w600 : FontWeight.w400,
              color: selected ? Colors.white : Colors.white70,
            ),
          ),
        ),
      ),
    );
  }
}
