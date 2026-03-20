class TodoItemModel {
  const TodoItemModel({required this.raw});

  final Map<String, dynamic> raw;

  factory TodoItemModel.fromJson(Map<String, dynamic> json) {
    return TodoItemModel(raw: Map<String, dynamic>.from(json));
  }

  TodoItemModel copyWithRaw(Map<String, dynamic> patch) {
    return TodoItemModel(raw: <String, dynamic>{...raw, ...patch});
  }

  int? get id => _toInt(raw['id']);

  String get title => (raw['title'] ?? '').toString().trim();

  int get sortOrder => _toInt(raw['sortOrder']) ?? 0;

  bool get completed {
    final value = raw['completed'];
    if (value is bool) return value;
    if (value is num) return value != 0;
    final text = value?.toString().trim().toLowerCase() ?? '';
    return text == 'true' || text == '1';
  }

  int? get ownerUserId => _toInt(raw['ownerUserId']);

  int? get createdBy => _toInt(raw['createdBy']);

  int? get updatedBy => _toInt(raw['updatedBy']);

  DateTime? get createdAt => _parseDateTime(raw['createdAt']);

  DateTime? get updatedAt => _parseDateTime(raw['updatedAt']);

  DateTime? get completedAt => _parseDateTime(raw['completedAt']);

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.toInt();
    return int.tryParse(value?.toString() ?? '');
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
