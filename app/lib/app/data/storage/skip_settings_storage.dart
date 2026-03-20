import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

class SkipSettings {
  const SkipSettings({required this.intro, required this.outro});

  final Duration intro;
  final Duration outro;

  Map<String, dynamic> toJson() => <String, dynamic>{
    'introSeconds': intro.inSeconds,
    'outroSeconds': outro.inSeconds,
  };

  factory SkipSettings.fromJson(Map<String, dynamic> json) {
    int toInt(dynamic value) {
      if (value is int) return value;
      if (value is String) return int.tryParse(value) ?? 0;
      return 0;
    }

    final introSeconds = toInt(json['introSeconds']);
    final outroSeconds = toInt(json['outroSeconds']);
    return SkipSettings(
      intro: introSeconds > 0 ? Duration(seconds: introSeconds) : Duration.zero,
      outro: outroSeconds > 0 ? Duration(seconds: outroSeconds) : Duration.zero,
    );
  }
}

class SkipSettingsStorage {
  SkipSettingsStorage({String keyPrefix = 'skip_settings_'})
    : _keyPrefix = keyPrefix;

  final String _keyPrefix;
  SharedPreferences? _prefs;

  Future<SharedPreferences> _instance() async {
    final cached = _prefs;
    if (cached != null) return cached;
    final prefs = await SharedPreferences.getInstance();
    _prefs = prefs;
    return prefs;
  }

  Future<void> save(String folderPath, SkipSettings settings) async {
    final folder = folderPath.trim();
    if (folder.isEmpty) return;
    final prefs = await _instance();
    await prefs.setString(_keyFor(folder), jsonEncode(settings.toJson()));
  }

  Future<SkipSettings?> read(String folderPath) async {
    final folder = folderPath.trim();
    if (folder.isEmpty) return null;
    final prefs = await _instance();
    final raw = prefs.getString(_keyFor(folder));
    if (raw == null || raw.trim().isEmpty) return null;
    try {
      final json = jsonDecode(raw);
      if (json is Map<String, dynamic>) return SkipSettings.fromJson(json);
    } catch (_) {
      return null;
    }
    return null;
  }

  Future<void> clear(String folderPath) async {
    final folder = folderPath.trim();
    if (folder.isEmpty) return;
    final prefs = await _instance();
    await prefs.remove(_keyFor(folder));
  }

  String _keyFor(String folderPath) {
    final normalized = folderPath.replaceAll('\\', '/').trim();
    return _keyPrefix + base64Url.encode(utf8.encode(normalized));
  }
}
