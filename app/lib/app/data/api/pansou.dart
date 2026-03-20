import '../../utils/http_client.dart';
import '../models/pansou_resource_item.dart';

class PansouRepository {
  PansouRepository({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<List<PansouResourceItem>> searchQuark(
    String keyword, {
    bool refresh = false,
  }) {
    final kw = keyword.trim();
    if (kw.isEmpty) return Future.value(const <PansouResourceItem>[]);

    return _httpClient.post<List<PansouResourceItem>>(
      'pansou/search',
      data: <String, dynamic>{
        'kw': kw,
        'res': 'merge',
        'cloud_types': <String>['quark'],
        if (refresh) 'refresh': true,
      },
      decoder: (data) {
        final list = _extractMergedList(data, type: 'quark');
        return list
            .map((e) => PansouResourceItem.fromJson(e, type: 'quark'))
            .where((e) => !e.isEmpty)
            .toList(growable: false);
      },
    );
  }

  static List<Map<String, dynamic>> _extractMergedList(
    dynamic data, {
    required String type,
  }) {
    if (data is! Map) return const <Map<String, dynamic>>[];

    final mergedByType = data['merged_by_type'];
    final dynamic list = mergedByType is Map ? mergedByType[type] : data[type];
    if (list is! List) return const <Map<String, dynamic>>[];

    return list
        .whereType<Map>()
        .map((e) => e.map((k, v) => MapEntry(k.toString(), v)))
        .toList(growable: false);
  }
}
