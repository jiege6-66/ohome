import 'package:get/get.dart';

import '../controllers/quark_transfer_tasks_controller.dart';

class QuarkTransferTasksBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<QuarkTransferTasksController>(
      () => QuarkTransferTasksController(),
    );
  }
}
