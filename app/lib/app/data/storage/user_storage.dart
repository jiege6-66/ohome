import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import '../models/user_model.dart';

class UserStorage {
  static const _kUserInfo = 'userInfo';

  Future<UserModel?> readUser() async {
    final prefs = await SharedPreferences.getInstance();
    final raw = prefs.getString(_kUserInfo);
    if (raw == null || raw.trim().isEmpty) return null;
    final decoded = jsonDecode(raw);
    if (decoded is! Map<String, dynamic>) return null;
    return UserModel.fromJson(decoded);
  }

  Future<void> writeUser(UserModel user) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_kUserInfo, jsonEncode(user.toJson()));
  }

  Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_kUserInfo);
  }
}
