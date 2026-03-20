import 'package:get/get.dart';

import '../controllers/music_player_controller.dart';

class MusicPlayerBinding extends Bindings {
  @override
  void dependencies() {
    if (!Get.isRegistered<MusicPlayerController>()) {
      Get.put<MusicPlayerController>(
        MusicPlayerController(defaultApplicationType: 'music'),
        permanent: true,
      );
    }
  }
}
