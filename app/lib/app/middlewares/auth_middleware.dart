import 'package:flutter/widgets.dart';
import 'package:get/get.dart';

import '../services/auth_service.dart';

class AuthMiddleware extends GetMiddleware {
  AuthMiddleware({int priority = 1}) : _priority = priority;

  final int _priority;

  @override
  int? get priority => _priority;

  @override
  RouteSettings? redirect(String? route) {
    final auth = Get.find<AuthService>();
    if (!auth.isLoggedIn) {
      return const RouteSettings(name: '/login');
    }
    return null;
  }
}
