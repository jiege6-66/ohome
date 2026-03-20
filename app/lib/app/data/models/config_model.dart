class ConfigModel {
  const ConfigModel({required this.raw});

  final Map<String, dynamic> raw;

  factory ConfigModel.fromJson(Map<String, dynamic> json) {
    return ConfigModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id {
    final value = raw['id'];
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }

  String get name => (raw['name'] ?? '').toString().trim();

  String get key => (raw['key'] ?? '').toString().trim();

  String get value => (raw['value'] ?? '').toString();

  String get isLock => (raw['isLock'] ?? '').toString().trim();

  String get remark => (raw['remark'] ?? '').toString().trim();

  DateTime? get updatedAt {
    final value = raw['updatedAt'];
    if (value is String) {
      final text = value.trim();
      if (text.isEmpty) return null;
      return DateTime.tryParse(text)?.toLocal();
    }
    return null;
  }
}
