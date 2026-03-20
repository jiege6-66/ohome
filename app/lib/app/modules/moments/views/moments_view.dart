import 'package:flutter/material.dart';

import 'package:get/get.dart';

import '../../../widgets/under_development_page.dart';
import '../controllers/moments_controller.dart';

class MomentsView extends GetView<MomentsController> {
  const MomentsView({super.key});
  @override
  Widget build(BuildContext context) {
    return const SafeArea(
      child: UnderDevelopmentPage(title: '动态', icon: Icons.explore_rounded),
    );
  }
}
