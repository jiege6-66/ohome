import 'package:shared_preferences/shared_preferences.dart';

class SearchHistoryStorage {
  static const _key = 'search_history_keywords';
  SharedPreferences? _prefs;

  Future<SharedPreferences> _instance() async {
    final cached = _prefs;
    if (cached != null) return cached;
    final prefs = await SharedPreferences.getInstance();
    _prefs = prefs;
    return prefs;
  }

  Future<List<String>> readAll() async {
    final prefs = await _instance();
    final list = prefs.getStringList(_key) ?? const <String>[];
    return list
        .map((item) => item.trim())
        .where((item) => item.isNotEmpty)
        .toList(growable: false);
  }

  Future<void> writeAll(List<String> keywords) async {
    final prefs = await _instance();
    final values = keywords
        .map((item) => item.trim())
        .where((item) => item.isNotEmpty)
        .toList(growable: false);
    await prefs.setStringList(_key, values);
  }

  Future<void> clear() async {
    final prefs = await _instance();
    await prefs.remove(_key);
  }
}
