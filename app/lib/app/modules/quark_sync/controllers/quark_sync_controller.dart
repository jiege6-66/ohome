import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/api/quark_auto_save_task.dart';
import '../../../data/models/quark_auto_save_task_model.dart';
import '../views/quark_sync_form_view.dart';

class QuarkSyncController extends GetxController {
  QuarkSyncController({QuarkAutoSaveTaskApi? taskApi})
    : _taskApi = taskApi ?? Get.find<QuarkAutoSaveTaskApi>();

  static const int _pageSize = 20;

  final QuarkAutoSaveTaskApi _taskApi;

  final taskNameController = TextEditingController();
  final scrollController = ScrollController();

  final tasks = <QuarkAutoSaveTaskModel>[].obs;
  final loading = false.obs;
  final loadingMore = false.obs;
  final hasMore = true.obs;
  final actioningTaskIds = <int>[].obs;

  int _page = 1;
  int _loadToken = 0;

  @override
  void onInit() {
    super.onInit();
    scrollController.addListener(_handleScroll);
    loadTasks(refresh: true);
  }

  @override
  void onClose() {
    taskNameController.dispose();
    scrollController.dispose();
    super.onClose();
  }

  bool isTaskBusy(int? id) => id != null && actioningTaskIds.contains(id);

  Future<void> loadTasks({required bool refresh}) async {
    late final int token;
    if (refresh) {
      token = ++_loadToken;
      _page = 1;
      hasMore.value = true;
      loading.value = true;
      loadingMore.value = false;
    } else {
      if (loading.value || loadingMore.value || !hasMore.value) {
        return;
      }
      token = _loadToken;
      loadingMore.value = true;
    }

    try {
      final result = await _taskApi.getTaskList(
        taskName: taskNameController.text,
        page: _page,
        limit: _pageSize,
      );
      if (token != _loadToken) return;

      if (refresh) {
        tasks.assignAll(result.records);
      } else {
        tasks.addAll(result.records);
      }

      hasMore.value = tasks.length < result.total;
      if (hasMore.value) {
        _page += 1;
      }
    } catch (_) {
      return;
    } finally {
      if (token == _loadToken) {
        if (refresh) {
          loading.value = false;
        } else {
          loadingMore.value = false;
        }
      }
    }
  }

  Future<void> search() => loadTasks(refresh: true);

  Future<void> resetFilters() async {
    taskNameController.clear();
    await loadTasks(refresh: true);
  }

  Future<void> openCreatePage() async {
    final changed = await Get.to<bool>(() => const QuarkSyncFormView());
    if (changed == true) {
      await loadTasks(refresh: true);
    }
  }

  Future<void> openEditPage(QuarkAutoSaveTaskModel task) async {
    final changed = await Get.to<bool>(
      () => QuarkSyncFormView(initialTask: task),
    );
    if (changed == true) {
      await loadTasks(refresh: true);
    }
  }

  Future<void> confirmDelete(QuarkAutoSaveTaskModel task) async {
    final id = task.id;
    if (id == null || isTaskBusy(id)) return;

    final confirmed = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('确认删除任务'),
        content: Text(
          '删除 ${task.taskName.isEmpty ? '该任务' : task.taskName} 后不可恢复。',
        ),
        actions: [
          TextButton(
            onPressed: () => Get.back(result: false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Get.back(result: true),
            child: const Text('删除'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    _markBusy(id);
    try {
      await _taskApi.deleteTask(id);
      Get.snackbar('提示', '任务已删除');
      await loadTasks(refresh: true);
    } catch (_) {
      return;
    } finally {
      _clearBusy(id);
    }
  }

  Future<void> runTaskOnce(QuarkAutoSaveTaskModel task) async {
    final id = task.id;
    if (id == null || isTaskBusy(id)) return;

    final confirmed = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('执行一次'),
        content: Text(
          '立即执行 ${task.taskName.isEmpty ? '该任务' : task.taskName} 的转存同步？',
        ),
        actions: [
          TextButton(
            onPressed: () => Get.back(result: false),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Get.back(result: true),
            child: const Text('执行'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    _markBusy(id);
    try {
      await _taskApi.runOnce(id);
      Get.snackbar('提示', '任务已提交执行');
      await loadTasks(refresh: true);
    } catch (_) {
      return;
    } finally {
      _clearBusy(id);
    }
  }

  void _handleScroll() {
    if (!scrollController.hasClients) return;
    final position = scrollController.position;
    if (position.maxScrollExtent - position.pixels <= 200) {
      loadTasks(refresh: false);
    }
  }

  void _markBusy(int id) {
    if (!actioningTaskIds.contains(id)) {
      actioningTaskIds.add(id);
    }
  }

  void _clearBusy(int id) {
    actioningTaskIds.remove(id);
  }
}
