import 'quark_transfer_task_model.dart';

class QuarkTransferTaskListResult {
  const QuarkTransferTaskListResult({
    required this.records,
    required this.total,
  });

  final List<QuarkTransferTaskModel> records;
  final int total;

  factory QuarkTransferTaskListResult.fromJson(Map<String, dynamic> json) {
    final records = <QuarkTransferTaskModel>[];
    final rawRecords = json['records'];
    if (rawRecords is List) {
      for (final item in rawRecords) {
        if (item is Map<String, dynamic>) {
          records.add(QuarkTransferTaskModel.fromJson(item));
        } else if (item is Map) {
          records.add(
            QuarkTransferTaskModel.fromJson(item.cast<String, dynamic>()),
          );
        }
      }
    }

    final rawTotal = json['total'];
    final total = rawTotal is int
        ? rawTotal
        : int.tryParse(rawTotal?.toString() ?? '') ?? records.length;

    return QuarkTransferTaskListResult(records: records, total: total);
  }
}
