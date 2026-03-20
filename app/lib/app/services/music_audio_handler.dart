import 'dart:async';

import 'package:audio_service/audio_service.dart';

class MusicAudioHandler extends BaseAudioHandler with QueueHandler, SeekHandler {
  MusicAudioHandler({
    required Future<void> Function() onPlay,
    required Future<void> Function() onPause,
    required Future<void> Function() onStop,
    required Future<void> Function() onSkipToNext,
    required Future<void> Function() onSkipToPrevious,
    required Future<void> Function(Duration position) onSeek,
    required Future<void> Function(int index) onSkipToQueueItem,
  }) : _onPlay = onPlay,
       _onPause = onPause,
       _onStop = onStop,
       _onSkipToNext = onSkipToNext,
       _onSkipToPrevious = onSkipToPrevious,
       _onSeek = onSeek,
       _onSkipToQueueItem = onSkipToQueueItem;

  final Future<void> Function() _onPlay;
  final Future<void> Function() _onPause;
  final Future<void> Function() _onStop;
  final Future<void> Function() _onSkipToNext;
  final Future<void> Function() _onSkipToPrevious;
  final Future<void> Function(Duration position) _onSeek;
  final Future<void> Function(int index) _onSkipToQueueItem;

  void setQueueItems(List<MediaItem> items) {
    queue.add(List<MediaItem>.unmodifiable(items));
  }

  void setMediaItemData(MediaItem? item) {
    if (item == null) return;
    mediaItem.add(item);
  }

  void clearMediaItem() {
    mediaItem.add(null);
  }

  void setPlaybackStateData({
    required bool playing,
    required AudioProcessingState processingState,
    required Duration position,
    required Duration bufferedPosition,
    required double speed,
    required int? queueIndex,
  }) {
    final controls = <MediaControl>[
      MediaControl.skipToPrevious,
      if (playing) MediaControl.pause else MediaControl.play,
      MediaControl.skipToNext,
      MediaControl.stop,
    ];
    playbackState.add(
      PlaybackState(
        controls: controls,
        systemActions: const {
          MediaAction.seek,
          MediaAction.seekForward,
          MediaAction.seekBackward,
          MediaAction.skipToNext,
          MediaAction.skipToPrevious,
          MediaAction.playPause,
          MediaAction.stop,
        },
        androidCompactActionIndices: const [0, 1, 2],
        processingState: processingState,
        playing: playing,
        updatePosition: position,
        bufferedPosition: bufferedPosition,
        speed: speed,
        queueIndex: queueIndex,
      ),
    );
  }

  @override
  Future<void> play() => _onPlay();

  @override
  Future<void> pause() => _onPause();

  @override
  Future<void> stop() => _onStop();

  @override
  Future<void> seek(Duration position) => _onSeek(position);

  @override
  Future<void> skipToNext() => _onSkipToNext();

  @override
  Future<void> skipToPrevious() => _onSkipToPrevious();

  @override
  Future<void> skipToQueueItem(int index) => _onSkipToQueueItem(index);
}
