import 'dart:async';
import 'dart:convert';
import 'dart:math';

import 'package:get/get.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import '../data/models/app_message_push_event.dart';
import '../utils/app_env.dart';
import 'auth_service.dart';

class AppMessagePushService extends GetxService {
  AppMessagePushService({AuthService? authService})
    : _authService = authService ?? Get.find<AuthService>();

  final AuthService _authService;
  final _eventsController = StreamController<AppMessagePushEvent>.broadcast();

  final connected = false.obs;

  Stream<AppMessagePushEvent> get events => _eventsController.stream;

  Worker? _authWorker;
  WebSocketChannel? _channel;
  StreamSubscription<dynamic>? _channelSubscription;
  Timer? _reconnectTimer;
  int _session = 0;
  int _reconnectAttempt = 0;

  @override
  void onInit() {
    super.onInit();
    _authWorker = ever<String?>(_authService.accessToken, (_) {
      _handleTokenChanged();
    });
    _handleTokenChanged();
  }

  @override
  void onClose() {
    _authWorker?.dispose();
    _authWorker = null;
    _cancelReconnect();
    unawaited(_closeActiveChannel());
    unawaited(_eventsController.close());
    super.onClose();
  }

  void _handleTokenChanged() {
    final token = _authService.accessToken.value?.trim() ?? '';
    if (token.isEmpty) {
      _disconnect();
      return;
    }
    unawaited(_connect());
  }

  Future<void> reconnect() async {
    if (!_authService.isLoggedIn) return;
    await _connect(force: true);
  }

  Future<void> _connect({bool force = false}) async {
    final token = _authService.accessToken.value?.trim() ?? '';
    if (token.isEmpty) {
      _disconnect();
      return;
    }

    if (!force && connected.value && _channel != null) {
      return;
    }

    final session = ++_session;
    _cancelReconnect();
    await _closeActiveChannel();

    try {
      final channel = WebSocketChannel.connect(_buildUri(token));
      _channel = channel;
      await channel.ready;
      if (session != _session) {
        await channel.sink.close();
        return;
      }

      connected.value = true;
      _reconnectAttempt = 0;
      _channelSubscription = channel.stream.listen(
        _handleMessage,
        onError: (Object error, StackTrace stackTrace) =>
            _handleDisconnect(session),
        onDone: () => _handleDisconnect(session),
        cancelOnError: true,
      );
    } catch (_) {
      if (session != _session) return;
      connected.value = false;
      await _closeActiveChannel();
      _scheduleReconnect();
    }
  }

  void _disconnect() {
    _session++;
    connected.value = false;
    _reconnectAttempt = 0;
    _cancelReconnect();
    unawaited(_closeActiveChannel());
  }

  void _handleDisconnect(int session) {
    if (session != _session) return;
    connected.value = false;
    unawaited(_closeActiveChannel());
    _scheduleReconnect();
  }

  void _scheduleReconnect() {
    if (!_authService.isLoggedIn || _reconnectTimer != null) {
      return;
    }

    final delay = Duration(seconds: min(30, 1 << min(_reconnectAttempt, 4)));
    _reconnectAttempt += 1;
    _reconnectTimer = Timer(delay, () {
      _reconnectTimer = null;
      unawaited(_connect());
    });
  }

  void _cancelReconnect() {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
  }

  Future<void> _closeActiveChannel() async {
    final subscription = _channelSubscription;
    _channelSubscription = null;
    await subscription?.cancel();

    final channel = _channel;
    _channel = null;
    await channel?.sink.close();
  }

  void _handleMessage(dynamic data) {
    try {
      final decoded = _decodePayload(data);
      if (decoded is! Map) return;
      final event = AppMessagePushEvent.fromJson(
        decoded.cast<String, dynamic>(),
      );
      if (event.event.isEmpty) return;
      _eventsController.add(event);
    } catch (_) {
      return;
    }
  }

  dynamic _decodePayload(dynamic data) {
    if (data is String) {
      return jsonDecode(data);
    }
    if (data is List<int>) {
      return jsonDecode(utf8.decode(data));
    }
    return null;
  }

  Uri _buildUri(String accessToken) {
    final baseUri = Uri.parse(AppEnv.instance.apiBaseUrl);
    final pathSegments =
        baseUri.pathSegments
            .where((segment) => segment.trim().isNotEmpty)
            .toList(growable: true)
          ..addAll(const <String>['appMessage', 'ws']);

    return baseUri.replace(
      scheme: baseUri.scheme == 'https' ? 'wss' : 'ws',
      pathSegments: pathSegments,
      queryParameters: <String, String>{'access_token': accessToken},
    );
  }
}
