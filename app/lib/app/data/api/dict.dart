import '../../utils/http_client.dart';
import '../models/dict_data_model.dart';

class DictApi {
  DictApi({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<List<DictDataModel>> getDataList({
    required String dictType,
    int page = 1,
    int limit = 100,
  }) {
    return _httpClient.post<List<DictDataModel>>(
      'public/dict/data_list',
      showErrorToast: false,
      data: <String, dynamic>{
        'dictType': dictType.trim(),
        'page': page,
        'limit': limit,
      },
      decoder: (data) {
        final payload = _asMap(data);
        final records = payload['records'];
        if (records is! List) return const <DictDataModel>[];
        return records
            .whereType<Map>()
            .map((item) => DictDataModel.fromJson(item.cast<String, dynamic>()))
            .toList(growable: false);
      },
    );
  }

  static Map<String, dynamic> _asMap(dynamic data) {
    if (data is Map<String, dynamic>) return data;
    if (data is Map) return data.cast<String, dynamic>();
    throw ApiException('响应格式错误');
  }
}
