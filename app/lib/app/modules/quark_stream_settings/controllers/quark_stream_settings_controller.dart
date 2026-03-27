import 'package:get/get.dart';

import '../../../data/api/config.dart';
import '../../../data/models/config_model.dart';
import '../../../data/models/config_upsert_payload.dart';

class QuarkStreamSettingsController extends GetxController {
  QuarkStreamSettingsController({ConfigApi? configApi})
    : _configApi = configApi ?? Get.find<ConfigApi>();

  static const String webProxyModeKey = 'quark_fs_web_proxy_mode';

  static const String defaultWebProxyMode = 'native_proxy';

  final ConfigApi _configApi;

  final loading = false.obs;
  final saving = false.obs;
  final configs = <String, ConfigModel>{}.obs;
  final selectedMode = defaultWebProxyMode.obs;

  @override
  void onInit() {
    super.onInit();
    loadConfigs();
  }

  Future<void> loadConfigs() async {
    loading.value = true;
    try {
      final result = await _configApi.findConfigsByKeys(_configKeys);
      configs.assignAll(result);
      selectedMode.value = _normalizeMode(
        result[webProxyModeKey]?.value ?? defaultWebProxyMode,
      );
    } finally {
      loading.value = false;
    }
  }

  Future<void> save() async {
    if (saving.value) return;

    saving.value = true;
    try {
      await _saveSingle(
        key: webProxyModeKey,
        name: '夸克播放代理模式',
        remark: '夸克在线播放代理模式：native_proxy=本地代理，302_redirect=302直连',
        value: selectedMode.value,
      );
      await loadConfigs();
      Get.snackbar('提示', '夸克播放配置已保存，新的播放请求会立即生效');
    } finally {
      saving.value = false;
    }
  }

  DateTime? updatedAtFor(String key) => configs[key]?.updatedAt;

  Future<void> _saveSingle({
    required String key,
    required String name,
    required String remark,
    required String value,
  }) async {
    final existing = configs[key];
    final payload = existing != null
        ? ConfigUpsertPayload.fromConfig(existing, value: value)
        : ConfigUpsertPayload(
            name: name,
            key: key,
            value: value,
            isLock: '1',
            remark: remark,
          );
    await _configApi.saveConfig(payload);
  }

  String _normalizeMode(String raw) {
    switch (raw.trim().toLowerCase()) {
      case 'redirect':
      case '302':
      case '302_redirect':
      case 'direct':
        return '302_redirect';
      default:
        return defaultWebProxyMode;
    }
  }

  List<String> get _configKeys => const <String>[webProxyModeKey];
}
