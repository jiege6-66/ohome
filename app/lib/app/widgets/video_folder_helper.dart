import '../services/playback_entry_service.dart';
import 'folder_playback_helper.dart';

typedef VideoStreamUriBuilder = PlaybackStreamUriBuilder;

class VideoFolderHelper extends FolderPlaybackHelper {
  VideoFolderHelper({
    required String applicationType,
    required String playerRoute,
    required VideoStreamUriBuilder streamUriBuilder,
    super.entryService,
  }) : super(
         config: PlaybackEntryConfig(
           applicationType: applicationType.trim().toLowerCase(),
           playerRoute: playerRoute,
           streamUriBuilder: streamUriBuilder,
           supportedExtensions: const <String>[
             '.mp4',
             '.mkv',
             '.mov',
             '.m4v',
             '.avi',
             '.wmv',
             '.flv',
             '.webm',
           ],
           defaultTitle: applicationType.trim().toLowerCase() == 'playlet'
               ? '短剧'
               : '影视',
           isVideo: true,
         ),
       );
}
