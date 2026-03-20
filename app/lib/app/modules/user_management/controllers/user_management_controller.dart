import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/api/user.dart';
import '../../../data/models/user_model.dart';
import '../../../services/auth_service.dart';
import '../views/user_management_form_view.dart';

class UserManagementController extends GetxController {
  UserManagementController({UserApi? userApi, AuthService? authService})
    : _userApi = userApi ?? Get.find<UserApi>(),
      _authService = authService ?? Get.find<AuthService>();

  static const int _pageSize = 20;

  final UserApi _userApi;
  final AuthService _authService;

  final nameController = TextEditingController();
  final scrollController = ScrollController();

  final users = <UserModel>[].obs;
  final loading = false.obs;
  final loadingMore = false.obs;
  final hasMore = true.obs;
  final _keyword = ''.obs;

  Worker? _searchWorker;

  int _page = 1;
  int _loadToken = 0;

  int? get currentUserId => _authService.user.value?.id;

  bool isCurrentUser(UserModel user) =>
      user.id != null && user.id == currentUserId;

  bool canDeleteUser(UserModel user) => !isCurrentUser(user);

  @override
  void onInit() {
    super.onInit();
    _searchWorker = debounce<String>(
      _keyword,
      (_) => loadUsers(refresh: true),
      time: const Duration(milliseconds: 250),
    );
    scrollController.addListener(_handleScroll);
    loadUsers(refresh: true);
  }

  @override
  void onClose() {
    _searchWorker?.dispose();
    nameController.dispose();
    scrollController.dispose();
    super.onClose();
  }

  Future<void> loadUsers({required bool refresh}) async {
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
      final result = await _userApi.getUserList(
        name: nameController.text,
        page: _page,
        limit: _pageSize,
      );
      if (token != _loadToken) return;

      if (refresh) {
        users.assignAll(result.records);
      } else {
        users.addAll(result.records);
      }

      hasMore.value = users.length < result.total;
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

  Future<void> search() => loadUsers(refresh: true);

  void onKeywordChanged(String value) {
    final keyword = value.trim();
    if (_keyword.value == keyword) {
      return;
    }
    _keyword.value = keyword;
  }

  Future<void> openCreatePage() async {
    final changed = await Get.to<bool>(() => const UserManagementFormView());
    if (changed == true) {
      await loadUsers(refresh: true);
    }
  }

  Future<void> openEditPage(UserModel user) async {
    final changed = await Get.to<bool>(
      () => UserManagementFormView(initialUser: user),
    );
    if (changed == true) {
      if (_authService.user.value?.isSuperAdmin != true) {
        Get.back<void>();
        Get.snackbar('提示', '角色已更新，当前账号不再具有用户管理权限');
        return;
      }
      await loadUsers(refresh: true);
    }
  }

  Future<void> confirmDelete(UserModel user) async {
    final id = user.id;
    if (id == null) return;
    if (!canDeleteUser(user)) {
      Get.snackbar('提示', '不能删除当前登录账号');
      return;
    }

    final confirmed = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('确认删除用户'),
        content: Text('删除 ${user.name.isEmpty ? '该用户' : user.name} 后不可恢复。'),
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

    try {
      await _userApi.deleteUser(id);
      Get.snackbar('提示', '用户已删除');
      await loadUsers(refresh: true);
    } catch (_) {
      return;
    }
  }

  Future<void> confirmResetPassword(UserModel user) async {
    final id = user.id;
    if (id == null) return;

    final confirmed = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('确认重置密码'),
        content: const Text('密码将重置为系统默认密码。'),
        actions: [
          TextButton(
            onPressed: () => Get.back(result: false),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Get.back(result: true),
            child: const Text('确认'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    try {
      await _userApi.resetPassword(id);
      Get.snackbar('提示', '密码已重置为系统默认密码');
    } catch (_) {
      return;
    }
  }

  void _handleScroll() {
    if (!scrollController.hasClients) return;
    final position = scrollController.position;
    if (position.maxScrollExtent - position.pixels <= 200) {
      loadUsers(refresh: false);
    }
  }
}
