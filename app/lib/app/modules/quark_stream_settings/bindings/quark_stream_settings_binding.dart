import 'package:get/get.dart';

import '../controllers/quark_stream_settings_controller.dart';

class QuarkStreamSettingsBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<QuarkStreamSettingsController>(
      () => QuarkStreamSettingsController(),
    );
  }
}
