import 'package:shared_preferences/shared_preferences.dart';

class VolumeSettingsStorage {
  static const _key = 'music_volume';
  SharedPreferences? _prefs;

  Future<SharedPreferences> _instance() async {
    final cached = _prefs;
    if (cached != null) return cached;
    final prefs = await SharedPreferences.getInstance();
    _prefs = prefs;
    return prefs;
  }

  Future<void> save(double volume) async {
    final prefs = await _instance();
    await prefs.setDouble(_key, volume.clamp(0.0, 1.0));
  }

  Future<double?> read() async {
    final prefs = await _instance();
    final value = prefs.getDouble(_key);
    if (value == null) return null;
    return value.clamp(0.0, 1.0);
  }
}
