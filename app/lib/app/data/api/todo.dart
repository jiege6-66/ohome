import '../../utils/http_client.dart';
import '../models/todo_item_model.dart';
import '../models/todo_list_result.dart';

class TodoApi {
  TodoApi({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<TodoListResult> getTodoList({
    int page = 1,
    int limit = 200,
    bool showErrorToast = true,
  }) {
    return _httpClient.post<TodoListResult>(
      '/todo/list',
      data: <String, dynamic>{'page': page, 'limit': limit},
      showErrorToast: showErrorToast,
      decoder: (data) => TodoListResult.fromJson(_asMap(data)),
    );
  }

  Future<TodoItemModel> addTodoItem({required String title}) {
    return _httpClient.post<TodoItemModel>(
      '/todo/add',
      data: <String, dynamic>{'title': title.trim()},
      decoder: (data) => TodoItemModel.fromJson(_asMap(data)),
    );
  }

  Future<TodoItemModel> updateTodoItem({
    required int id,
    required String title,
  }) {
    return _httpClient.put<TodoItemModel>(
      '/todo/$id',
      data: <String, dynamic>{'title': title.trim()},
      decoder: (data) => TodoItemModel.fromJson(_asMap(data)),
    );
  }

  Future<TodoItemModel> updateTodoStatus({
    required int id,
    required bool completed,
  }) {
    return _httpClient.put<TodoItemModel>(
      '/todo/$id/status',
      data: <String, dynamic>{'completed': completed},
      decoder: (data) => TodoItemModel.fromJson(_asMap(data)),
    );
  }

  Future<void> reorderTodoItems({required List<int> ids}) {
    return _httpClient.put<void>(
      '/todo/reorder',
      data: <String, dynamic>{'ids': ids},
      decoder: (_) {},
    );
  }

  Future<void> deleteTodoItem(int id) {
    return _httpClient.delete<void>('/todo/$id', decoder: (_) {});
  }

  static Map<String, dynamic> _asMap(dynamic data) {
    if (data is Map<String, dynamic>) return data;
    if (data is Map) return data.cast<String, dynamic>();
    throw ApiException('待办响应格式错误');
  }
}
