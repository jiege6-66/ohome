import 'package:get/get.dart';

import '../controllers/drops_controller.dart';

class DropsBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<DropsController>(() => DropsController());
  }
}
