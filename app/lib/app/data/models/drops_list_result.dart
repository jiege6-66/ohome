import 'drops_event_model.dart';
import 'drops_item_model.dart';

class DropsItemsListResult {
  const DropsItemsListResult({required this.records, required this.total});

  final List<DropsItemModel> records;
  final int total;

  factory DropsItemsListResult.fromJson(Map<String, dynamic> json) {
    final recordsData = json['records'];
    final records = recordsData is List
        ? recordsData
              .whereType<Map>()
              .map(
                (item) => DropsItemModel.fromJson(item.cast<String, dynamic>()),
              )
              .toList(growable: false)
        : const <DropsItemModel>[];
    return DropsItemsListResult(
      records: records,
      total: _toInt(json['total']) ?? records.length,
    );
  }
}

class DropsEventsListResult {
  const DropsEventsListResult({required this.records, required this.total});

  final List<DropsEventModel> records;
  final int total;

  factory DropsEventsListResult.fromJson(Map<String, dynamic> json) {
    final recordsData = json['records'];
    final records = recordsData is List
        ? recordsData
              .whereType<Map>()
              .map(
                (item) =>
                    DropsEventModel.fromJson(item.cast<String, dynamic>()),
              )
              .toList(growable: false)
        : const <DropsEventModel>[];
    return DropsEventsListResult(
      records: records,
      total: _toInt(json['total']) ?? records.length,
    );
  }
}

int? _toInt(dynamic value) {
  if (value is int) return value;
  if (value is String) return int.tryParse(value);
  return null;
}
