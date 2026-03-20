import 'discovery_info.dart';

enum DiscoverySource {
  previousSuccess,
  mdns,
  subnetScan;

  String get label {
    switch (this) {
      case DiscoverySource.previousSuccess:
        return '上次成功';
      case DiscoverySource.mdns:
        return 'mDNS';
      case DiscoverySource.subnetScan:
        return '子网扫描';
    }
  }
}

class DiscoveredServer {
  DiscoveredServer({
    required this.info,
    required this.origin,
    Set<DiscoverySource>? sources,
  }) : sources = sources ?? <DiscoverySource>{};

  final DiscoveryInfo info;
  final String origin;
  final Set<DiscoverySource> sources;

  bool get isPreviousSuccess =>
      sources.contains(DiscoverySource.previousSuccess);

  bool get isFromMdns => sources.contains(DiscoverySource.mdns);

  bool get isFromSubnetScan => sources.contains(DiscoverySource.subnetScan);

  String get serviceName => info.serviceName;

  String get version => info.version;

  int get port => info.port;

  String get instanceId => info.instanceId;

  String get apiBaseUrl => info.apiBaseUrl;

  int get rank {
    var score = 0;
    if (isPreviousSuccess) score += 100;
    if (isFromMdns) score += 10;
    if (isFromSubnetScan) score += 1;
    return score;
  }

  List<String> get sourceLabels => DiscoverySource.values
      .where(sources.contains)
      .map((source) => source.label)
      .toList(growable: false);

  DiscoveredServer merge(DiscoveredServer other) {
    final preferred = other.rank > rank ? other : this;
    return DiscoveredServer(
      info: preferred.info,
      origin: preferred.origin,
      sources: <DiscoverySource>{...sources, ...other.sources},
    );
  }
}
