class DropsEventModel {
  const DropsEventModel({required this.raw});

  final Map<String, dynamic> raw;

  factory DropsEventModel.fromJson(Map<String, dynamic> json) {
    return DropsEventModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id => _toInt(raw['id']);

  String get scopeType => (raw['scopeType'] ?? '').toString().trim();

  String get title => (raw['title'] ?? '').toString().trim();

  String get eventType => (raw['eventType'] ?? '').toString().trim();

  String get calendarType => (raw['calendarType'] ?? '').toString().trim();

  int get eventYear => _toInt(raw['eventYear']) ?? 0;

  int get eventMonth => _toInt(raw['eventMonth']) ?? 0;

  int get eventDay => _toInt(raw['eventDay']) ?? 0;

  bool get isLeapMonth {
    final value = raw['isLeapMonth'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() == 'true';
  }

  bool get repeatYearly {
    final value = raw['repeatYearly'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() != 'false';
  }

  String get remark => (raw['remark'] ?? '').toString().trim();

  bool get enabled {
    final value = raw['enabled'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    return value?.toString().trim().toLowerCase() != 'false';
  }

  DateTime? get nextOccurAt => _parseDateTime(raw['nextOccurAt']);

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
