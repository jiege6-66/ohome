import '../../utils/backend_url_resolver.dart';

class DropsItemPhotoModel {
  const DropsItemPhotoModel({required this.raw});

  final Map<String, dynamic> raw;

  factory DropsItemPhotoModel.fromJson(Map<String, dynamic> json) {
    return DropsItemPhotoModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id => _toInt(raw['id']);

  String get fileName => (raw['fileName'] ?? '').toString().trim();

  String get filePath => (raw['filePath'] ?? '').toString().trim();

  String get url => (raw['url'] ?? '').toString().trim();

  String get resolvedUrl => BackendUrlResolver.resolve(url);

  bool get isCover {
    final value = raw['isCover'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() == 'true';
  }

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }
}
