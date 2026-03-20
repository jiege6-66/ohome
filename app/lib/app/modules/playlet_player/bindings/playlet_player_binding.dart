import 'package:get/get.dart';

import '../../player/controllers/player_controller.dart';

class PlayletPlayerBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<PlayerController>(
      () => PlayerController(applicationType: 'playlet'),
    );
  }
}
