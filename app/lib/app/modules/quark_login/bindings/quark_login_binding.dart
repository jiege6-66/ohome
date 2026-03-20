import 'package:get/get.dart';

import '../controllers/quark_login_controller.dart';

class QuarkLoginBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<QuarkLoginController>(() => QuarkLoginController());
  }
}
