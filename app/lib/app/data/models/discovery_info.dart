class DiscoveryInfo {
  DiscoveryInfo({
    required this.instanceId,
    required this.serviceName,
    required this.version,
    required this.apiBaseUrl,
    required this.port,
    required this.capabilities,
  });

  final String instanceId;
  final String serviceName;
  final String version;
  final String apiBaseUrl;
  final int port;
  final List<String> capabilities;

  factory DiscoveryInfo.fromJson(Map<String, dynamic> json) {
    return DiscoveryInfo(
      instanceId: (json['instanceId'] as String? ?? '').trim(),
      serviceName: (json['serviceName'] as String? ?? '').trim(),
      version: (json['version'] as String? ?? '').trim(),
      apiBaseUrl: (json['apiBaseUrl'] as String? ?? '').trim(),
      port: _parseInt(json['port']),
      capabilities: _parseCapabilities(json['capabilities']),
    );
  }

  String get origin {
    final uri = Uri.parse(apiBaseUrl);
    final portText = uri.hasPort ? ':${uri.port}' : '';
    return '${uri.scheme}://${uri.host}$portText';
  }

  static int _parseInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value.trim()) ?? 0;
    return 0;
  }

  static List<String> _parseCapabilities(dynamic value) {
    if (value is! List) return const <String>[];
    return value
        .map((item) => item?.toString().trim() ?? '')
        .where((item) => item.isNotEmpty)
        .toList(growable: false);
  }
}
