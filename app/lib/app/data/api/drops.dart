import 'package:dio/dio.dart';
import 'package:image_picker/image_picker.dart';

import '../../utils/http_client.dart';
import '../models/drops_event_model.dart';
import '../models/drops_item_model.dart';
import '../models/drops_list_result.dart';
import '../models/drops_overview_model.dart';

class DropsApi {
  DropsApi({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<DropsOverviewModel> getOverview() {
    return _httpClient.get<DropsOverviewModel>(
      '/dropsOverview',
      decoder: (data) => DropsOverviewModel.fromJson(_asMap(data)),
    );
  }

  Future<List<String>> getLocationSuggestions({String? keyword}) {
    return _httpClient.get<List<String>>(
      '/dropsLocation/suggestions',
      queryParameters: keyword == null || keyword.trim().isEmpty
          ? null
          : <String, dynamic>{'keyword': keyword.trim()},
      decoder: (data) {
        if (data is! List) return const <String>[];
        return data
            .map((item) => item.toString().trim())
            .toList(growable: false);
      },
    );
  }

  Future<DropsItemsListResult> getItemList({
    String? scopeType,
    String? category,
    String? keyword,
    int page = 1,
    int limit = 20,
  }) {
    return _httpClient.post<DropsItemsListResult>(
      '/dropsItem/list',
      data: <String, dynamic>{
        'page': page,
        'limit': limit,
        if (scopeType != null && scopeType.trim().isNotEmpty)
          'scopeType': scopeType.trim(),
        if (category != null && category.trim().isNotEmpty)
          'category': category.trim(),
        if (keyword != null && keyword.trim().isNotEmpty)
          'keyword': keyword.trim(),
      },
      decoder: (data) => DropsItemsListResult.fromJson(_asMap(data)),
    );
  }

  Future<DropsItemModel> getItemDetail(int id) {
    return _httpClient.get<DropsItemModel>(
      '/dropsItem/$id',
      decoder: (data) => DropsItemModel.fromJson(_asMap(data)),
    );
  }

  Future<DropsItemModel> createItem({
    required Map<String, dynamic> fields,
    required List<XFile> photos,
  }) async {
    if (photos.isEmpty) {
      throw ApiException('请至少拍摄一张物资照片');
    }
    final formData = FormData.fromMap({
      ...fields,
      'photos': await Future.wait(
        photos.map((file) async {
          final bytes = await file.readAsBytes();
          return MultipartFile.fromBytes(
            bytes,
            filename: file.name.trim().isEmpty ? 'photo.jpg' : file.name.trim(),
          );
        }),
      ),
    });
    return _httpClient.post<DropsItemModel>(
      '/dropsItem/create',
      data: formData,
      decoder: (data) => DropsItemModel.fromJson(_asMap(data)),
    );
  }

  Future<DropsItemModel> updateItem({
    required int id,
    required Map<String, dynamic> fields,
  }) {
    return _httpClient.put<DropsItemModel>(
      '/dropsItem/$id',
      data: fields,
      decoder: (data) => DropsItemModel.fromJson(_asMap(data)),
    );
  }

  Future<DropsItemModel> addItemPhotos({
    required int id,
    required List<XFile> photos,
  }) async {
    if (photos.isEmpty) {
      throw ApiException('请至少拍摄一张照片');
    }
    final formData = FormData.fromMap({
      'photos': await Future.wait(
        photos.map((file) async {
          final bytes = await file.readAsBytes();
          return MultipartFile.fromBytes(
            bytes,
            filename: file.name.trim().isEmpty ? 'photo.jpg' : file.name.trim(),
          );
        }),
      ),
    });
    return _httpClient.post<DropsItemModel>(
      '/dropsItem/$id/photos',
      data: formData,
      decoder: (data) => DropsItemModel.fromJson(_asMap(data)),
    );
  }

  Future<void> deleteItem(int id) {
    return _httpClient.delete<void>('/dropsItem/$id', decoder: (_) {});
  }

  Future<void> deleteItemPhoto({required int itemId, required int photoId}) {
    return _httpClient.delete<void>(
      '/dropsItem/$itemId/photos/$photoId',
      decoder: (_) {},
    );
  }

  Future<DropsEventsListResult> getEventList({
    String? scopeType,
    String? eventType,
    int? month,
    String? keyword,
    int page = 1,
    int limit = 20,
  }) {
    return _httpClient.post<DropsEventsListResult>(
      '/dropsEvent/list',
      data: <String, dynamic>{
        'page': page,
        'limit': limit,
        if (scopeType != null && scopeType.trim().isNotEmpty)
          'scopeType': scopeType.trim(),
        if (eventType != null && eventType.trim().isNotEmpty)
          'eventType': eventType.trim(),
        if (month != null && month > 0) 'month': month,
        if (keyword != null && keyword.trim().isNotEmpty)
          'keyword': keyword.trim(),
      },
      decoder: (data) => DropsEventsListResult.fromJson(_asMap(data)),
    );
  }

  Future<DropsEventModel> getEventDetail(int id) {
    return _httpClient.get<DropsEventModel>(
      '/dropsEvent/$id',
      decoder: (data) => DropsEventModel.fromJson(_asMap(data)),
    );
  }

  Future<DropsEventModel> saveEvent(Map<String, dynamic> payload) {
    final id = payload['id'];
    if (id != null) {
      return _httpClient.put<DropsEventModel>(
        '/dropsEvent/$id',
        data: payload,
        decoder: (data) => DropsEventModel.fromJson(_asMap(data)),
      );
    } else {
      return _httpClient.post<DropsEventModel>(
        '/dropsEvent/add',
        data: payload,
        decoder: (data) => DropsEventModel.fromJson(_asMap(data)),
      );
    }
  }

  Future<void> deleteEvent(int id) {
    return _httpClient.delete<void>('/dropsEvent/$id', decoder: (_) {});
  }

  static Map<String, dynamic> _asMap(dynamic data) {
    if (data is Map<String, dynamic>) return data;
    if (data is Map) return data.cast<String, dynamic>();
    throw ApiException('响应格式错误');
  }
}
