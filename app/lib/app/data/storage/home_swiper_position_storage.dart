import 'package:shared_preferences/shared_preferences.dart';

class HomeSwiperPositionStorage {
  static const _key = 'home_swiper_index';
  SharedPreferences? _prefs;

  Future<SharedPreferences> _instance() async {
    final cached = _prefs;
    if (cached != null) return cached;
    final prefs = await SharedPreferences.getInstance();
    _prefs = prefs;
    return prefs;
  }

  Future<void> save(int index) async {
    final prefs = await _instance();
    await prefs.setInt(_key, index);
  }

  Future<int?> read() async {
    final prefs = await _instance();
    return prefs.getInt(_key);
  }
}
