import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/api/config.dart';
import '../../../data/models/config_model.dart';
import '../../../data/models/config_upsert_payload.dart';

class QuarkStreamSettingsController extends GetxController {
  QuarkStreamSettingsController({ConfigApi? configApi})
    : _configApi = configApi ?? Get.find<ConfigApi>();

  static const String webProxyModeKey = 'quark_fs_web_proxy_mode';
  static const String concurrencyKey = 'quark_fs_concurrency';
  static const String partSizeMBKey = 'quark_fs_part_size_mb';
  static const String chunkMaxRetriesKey = 'quark_fs_chunk_max_retries';

  static const String defaultWebProxyMode = 'native_proxy';
  static const String defaultConcurrency = '3';
  static const String defaultPartSizeMB = '10';
  static const String defaultChunkMaxRetries = '3';

  final ConfigApi _configApi;

  final loading = false.obs;
  final saving = false.obs;
  final configs = <String, ConfigModel>{}.obs;
  final selectedMode = defaultWebProxyMode.obs;

  final concurrencyController = TextEditingController();
  final partSizeController = TextEditingController();
  final chunkRetriesController = TextEditingController();

  @override
  void onInit() {
    super.onInit();
    loadConfigs();
  }

  @override
  void onClose() {
    concurrencyController.dispose();
    partSizeController.dispose();
    chunkRetriesController.dispose();
    super.onClose();
  }

  Future<void> loadConfigs() async {
    loading.value = true;
    try {
      final result = await _configApi.findConfigsByKeys(_configKeys);
      configs.assignAll(result);
      selectedMode.value = _normalizeMode(
        result[webProxyModeKey]?.value ?? defaultWebProxyMode,
      );
      concurrencyController.text =
          result[concurrencyKey]?.value.trim() ?? defaultConcurrency;
      partSizeController.text =
          result[partSizeMBKey]?.value.trim() ?? defaultPartSizeMB;
      chunkRetriesController.text =
          result[chunkMaxRetriesKey]?.value.trim() ?? defaultChunkMaxRetries;
    } finally {
      loading.value = false;
    }
  }

  Future<void> save() async {
    if (saving.value) return;

    final concurrency = _parsePositiveInt(
      concurrencyController.text,
      fallback: int.parse(defaultConcurrency),
      fieldName: '并发数',
    );
    if (concurrency == null) return;

    final partSizeMB = _parsePositiveInt(
      partSizeController.text,
      fallback: int.parse(defaultPartSizeMB),
      fieldName: '分片大小',
    );
    if (partSizeMB == null) return;

    final chunkRetries = _parseNonNegativeInt(
      chunkRetriesController.text,
      fallback: int.parse(defaultChunkMaxRetries),
      fieldName: '重试次数',
    );
    if (chunkRetries == null) return;

    saving.value = true;
    try {
      await _saveSingle(
        key: webProxyModeKey,
        name: '夸克播放代理模式',
        remark: '夸克在线播放代理模式：native_proxy=本地代理，302_redirect=302直连',
        value: selectedMode.value,
      );
      await _saveSingle(
        key: concurrencyKey,
        name: '夸克播放并发数',
        remark: '夸克在线播放并发回源分片数，建议 2-4',
        value: concurrency.toString(),
      );
      await _saveSingle(
        key: partSizeMBKey,
        name: '夸克播放分片大小MB',
        remark: '夸克在线播放每个分片大小，单位 MB',
        value: partSizeMB.toString(),
      );
      await _saveSingle(
        key: chunkMaxRetriesKey,
        name: '夸克播放分片重试次数',
        remark: '夸克在线播放单个分片失败后的最大重试次数',
        value: chunkRetries.toString(),
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

  int? _parsePositiveInt(
    String raw, {
    required int fallback,
    required String fieldName,
  }) {
    final text = raw.trim();
    if (text.isEmpty) return fallback;
    final value = int.tryParse(text);
    if (value == null || value <= 0) {
      Get.snackbar('提示', '$fieldName 必须是大于 0 的整数');
      return null;
    }
    return value;
  }

  int? _parseNonNegativeInt(
    String raw, {
    required int fallback,
    required String fieldName,
  }) {
    final text = raw.trim();
    if (text.isEmpty) return fallback;
    final value = int.tryParse(text);
    if (value == null || value < 0) {
      Get.snackbar('提示', '$fieldName 不能小于 0');
      return null;
    }
    return value;
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

  List<String> get _configKeys => const <String>[
    webProxyModeKey,
    concurrencyKey,
    partSizeMBKey,
    chunkMaxRetriesKey,
  ];
}
