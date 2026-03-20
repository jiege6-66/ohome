import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/api/app_message.dart';
import '../../../data/models/app_message_model.dart';
import '../../../data/models/app_message_push_event.dart';
import '../../../routes/app_pages.dart';
import '../../../services/app_message_push_service.dart';
import '../../drops/drops_catalog.dart';
import '../../drops/controllers/drops_controller.dart';
import '../../drops/views/drops_event_form_view.dart';
import '../../drops/views/drops_item_detail_view.dart';

class MessagesController extends GetxController {
  MessagesController({
    AppMessageApi? appMessageApi,
    AppMessagePushService? appMessagePushService,
  }) : _appMessageApi = appMessageApi ?? Get.find<AppMessageApi>(),
       _appMessagePushService =
           appMessagePushService ??
           (Get.isRegistered<AppMessagePushService>()
               ? Get.find<AppMessagePushService>()
               : null);

  static const int _pageSize = 20;

  final AppMessageApi _appMessageApi;
  final AppMessagePushService? _appMessagePushService;

  final scrollController = ScrollController();
  final messages = <AppMessageModel>[].obs;
  final loading = false.obs;
  final loadingMore = false.obs;
  final hasMore = true.obs;
  final unreadOnly = false.obs;
  final unreadCount = 0.obs;
  final sourceFilter = ''.obs;
  final deletingMessageIds = <int>{}.obs;

  int _page = 1;
  int _token = 0;
  bool _pendingPushSync = false;
  bool _syncingFromPush = false;
  StreamSubscription<AppMessagePushEvent>? _pushSubscription;
  Timer? _pushSyncTimer;

  @override
  void onInit() {
    super.onInit();
    scrollController.addListener(_handleScroll);
    _pushSubscription = _appMessagePushService?.events.listen(_handlePushEvent);
    loadMessages(refresh: true);
  }

  @override
  void onClose() {
    _pushSubscription?.cancel();
    _pushSyncTimer?.cancel();
    scrollController.dispose();
    super.onClose();
  }

  Future<void> loadMessages({
    required bool refresh,
    bool showErrorToast = true,
  }) async {
    late final int token;
    if (refresh) {
      token = ++_token;
      _page = 1;
      hasMore.value = true;
      loading.value = true;
      loadingMore.value = false;
    } else {
      if (loading.value || loadingMore.value || !hasMore.value) return;
      token = _token;
      loadingMore.value = true;
    }

    try {
      final result = await _appMessageApi.getMessageList(
        source: sourceFilter.value.isEmpty ? null : sourceFilter.value,
        readOnly: unreadOnly.value ? false : null,
        page: _page,
        limit: _pageSize,
        showErrorToast: showErrorToast,
      );
      if (token != _token) return;
      unreadCount.value = result.unreadCount;
      if (refresh) {
        messages.assignAll(result.records);
      } else {
        messages.addAll(result.records);
      }
      hasMore.value = messages.length < result.total;
      if (hasMore.value) {
        _page += 1;
      }
    } finally {
      if (token == _token) {
        if (refresh) {
          loading.value = false;
        } else {
          loadingMore.value = false;
        }
        if (_pendingPushSync && !loading.value && !loadingMore.value) {
          _schedulePushSync();
        }
      }
    }
  }

  void toggleUnreadOnly(bool value) {
    if (unreadOnly.value == value) return;
    unreadOnly.value = value;
    loadMessages(refresh: true);
  }

  void changeSourceFilter(String value) {
    final normalized = value.trim().toLowerCase();
    if (sourceFilter.value == normalized) return;
    sourceFilter.value = normalized;
    loadMessages(refresh: true);
  }

  Future<void> markAllRead() async {
    await _appMessageApi.markAllMessagesRead();
    await loadMessages(refresh: true);
    _refreshDropsOverviewIfNeeded();
    Get.snackbar('提示', '消息已全部标记为已读');
  }

  bool isDeleting(int? id) => id != null && deletingMessageIds.contains(id);

