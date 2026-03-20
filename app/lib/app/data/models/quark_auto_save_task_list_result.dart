import 'quark_auto_save_task_model.dart';

class QuarkAutoSaveTaskListResult {
  const QuarkAutoSaveTaskListResult({
    required this.records,
    required this.total,
  });

  final List<QuarkAutoSaveTaskModel> records;
  final int total;

  factory QuarkAutoSaveTaskListResult.fromJson(Map<String, dynamic> json) {
    final records = <QuarkAutoSaveTaskModel>[];
    final rawRecords = json['records'];
    if (rawRecords is List) {
      for (final item in rawRecords) {
        if (item is Map<String, dynamic>) {
          records.add(QuarkAutoSaveTaskModel.fromJson(item));
        } else if (item is Map) {
          records.add(QuarkAutoSaveTaskModel.fromJson(item.cast<String, dynamic>()));
        }
      }
    }

    final rawTotal = json['total'];
    final total = rawTotal is int
        ? rawTotal
        : int.tryParse(rawTotal?.toString() ?? '') ?? records.length;

    return QuarkAutoSaveTaskListResult(records: records, total: total);
  }
}
