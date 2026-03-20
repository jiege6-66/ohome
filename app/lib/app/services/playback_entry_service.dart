import 'package:flutter/foundation.dart';
import 'package:get/get.dart';

import '../data/models/media_history_entry.dart';
import '../data/models/quark_file_entry.dart';
import '../utils/media_path.dart';

typedef PlaybackStreamUriBuilder = Uri? Function(WebdavFileEntry file);

class PlaybackEntryConfig {
  const PlaybackEntryConfig({
    required this.applicationType,
    required this.playerRoute,
    required this.streamUriBuilder,
    required this.supportedExtensions,
    required this.defaultTitle,
    required this.isVideo,
  });

  factory PlaybackEntryConfig.forApplication(
    String applicationType, {
    PlaybackStreamUriBuilder? streamUriBuilder,
    List<String>? supportedExtensions,
  }) {
    final normalized = _normalizeApplicationType(applicationType);
    final isVideo = normalized == 'tv' || normalized == 'playlet';
    return PlaybackEntryConfig(
      applicationType: normalized,
      playerRoute: isVideo
          ? (normalized == 'playlet' ? '/playlet-player' : '/player')
          : '/music-player',
      streamUriBuilder:
          streamUriBuilder ??
          (file) => _defaultStreamUriBuilder(normalized, file),
      supportedExtensions:
          supportedExtensions ??
          (isVideo
              ? _defaultVideoExtensions
              : _defaultAudioExtensionsFor(normalized)),
      defaultTitle: switch (normalized) {
        'playlet' => '短剧',
        'music' => '音乐',
        'xiaoshuo' => '有声书',
        _ => '影视',
      },
      isVideo: isVideo,
    );
  }

  final String applicationType;
  final String playerRoute;
  final PlaybackStreamUriBuilder streamUriBuilder;
  final List<String> supportedExtensions;
  final String defaultTitle;
  final bool isVideo;

  String get resumePathArgumentKey =>
      isVideo ? 'resumeEpisodePath' : 'resumeTrackPath';

  bool isPlayableName(String name) {
    final lower = name.trim().toLowerCase();
    return supportedExtensions.any(lower.endsWith);
  }

  bool isPlayableEntry(WebdavFileEntry entry) {
    return !entry.isDir && isPlayableName(entry.name);
  }

  static const String loadCollectionOnEnterKey = 'loadCollectionOnEnter';
  static const String preferredItemPathKey = 'preferredItemPath';

  static const List<String> _defaultVideoExtensions = <String>[
    '.mp4',
    '.mkv',
    '.mov',
    '.m4v',
    '.avi',
    '.wmv',
    '.flv',
    '.webm',
  ];

  static List<String> _defaultAudioExtensionsFor(String applicationType) {
    final normalized = applicationType.trim().toLowerCase();
    if (normalized == 'xiaoshuo') {
      return const <String>[
        '.mp3',
        '.aac',
        '.m4a',
        '.m4b',
        '.flac',
        '.wav',
        '.ogg',
        '.opus',
        if (!kIsWeb) '.wma',
      ];
    }
    return const <String>[
      '.mp3',
      '.aac',
      '.m4a',
      '.flac',
      '.wav',
      '.ogg',
      '.opus',
      '.m4b',
      '.mpga',
      if (!kIsWeb) '.wma',
    ];
  }

  static String _normalizeApplicationType(String value) {
    final normalized = value.trim().toLowerCase();
    switch (normalized) {
      case 'tv':
      case 'playlet':
      case 'music':
      case 'xiaoshuo':
        return normalized;
      default:
        return 'tv';
    }
  }

  static Uri? _defaultStreamUriBuilder(
    String applicationType,
    WebdavFileEntry entry,
  ) {
    return entry.resolveStreamUri(applicationType: applicationType);
  }
}

class PlaybackLaunchPayload {
  const PlaybackLaunchPayload({required this.route, required this.arguments});

  final String route;
  final Map<String, dynamic> arguments;
}

class PlaybackEntryService extends GetxService {
  PlaybackEntryService();

  Future<PlaybackLaunchPayload?> buildFromHistoryEntry(
    MediaHistoryEntry entry,
  ) async {
    final config = PlaybackEntryConfig.forApplication(entry.applicationType);
    final folderPath = entry.folderPath.trim();
    final itemPath = MediaPath.normalize(entry.itemPath);
    if (folderPath.isEmpty || itemPath.isEmpty) {
      return null;
    }

    final arguments = _buildLazyArguments(
      config: config,
      folderPath: folderPath,
      folderTitle: _resolveFolderTitle(entry.folderTitle, folderPath, config),
      preferredPath: itemPath,
      resumePath: itemPath,
      resumePositionMs: entry.positionMs > 0 ? entry.positionMs : null,
    );
    return PlaybackLaunchPayload(
      route: config.playerRoute,
      arguments: arguments,
    );
  }

  Future<PlaybackLaunchPayload?> buildCollectionFromFolder({
    required PlaybackEntryConfig config,
    required String folderPath,
    required String folderTitle,
    String? preferredPath,
  }) async {
    final safeFolder = folderPath.trim();
    if (safeFolder.isEmpty) return null;

    final arguments = _buildLazyArguments(
      config: config,
      folderPath: safeFolder,
      folderTitle: folderTitle,
      preferredPath: preferredPath,
    );
    return PlaybackLaunchPayload(
      route: config.playerRoute,
      arguments: arguments,
    );
  }

  String titleFromPath(String path, {String fallback = '播放'}) {
    final normalized = path.replaceAll('\\', '/').trim();
    if (normalized.isEmpty) return fallback;
    final parts = normalized
        .split('/')
        .where((part) => part.trim().isNotEmpty)
        .toList(growable: false);
    return parts.isEmpty ? fallback : parts.last;
  }

  Map<String, dynamic> _buildLazyArguments({
    required PlaybackEntryConfig config,
    required String folderPath,
    required String folderTitle,
    String? preferredPath,
    String? resumePath,
    int? resumePositionMs,
  }) {
    final arguments = <String, dynamic>{
      'title': _resolveFolderTitle(
        '',
        folderPath,
        config,
        explicit: folderTitle,
      ),
      'folderPath': folderPath,
      'applicationType': config.applicationType,
      PlaybackEntryConfig.loadCollectionOnEnterKey: true,
    };

    final normalizedPreferredPath = MediaPath.normalize(preferredPath);
    if (normalizedPreferredPath.isNotEmpty) {
      arguments[PlaybackEntryConfig.preferredItemPathKey] =
          normalizedPreferredPath;
    }

    final normalizedResumePath = MediaPath.normalize(resumePath);
    if (normalizedResumePath.isNotEmpty) {
      arguments[config.resumePathArgumentKey] = normalizedResumePath;
    }
    if (resumePositionMs != null && resumePositionMs > 0) {
      arguments['resumePositionMs'] = resumePositionMs;
    }
    if (!config.isVideo) {
      arguments['supportedExtensions'] = config.supportedExtensions;
    }
    return arguments;
  }

  String _resolveFolderTitle(
    String rawFolderTitle,
    String folderPath,
    PlaybackEntryConfig config, {
    String? explicit,
  }) {
    final direct = explicit?.trim() ?? '';
    if (direct.isNotEmpty) return direct;
    final raw = rawFolderTitle.trim();
    if (raw.isNotEmpty) return raw;
    return titleFromPath(folderPath, fallback: config.defaultTitle);
  }
}
