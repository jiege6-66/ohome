import 'package:get/get.dart';

import '../../messages/controllers/messages_controller.dart';
import '../controllers/home_controller.dart';

class HomeBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<MessagesController>(() => MessagesController());
    Get.lazyPut<HomeController>(() => HomeController());
  }
}
