class DictDataModel {
  const DictDataModel({required this.raw});

  final Map<String, dynamic> raw;

  factory DictDataModel.fromJson(Map<String, dynamic> json) {
    return DictDataModel(raw: Map<String, dynamic>.from(json));
  }

  int get sort => _toInt(raw['sort']) ?? 0;

  String get label => (raw['label'] ?? '').toString().trim();

  String get value => (raw['value'] ?? '').toString().trim();

  String get dictType => (raw['dictType'] ?? '').toString().trim();

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }
}
