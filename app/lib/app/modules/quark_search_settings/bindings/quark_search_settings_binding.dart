import 'package:get/get.dart';

import '../controllers/quark_search_settings_controller.dart';

class QuarkSearchSettingsBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<QuarkSearchSettingsController>(
      () => QuarkSearchSettingsController(),
    );
  }
}
