import 'package:shared_preferences/shared_preferences.dart';

class RememberedServer {
  const RememberedServer({
    required this.origin,
    required this.instanceId,
    required this.port,
  });

  final String origin;
  final String instanceId;
  final int port;
}

class DiscoveryStorage {
  static const _kLastOrigin = 'discoveryLastOrigin';
  static const _kLastInstanceId = 'discoveryLastInstanceId';
  static const _kLastPort = 'discoveryLastPort';

  SharedPreferences? _prefs;

  Future<SharedPreferences> _instance() async {
    final cached = _prefs;
    if (cached != null) return cached;
    final prefs = await SharedPreferences.getInstance();
    _prefs = prefs;
    return prefs;
  }

  Future<RememberedServer?> readLastSuccessfulServer() async {
    final prefs = await _instance();
    final origin = prefs.getString(_kLastOrigin)?.trim() ?? '';
    final instanceId = prefs.getString(_kLastInstanceId)?.trim() ?? '';
    final port = prefs.getInt(_kLastPort) ?? 0;
    if (origin.isEmpty || instanceId.isEmpty || port <= 0) {
      return null;
    }
    return RememberedServer(origin: origin, instanceId: instanceId, port: port);
  }

  Future<void> writeLastSuccessfulServer(RememberedServer server) async {
    final prefs = await _instance();
    await prefs.setString(_kLastOrigin, server.origin.trim());
    await prefs.setString(_kLastInstanceId, server.instanceId.trim());
    await prefs.setInt(_kLastPort, server.port);
  }
}
