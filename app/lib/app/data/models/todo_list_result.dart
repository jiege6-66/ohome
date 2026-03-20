import 'todo_item_model.dart';

class TodoListResult {
  const TodoListResult({required this.records, required this.total});

  final List<TodoItemModel> records;
  final int total;

  factory TodoListResult.fromJson(Map<String, dynamic> json) {
    final records = <TodoItemModel>[];
    final rawRecords = json['records'];
    if (rawRecords is List) {
      for (final item in rawRecords) {
        if (item is Map<String, dynamic>) {
          records.add(TodoItemModel.fromJson(item));
        } else if (item is Map) {
          records.add(TodoItemModel.fromJson(item.cast<String, dynamic>()));
        }
      }
    }

    final rawTotal = json['total'];
    final total = rawTotal is int
        ? rawTotal
        : int.tryParse(rawTotal?.toString() ?? '') ?? records.length;

    return TodoListResult(records: records, total: total);
  }
}