  Future<void> deleteMessage(AppMessageModel message) async {
    final id = message.id;
    if (id == null || isDeleting(id)) return;

    final confirmed = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('删除消息'),
        content: Text(
          message.title.isEmpty ? '确定删除这条消息吗？' : '确定删除「${message.title}」吗？',
        ),
        actions: [
          TextButton(
            onPressed: () => Get.back(result: false),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Get.back(result: true),
            child: const Text('删除'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    deletingMessageIds.add(id);
    try {
      await _appMessageApi.deleteMessage(id);
      await loadMessages(refresh: true);
      _refreshDropsOverviewIfNeeded();
      Get.snackbar('提示', '消息已删除');
    } catch (_) {
      return;
    } finally {
      deletingMessageIds.remove(id);
    }
  }

  Future<void> openMessage(AppMessageModel message) async {
    final id = message.id;
    if (id != null && !message.read) {
      await _appMessageApi.markMessageRead(id);
    }

    if (message.source == 'drops' &&
        message.bizType == 'item' &&
        message.bizId != null) {
      await Get.to<bool>(() => DropsItemDetailView(itemId: message.bizId!));
    } else if (message.source == 'drops' &&
        message.bizType == 'event' &&
        message.bizId != null) {
      await Get.to<bool>(
        () => DropsEventFormView(initialEventId: message.bizId),
      );
    } else if (message.source == 'quark') {
      await Get.toNamed(Routes.QUARK_TRANSFER_TASKS);
    }

    await loadMessages(refresh: true);
    _refreshDropsOverviewIfNeeded();
  }

  String sourceLabelOf(AppMessageModel message) {
    return sourceLabel(message.source);
  }

  String sourceLabel(String source) {
    switch (source.trim().toLowerCase()) {
      case 'drops':
        return '点滴';
      case 'quark':
        return '夸克';
      case 'system':
        return '系统';
      default:
        return '通知';
    }
  }

  String sourceFilterLabel(String source) {
    switch (source.trim().toLowerCase()) {
      case '':
        return '全部';
      case 'drops':
        return '点滴';
      case 'quark':
        return '夸克';
      case 'system':
        return '系统';
      default:
        return sourceLabel(source);
    }
  }

  Color sourceColorOf(AppMessageModel message) {
    return sourceColor(message.source);
  }

  Color sourceColor(String source) {
    switch (source.trim().toLowerCase()) {
      case 'drops':
        return const Color(0xFFBB86FC);
      case 'quark':
        return const Color(0xFF42A5F5);
      case 'system':
        return const Color(0xFF26A69A);
      default:
        return const Color(0xFF9E9E9E);
    }
  }

  String get emptyStateText {
    switch (sourceFilter.value) {
      case 'drops':
        return '暂无点滴消息';
      case 'quark':
        return '暂无夸克消息';
      case 'system':
        return '暂无系统消息';
      default:
        return '暂无消息通知';
    }
  }

  String summaryOf(AppMessageModel message) {
    final summary = message.summary;
    if (summary.isEmpty) {
      return '点击查看详情';
    }
    if (message.source.trim().toLowerCase() != 'drops') {
      return summary;
    }
    return _translateDropsSummary(summary);
  }

  void _handleScroll() {
    if (!scrollController.hasClients) return;
    final position = scrollController.position;
    if (position.maxScrollExtent - position.pixels < 200) {
      loadMessages(refresh: false);
    }
  }

  void _handlePushEvent(AppMessagePushEvent event) {
    final unread = event.unreadCount;
    if (unreadCount.value != unread) {
      unreadCount.value = unread;
    }
    _schedulePushSync();
  }

  void _schedulePushSync() {
    if (isClosed) return;
    _pushSyncTimer?.cancel();
    _pushSyncTimer = Timer(const Duration(milliseconds: 250), () {
      _pushSyncTimer = null;
      unawaited(_syncMessagesFromPush());
    });
  }

  Future<void> _syncMessagesFromPush() async {
    if (_syncingFromPush || loading.value || loadingMore.value) {
      _pendingPushSync = true;
      return;
    }

    _pendingPushSync = false;
    _syncingFromPush = true;
    final token = ++_token;
    final currentLimit = messages.length > _pageSize
        ? messages.length
        : _pageSize;

    try {
      final result = await _appMessageApi.getMessageList(
        source: sourceFilter.value.isEmpty ? null : sourceFilter.value,
        readOnly: unreadOnly.value ? false : null,
        page: 1,
        limit: currentLimit,
        showErrorToast: false,
      );
      if (token != _token) return;

      unreadCount.value = result.unreadCount;
      messages.assignAll(result.records);
      hasMore.value = messages.length < result.total;
      final loadedPages = (messages.length / _pageSize).ceil();
      _page = loadedPages + 1;
    } finally {
      if (token == _token) {
        _syncingFromPush = false;
        if (_pendingPushSync) {
          _schedulePushSync();
        }
      }
    }
  }

  void _refreshDropsOverviewIfNeeded() {
    if (Get.isRegistered<DropsController>()) {
      Get.find<DropsController>().refreshOverview();
    }
  }

  String _translateDropsSummary(String summary) {
    var result = summary;
    final replacements = <String, String>{
      'birthday': dropsEventTypeLabel('birthday'),
      'anniversary': dropsEventTypeLabel('anniversary'),
      'custom': dropsEventTypeLabel('custom'),
      'solar': dropsCalendarLabel('solar'),
      'lunar': dropsCalendarLabel('lunar'),
      'kitchen': dropsCategoryLabel('kitchen'),
      'food': dropsCategoryLabel('food'),
      'medicine': dropsCategoryLabel('medicine'),
      'clothing': dropsCategoryLabel('clothing'),
      'other': dropsCategoryLabel('other'),
    };

    replacements.forEach((raw, label) {
      result = result.replaceAll(RegExp('\\b${RegExp.escape(raw)}\\b'), label);
    });
    return result;
  }
}
