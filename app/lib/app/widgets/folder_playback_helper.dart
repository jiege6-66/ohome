import 'package:get/get.dart';

import '../data/models/quark_file_entry.dart';
import '../services/playback_entry_service.dart';

class FolderPlaybackHelper {
  FolderPlaybackHelper({
    required PlaybackEntryConfig config,
    PlaybackEntryService? entryService,
  }) : _config = config,
       _entryService = entryService ?? Get.find<PlaybackEntryService>();

  final PlaybackEntryConfig _config;
  final PlaybackEntryService _entryService;

  bool isPlayableName(String name) => _config.isPlayableName(name);

  bool isPlayableEntry(WebdavFileEntry entry) => _config.isPlayableEntry(entry);

  bool shouldAutoOpen(List<WebdavFileEntry> entries) {
    if (entries.isEmpty) return false;
    final hasPlayable = entries.any(isPlayableEntry);
    final hasFolder = entries.any((entry) => entry.isDir);
    return hasPlayable && !hasFolder;
  }

  String titleFromPath(String path, {String? fallback}) {
    return _entryService.titleFromPath(
      path,
      fallback: fallback ?? _config.defaultTitle,
    );
  }

  Future<bool> openPlayer({
    required String folderPath,
    required String folderTitle,
    String? preferredPath,
  }) async {
    final launch = await _entryService.buildCollectionFromFolder(
      config: _config,
      folderPath: folderPath,
      folderTitle: folderTitle,
      preferredPath: preferredPath,
    );
    if (launch == null) return false;
    await Get.toNamed(launch.route, arguments: launch.arguments);
    return true;
  }
}
