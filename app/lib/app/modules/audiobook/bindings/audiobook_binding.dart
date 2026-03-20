import 'package:get/get.dart';

import '../controllers/audiobook_controller.dart';

class AudiobookBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<AudiobookController>(() => AudiobookController());
  }
}
