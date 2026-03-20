import 'package:get/get.dart';

import '../controllers/quark_sync_controller.dart';

class QuarkSyncBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<QuarkSyncController>(() => QuarkSyncController());
  }
}
