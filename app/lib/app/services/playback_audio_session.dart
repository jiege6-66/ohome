import 'dart:async';

import 'package:audio_session/audio_session.dart';

enum PlaybackAudioProfile { media, speech }

class PlaybackAudioSession {
  PlaybackAudioSession({
    required Future<void> Function() onPauseRequested,
    required bool Function() isPlaying,
  }) : _onPauseRequested = onPauseRequested,
       _isPlaying = isPlaying;

  final Future<void> Function() _onPauseRequested;
  final bool Function() _isPlaying;

  AudioSession? _session;
  Future<void>? _initialization;
  StreamSubscription<AudioInterruptionEvent>? _interruptionSubscription;
  StreamSubscription<void>? _becomingNoisySubscription;

  Future<void> initialize() {
    return _initialization ??= _initializeInternal();
  }

  Future<bool> activate(PlaybackAudioProfile profile) async {
    final session = await _ensureSession();
    final configuration = _configurationFor(profile);
    await session.configure(configuration);
    return session.setActive(true, fallbackConfiguration: configuration);
  }

  Future<void> deactivate() async {
    final session = _session;
    if (session == null) return;
    try {
      await session.setActive(false);
    } catch (_) {}
  }

  Future<void> dispose() async {
    await _interruptionSubscription?.cancel();
    await _becomingNoisySubscription?.cancel();
    await deactivate();
  }

  Future<void> _initializeInternal() async {
    final session = await AudioSession.instance;
    _session = session;
    _interruptionSubscription = session.interruptionEventStream.listen(
      _handleInterruption,
    );
    _becomingNoisySubscription = session.becomingNoisyEventStream.listen(
      _handleBecomingNoisy,
    );
  }

  Future<AudioSession> _ensureSession() async {
    await initialize();
    return _session!;
  }

  AudioSessionConfiguration _configurationFor(PlaybackAudioProfile profile) {
    switch (profile) {
      case PlaybackAudioProfile.speech:
        return const AudioSessionConfiguration.speech().copyWith(
          androidWillPauseWhenDucked: false,
        );
      case PlaybackAudioProfile.media:
        return const AudioSessionConfiguration.music().copyWith(
          androidWillPauseWhenDucked: false,
        );
    }
  }

  void _handleInterruption(AudioInterruptionEvent event) {
    if (!event.begin || !_isPlaying()) return;
    if (event.type == AudioInterruptionType.duck) return;
    unawaited(_onPauseRequested());
  }

  void _handleBecomingNoisy(void _) {
    if (!_isPlaying()) return;
    unawaited(_onPauseRequested());
  }
}
