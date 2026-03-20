import 'app_message_model.dart';

class AppMessageListResult {
  const AppMessageListResult({
    required this.records,
    required this.total,
    required this.unreadCount,
  });

  final List<AppMessageModel> records;
  final int total;
  final int unreadCount;

  factory AppMessageListResult.fromJson(Map<String, dynamic> json) {
    final recordsData = json['records'];
    final records = recordsData is List
        ? recordsData
              .whereType<Map>()
              .map(
                (item) =>
                    AppMessageModel.fromJson(item.cast<String, dynamic>()),
              )
              .toList(growable: false)
        : const <AppMessageModel>[];

    return AppMessageListResult(
      records: records,
      total: _toInt(json['total']) ?? records.length,
      unreadCount: _toInt(json['unreadCount']) ?? 0,
    );
  }
}

int? _toInt(dynamic value) {
  if (value is int) return value;
  if (value is String) return int.tryParse(value);
  return null;
}
