import 'package:get/get.dart';
import 'package:ohome/app/data/storage/discovery_storage.dart';
import 'package:ohome/app/services/discovery_service.dart';

import '../controllers/login_controller.dart';

class LoginBinding extends Bindings {
  @override
  void dependencies() {
    if (!Get.isRegistered<DiscoveryService>()) {
      Get.put<DiscoveryService>(
        DiscoveryService(storage: DiscoveryStorage()),
        permanent: true,
      );
    }

    Get.lazyPut<LoginController>(
      () => LoginController(discoveryService: Get.find<DiscoveryService>()),
    );
  }
}
