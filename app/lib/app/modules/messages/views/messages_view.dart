import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';

import '../../../data/models/app_message_model.dart';
import '../../../theme/app_theme.dart';
import '../controllers/messages_controller.dart';

class MessagesView extends StatefulWidget {
  const MessagesView({super.key});

  @override
  State<MessagesView> createState() => _MessagesViewState();
}

class _MessagesViewState extends State<MessagesView> {
  late final MessagesController controller;
  late final PageController _pageController;

  @override
  void initState() {
    super.initState();
    controller = Get.find<MessagesController>();
    _pageController = PageController(
      initialPage: controller.selectedTabIndex.value,
    );
    controller.ensureTabLoaded(controller.selectedTabIndex.value);
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

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
              padding: EdgeInsets.fromLTRB(16.w, 16.h, 16.w, 12.h),
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
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: controller.tabCount,
                onPageChanged: controller.changeTab,
                itemBuilder: (context, tabIndex) =>
                    _buildList(tabIndex, bottomInset),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildList(int tabIndex, double bottomInset) {
    return Obx(() {
      final records = controller.messagesOf(tabIndex);
      final loading = controller.loadingOf(tabIndex).value;
      final loadingMore = controller.loadingMoreOf(tabIndex).value;
      final hasMore = controller.hasMoreOf(tabIndex).value;

      if (loading && records.isEmpty) {
        return const Center(child: CircularProgressIndicator());
      }

      if (records.isEmpty) {
        return RefreshIndicator(
          onRefresh: () =>
              controller.loadMessages(refresh: true, tabIndex: tabIndex),
          child: ListView(
            controller: controller.scrollControllerOf(tabIndex),
            physics: const AlwaysScrollableScrollPhysics(),
            padding: EdgeInsets.only(bottom: bottomInset),
            children: [
              SizedBox(height: 120.h),
              Center(
                child: Text(
                  controller.emptyStateTextOf(tabIndex),
                  style: const TextStyle(color: Colors.white54),
                ),
              ),
            ],
          ),
        );
      }

      return RefreshIndicator(
        onRefresh: () =>
            controller.loadMessages(refresh: true, tabIndex: tabIndex),
        child: ListView.separated(
          controller: controller.scrollControllerOf(tabIndex),
          padding: EdgeInsets.fromLTRB(16.w, 8.h, 16.w, bottomInset),
          itemCount: records.length + (hasMore ? 1 : 0),
          separatorBuilder: (_, _) => SizedBox(height: 12.h),
          itemBuilder: (_, index) {
            if (index >= records.length) {
              return loadingMore
                  ? const Padding(
                      padding: EdgeInsets.symmetric(vertical: 12),
                      child: Center(child: CircularProgressIndicator()),
                    )
                  : const SizedBox.shrink();
            }
            final message = records[index];
            return _buildMessageCard(message);
          },
        ),
      );
    });
  }

  Widget _buildMessageCard(AppMessageModel message) {
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
            border: Border.all(color: Colors.white.withValues(alpha: 0.06)),
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          width: 10.w,
                          height: 10.w,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: message.read
                                ? Colors.white24
                                : const Color(0xFFE53935),
                          ),
                        ),
                        SizedBox(width: 12.w),
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
                          final deleting = controller.isDeleting(message.id);
                          return IconButton(
                            onPressed: deleting
                                ? null
                                : () => controller.deleteMessage(message),
                            tooltip: '删除消息',
                            splashRadius: 18.r,
                            visualDensity: VisualDensity.compact,
                            padding: EdgeInsets.zero,
                            constraints: BoxConstraints.tightFor(
                              width: 28.w,
                              height: 28.w,
                            ),
                            icon: deleting
                                ? SizedBox(
                                    width: 16.w,
                                    height: 16.w,
                                    child: const CircularProgressIndicator(
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
                      style: TextStyle(fontSize: 12.sp, color: Colors.white60),
                    ),
                    SizedBox(height: 8.h),
                    Text(
                      message.triggerDate == null
                          ? '提醒时间未知'
                          : '提醒日期 ${_formatDate(message.triggerDate!)}',
                      style: TextStyle(fontSize: 11.sp, color: Colors.white38),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSourceFilter() {
    return Obx(() {
      final labels = controller.tabLabels;
      final selected = controller.selectedTabIndex.value;
      final children = <Widget>[];

      for (var index = 0; index < labels.length; index++) {
        final isSelected = selected == index;
        children.add(
          Expanded(
            child: GestureDetector(
              onTap: () {
                controller.changeTab(index);
                _pageController.animateToPage(
                  index,
                  duration: const Duration(milliseconds: 220),
                  curve: Curves.easeOutCubic,
                );
              },
              child: Stack(
                children: [
                  AnimatedContainer(
                    duration: const Duration(milliseconds: 200),
                    curve: Curves.easeOut,
                    padding: EdgeInsets.symmetric(vertical: 10.h),
                    decoration: BoxDecoration(
                      color: const Color(0xFF1E1E1E),
                      borderRadius: BorderRadius.circular(24.r),
                    ),
                    child: Center(
                      child: Text(
                        labels[index],
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.6),
                          fontWeight: FontWeight.w400,
                          fontSize: 13.sp,
                        ),
                      ),
                    ),
                  ),
                  Positioned.fill(
                    child: AnimatedOpacity(
                      opacity: isSelected ? 1.0 : 0.0,
                      duration: const Duration(milliseconds: 200),
                      curve: Curves.easeOut,
                      child: Container(
                        padding: EdgeInsets.symmetric(vertical: 10.h),
                        decoration: BoxDecoration(
                          gradient: const LinearGradient(
                            colors: [Color(0xFF7C4DFF), Color(0xFF448AFF)],
                          ),
                          borderRadius: BorderRadius.circular(24.r),
                          boxShadow: [
                            BoxShadow(
                              color: const Color(
                                0xFF7C4DFF,
                              ).withValues(alpha: 0.3),
                              blurRadius: 10.r,
                              offset: Offset(0, 3.h),
                            ),
                          ],
                        ),
                        child: Center(
                          child: Text(
                            labels[index],
                            style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.w600,
                              fontSize: 13.sp,
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
        );

        if (index < labels.length - 1) {
          children.add(SizedBox(width: 8.w));
        }
      }

      return Row(children: children);
    });
  }

  static String _formatDate(DateTime value) {
    String two(int number) => number.toString().padLeft(2, '0');
    return '${value.year}-${two(value.month)}-${two(value.day)}';
  }
}
