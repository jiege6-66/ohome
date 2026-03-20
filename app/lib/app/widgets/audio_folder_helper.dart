import '../services/playback_entry_service.dart';
import 'folder_playback_helper.dart';

typedef StreamUriBuilder = PlaybackStreamUriBuilder;

class AudioFolderHelper extends FolderPlaybackHelper {
  AudioFolderHelper({
    required String applicationType,
    required StreamUriBuilder streamUriBuilder,
    super.entryService,
    List<String>? supportedExtensions,
  }) : super(
         config: PlaybackEntryConfig.forApplication(
           applicationType,
           streamUriBuilder: streamUriBuilder,
           supportedExtensions: supportedExtensions,
         ),
       );
}
