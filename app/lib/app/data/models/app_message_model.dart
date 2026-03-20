class AppMessageModel {
  const AppMessageModel({required this.raw});

  final Map<String, dynamic> raw;

  factory AppMessageModel.fromJson(Map<String, dynamic> json) {
    return AppMessageModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id => _toInt(raw['id']);

  String get source => (raw['source'] ?? '').toString().trim();

  String get sourceKey => (raw['sourceKey'] ?? '').toString().trim();

  String get messageType => (raw['messageType'] ?? '').toString().trim();

  String get bizType => (raw['bizType'] ?? '').toString().trim();

  int? get bizId => _toInt(raw['bizId']);

  String get title => (raw['title'] ?? '').toString().trim();

  String get summary => (raw['summary'] ?? '').toString().trim();

  DateTime? get triggerDate => _parseDateTime(raw['triggerDate']);

  bool get read {
    final value = raw['read'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() == 'true';
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
