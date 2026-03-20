import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/models/quark_file_entry.dart';
import '../../../widgets/audio_folder_helper.dart';
import '../../../widgets/resource_card_page.dart';
import '../controllers/audiobook_controller.dart';

class AudiobookView extends GetView<AudiobookController> {
  AudiobookView({super.key});

  late final AudioFolderHelper _audioHelper = AudioFolderHelper(
    applicationType: 'xiaoshuo',
    streamUriBuilder: controller.buildStreamUri,
    supportedExtensions: const <String>[
      '.mp3',
      '.aac',
      '.m4a',
      '.m4b',
      '.flac',
      '.wav',
      '.ogg',
      '.opus',
      '.wma',
    ],
  );

  @override
  Widget build(BuildContext context) {
    return ResourceCardPage(
      title: '有声书',
      controller: controller,
      iconBuilder: (entry) =>
          entry.isDir ? Icons.folder_rounded : Icons.headphones,
      iconColorBuilder: (entry) =>
          entry.isDir ? Colors.amber : Colors.lightBlueAccent,
      statusBuilder: (entry) => entry.isDir ? '文件夹' : '点击播放',
      onFolderTap: (entry, _, _) => _onFolderTap(entry),
      onFileTap: (entry, _, currentPath) => _onFileTap(entry, currentPath),
      shouldAutoOpen: (entries, _) => _audioHelper.shouldAutoOpen(entries),
      onAutoOpen: (_, currentPath) => _onAutoOpen(currentPath),
    );
  }

  Future<void> _onFolderTap(WebdavFileEntry entry) async {
    if (!entry.isDir) return;
    await controller.enterDir(entry);
  }

  Future<void> _onFileTap(WebdavFileEntry entry, String currentPath) async {
    if (!_audioHelper.isPlayableEntry(entry)) {
      Get.snackbar('提示', '不支持的文件类型');
      return;
    }
    await _openPlayerIfPossible(
      folderPath: currentPath,
      preferredPath: entry.path,
    );
  }

  Future<void> _onAutoOpen(String currentPath) async {
    await _openPlayerIfPossible(
      folderPath: currentPath,
      preferredPath: null,
      popAfter: controller.canGoBack,
    );
  }

  Future<void> _openPlayerIfPossible({
    required String folderPath,
    String? preferredPath,
    bool popAfter = false,
  }) async {
    final opened = await _audioHelper.openPlayer(
      folderPath: folderPath,
      folderTitle: _audioHelper.titleFromPath(folderPath),
      preferredPath: preferredPath,
    );
    if (!opened) return;
    if (popAfter && controller.canGoBack) {
      await controller.popDir();
    }
  }
}
