import 'package:get/get.dart';

import '../controllers/plugin_controller.dart';

class PluginBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<PluginController>(
      () => PluginController(),
    );
  }
}
