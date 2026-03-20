import 'package:dio/dio.dart';
import 'package:get/get.dart';

import '../../../data/api/config.dart';
import '../../../data/models/config_model.dart';
import '../../../data/models/config_upsert_payload.dart';
import '../views/quark_cookie_web_login_view.dart';

class QuarkLoginController extends GetxController {
  QuarkLoginController({ConfigApi? configApi})
    : _configApi = configApi ?? Get.find<ConfigApi>();

  static const String _quarkCookiesKey = 'quark_cookies';
  static const String _quarkLogoutUrl =
      'https://pan.quark.cn/account/logout?callback=';
  static const String _quarkWebReferer = 'https://pan.quark.cn/';
  static const String _quarkUserAgent =
      'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 '
      '(KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36';

  final ConfigApi _configApi;

  final loading = false.obs;
  final saving = false.obs;
  final currentConfig = Rxn<ConfigModel>();

  @override
  void onInit() {
    super.onInit();
    loadConfig();
  }

  DateTime? get updatedAt => currentConfig.value?.updatedAt;

  bool get configured => currentConfig.value?.value.trim().isNotEmpty == true;

  Future<void> loadConfig() async {
    loading.value = true;
    try {
      final config = await _configApi.findConfigByKey(_quarkCookiesKey);
      currentConfig.value = config;
    } catch (_) {
      return;
    } finally {
      loading.value = false;
    }
  }

  Future<void> openWebLogin() async {
    if (saving.value) return;
    final shouldResetSession = configured;
    if (shouldResetSession) {
      await _logoutCurrentCookies();
    }
    final result = await Get.to<String>(
      () => QuarkCookieWebLoginView(resetSessionOnOpen: shouldResetSession),
    );
    final value = result?.trim() ?? '';
    if (value.isEmpty) {
      return;
    }
    await saveCookies(value);
  }

  Future<void> _logoutCurrentCookies() async {
    final existing = currentConfig.value;
    final cookie = existing?.value.trim() ?? '';
    if (cookie.isEmpty) {
      return;
    }

    saving.value = true;
    try {
      final dio = Dio(
        BaseOptions(
          connectTimeout: const Duration(seconds: 15),
          receiveTimeout: const Duration(seconds: 15),
          followRedirects: false,
          validateStatus: (status) => status != null && status < 400,
          headers: <String, String>{
            'Cookie': cookie,
            'Referer': _quarkWebReferer,
            'User-Agent': _quarkUserAgent,
          },
        ),
      );
      await dio.get<void>(_quarkLogoutUrl);
    } catch (_) {
      // Ignore logout request errors and still clear the locally stored cookie.
    } finally {
      await _clearSavedCookies();
      saving.value = false;
    }
  }

  Future<void> _clearSavedCookies() async {
    final existing = currentConfig.value;
    if (existing == null) {
      currentConfig.value = null;
      return;
    }
    await _configApi.saveConfig(
      ConfigUpsertPayload.fromConfig(existing, value: ''),
    );
    await loadConfig();
  }

  Future<void> saveCookies(String value) async {
    final normalized = value.trim();
    if (normalized.isEmpty) {
      Get.snackbar('提示', '未检测到可保存的夸克 Cookies');
      return;
    }

    saving.value = true;
    try {
      final existing = currentConfig.value;
      final payload = existing != null
          ? ConfigUpsertPayload.fromConfig(existing, value: normalized)
          : ConfigUpsertPayload(
              name: '夸克cookies',
              key: _quarkCookiesKey,
              value: normalized,
              isLock: '1',
              remark: '',
            );
      await _configApi.saveConfig(payload);
      Get.snackbar('提示', '夸克网页登录 Cookies 已保存');
      await loadConfig();
    } catch (_) {
      return;
    } finally {
      saving.value = false;
    }
  }
}
