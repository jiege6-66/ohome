import 'package:get/get.dart';

import '../controllers/playlet_controller.dart';

class PlayLetBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<PlayLetController>(() => PlayLetController());
  }
}
