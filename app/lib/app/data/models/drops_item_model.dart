import '../../utils/backend_url_resolver.dart';
import 'drops_item_photo_model.dart';

class DropsItemModel {
  const DropsItemModel({required this.raw});

  final Map<String, dynamic> raw;

  factory DropsItemModel.fromJson(Map<String, dynamic> json) {
    return DropsItemModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id => _toInt(raw['id']);

  String get scopeType => (raw['scopeType'] ?? '').toString().trim();

  String get category => (raw['category'] ?? '').toString().trim();

  String get name => (raw['name'] ?? '').toString().trim();

  String get location => (raw['location'] ?? '').toString().trim();

  DateTime? get expireAt => _parseDateTime(raw['expireAt']);

  String get remark => (raw['remark'] ?? '').toString().trim();

  bool get enabled {
    final value = raw['enabled'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() == 'true';
  }

  String get coverUrl =>
      BackendUrlResolver.resolve((raw['coverUrl'] ?? '').toString().trim());

  int get photoCount => _toInt(raw['photoCount']) ?? photos.length;

  List<DropsItemPhotoModel> get photos {
    final data = raw['photos'];
    if (data is! List) return const <DropsItemPhotoModel>[];
    return data
        .whereType<Map>()
        .map(
          (item) => DropsItemPhotoModel.fromJson(item.cast<String, dynamic>()),
        )
        .toList(growable: false);
  }

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }

  static DateTime? _parseDateTime(dynamic value) {
    if (value is String) {
      final text = value.trim();
      if (text.isEmpty) return null;
      return DateTime.tryParse(text)?.toLocal();
    }
    return null;
  }
}
