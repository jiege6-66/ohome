class QuarkAutoSaveTaskModel {
  const QuarkAutoSaveTaskModel({required this.raw});

  final Map<String, dynamic> raw;

  factory QuarkAutoSaveTaskModel.fromJson(Map<String, dynamic> json) {
    return QuarkAutoSaveTaskModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id {
    final value = raw['id'];
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }

  String get taskName => (raw['taskName'] ?? '').toString().trim();

  String get shareUrl => (raw['shareUrl'] ?? '').toString().trim();

  String get savePath => (raw['savePath'] ?? '').toString().trim();

  String get scheduleType => (raw['scheduleType'] ?? '').toString().trim();

  String get runTime => (raw['runTime'] ?? '').toString().trim();

  String get runWeek => (raw['runWeek'] ?? '').toString().trim();

  bool get enabled {
    final value = raw['enabled'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    final text = value?.toString().trim().toLowerCase() ?? '';
    return text == '1' || text == 'true';
  }

  DateTime? get updatedAt => _parseDateTime(raw['updatedAt']);

  DateTime? get lastRunAt => _parseDateTime(raw['lastRunAt']);

  List<int> get runWeekDays {
    return runWeek
        .split(',')
        .map((value) => int.tryParse(value.trim()))
        .whereType<int>()
        .where((value) => value >= 1 && value <= 7)
        .toList(growable: false);
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
