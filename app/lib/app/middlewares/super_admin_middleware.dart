import 'package:flutter/widgets.dart';
import 'package:get/get.dart';

import '../services/auth_service.dart';

class SuperAdminMiddleware extends GetMiddleware {
  SuperAdminMiddleware({int priority = 2}) : _priority = priority;

  final int _priority;

  @override
  int? get priority => _priority;

  @override
  RouteSettings? redirect(String? route) {
    final auth = Get.find<AuthService>();
    if (!auth.isLoggedIn) {
      return const RouteSettings(name: '/login');
    }
    if (auth.user.value?.isSuperAdmin == true) {
      return null;
    }
    if (Get.overlayContext != null) {
      if (Get.isSnackbarOpen) {
        Get.closeCurrentSnackbar();
      }
      Get.snackbar('提示', '仅超级管理员可访问');
    }
    return const RouteSettings(name: '/main');
  }
}
