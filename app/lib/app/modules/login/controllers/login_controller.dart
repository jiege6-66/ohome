import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:ohome/app/data/models/discovered_server.dart';
import 'package:ohome/app/routes/app_pages.dart';
import 'package:ohome/app/services/discovery_service.dart';
import 'package:ohome/app/utils/app_env.dart';
import 'package:ohome/app/utils/http_client.dart';

import '../../../services/auth_service.dart';

class LoginController extends GetxController {
  LoginController({required DiscoveryService discoveryService})
    : _discoveryService = discoveryService;

  final GlobalKey<FormState> loginFormKey = GlobalKey<FormState>();
  final DiscoveryService _discoveryService;

  final autoValidateMode = AutovalidateMode.disabled.obs;

  late TextEditingController apiBaseUrlController;
  late TextEditingController nameController;
  late TextEditingController passwordController;

  final isLoading = false.obs;
  final isDiscovering = false.obs;
  final discoveryErrorMessage = RxnString();
  final discoveredServers = <DiscoveredServer>[].obs;
  final selectedServer = Rxn<DiscoveredServer>();

  var _manualApiBaseUrlEdited = false;
  var _isApplyingSelection = false;
  var _hasUserSelectedServer = false;

  @override
  void onInit() {
    super.onInit();
    apiBaseUrlController = TextEditingController(
      text: AppEnv.instance.apiBaseUrlInputValue,
    );
    apiBaseUrlController.addListener(_handleApiBaseUrlChanged);
    nameController = TextEditingController();
    passwordController = TextEditingController();
    Future<void>.microtask(refreshDiscovery);
  }

  String? validateApiBaseUrl(String? value) {
    try {
      AppEnv.normalizeApiBaseUrlInput(value ?? '');
      return null;
    } on FormatException catch (error) {
      return error.message.toString();
    }
  }

  String? validateName(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入用户名';
    }
    return null;
  }

  String? validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return '请输入密码';
    }
    return null;
  }

  Future<void> login() async {
    autoValidateMode.value = AutovalidateMode.onUserInteraction;
    final apiBaseUrlError = validateApiBaseUrl(apiBaseUrlController.text);
    if (apiBaseUrlError != null) {
      Get.snackbar('提示', apiBaseUrlError, duration: const Duration(seconds: 2));
      return;
    }
    if (loginFormKey.currentState!.validate()) {
      loginFormKey.currentState!.save();
      try {
        isLoading.value = true;
        await AppEnv.instance.updateApiBaseUrl(apiBaseUrlController.text);
        HttpClient.syncBaseUrl();
        final auth = Get.find<AuthService>();
        await auth.login(
          name: nameController.text.trim(),
          password: passwordController.text,
        );
        await _discoveryService.rememberSuccessfulServer(
          apiBaseUrlInput: apiBaseUrlController.text,
          selectedServer: selectedServer.value,
        );
        Get.offAllNamed(Routes.MAIN);
      } finally {
        isLoading.value = false;
      }
    }
  }

  Future<void> refreshDiscovery() async {
    try {
      isDiscovering.value = true;
      discoveryErrorMessage.value = null;

      final servers = await _discoveryService.discoverServers();
      discoveredServers.assignAll(servers);

      final selected = selectedServer.value;
      if (selected != null) {
        final stillExists = servers.any(
          (server) =>
              server.instanceId == selected.instanceId &&
              server.origin == selected.origin,
        );
        if (!stillExists) {
          selectedServer.value = null;
        }
      }

      if (servers.isEmpty) {
        return;
      }

      if (servers.length == 1 &&
          !_manualApiBaseUrlEdited &&
          !_hasUserSelectedServer) {
        _applySelectedServer(servers.first, userInitiated: false);
        return;
      }

      if (_manualApiBaseUrlEdited || _hasUserSelectedServer) {
        return;
      }
    } catch (error) {
      discoveryErrorMessage.value = error.toString();
    } finally {
      isDiscovering.value = false;
    }
  }

  void selectDiscoveredServer(DiscoveredServer server) {
    _applySelectedServer(server, userInitiated: true);
  }

  bool get hasFoundServer =>
      selectedServer.value != null || discoveredServers.isNotEmpty;

  String get serverStatusText {
    final selected = selectedServer.value;
    if (selected != null) {
      return selected.serviceName;
    }
    if (isDiscovering.value) {
      return '正在查找局域网服务';
    }
    if (discoveredServers.isNotEmpty) {
      return '已找到 ${discoveredServers.length} 个局域网服务';
    }
    return '手动输入或重新扫描';
  }

  void _applySelectedServer(
    DiscoveredServer server, {
    required bool userInitiated,
  }) {
    selectedServer.value = server;
    _isApplyingSelection = true;
    apiBaseUrlController.text = server.origin;
    _isApplyingSelection = false;
    if (userInitiated) {
      _manualApiBaseUrlEdited = false;
      _hasUserSelectedServer = true;
    }
  }

  void _handleApiBaseUrlChanged() {
    if (_isApplyingSelection) {
      return;
    }

    final current = apiBaseUrlController.text.trim();
    final selected = selectedServer.value;
    if (selected != null && current == selected.origin) {
      return;
    }

    _manualApiBaseUrlEdited = true;
    _hasUserSelectedServer = false;
    selectedServer.value = null;
  }

  @override
  void onClose() {
    apiBaseUrlController.removeListener(_handleApiBaseUrlChanged);
    apiBaseUrlController.dispose();
    nameController.dispose();
    passwordController.dispose();
    super.onClose();
  }
}
