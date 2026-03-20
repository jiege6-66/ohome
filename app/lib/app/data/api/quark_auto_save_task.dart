import '../../utils/http_client.dart';
import '../models/quark_auto_save_task_list_result.dart';
import '../models/quark_auto_save_task_upsert_payload.dart';

class QuarkAutoSaveTaskApi {
  QuarkAutoSaveTaskApi({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<QuarkAutoSaveTaskListResult> getTaskList({
    String? taskName,
    int page = 1,
    int limit = 20,
  }) {
    final payload = <String, dynamic>{'page': page, 'limit': limit};
    if (taskName != null && taskName.trim().isNotEmpty) {
      payload['taskName'] = taskName.trim();
    }

    return _httpClient.post<QuarkAutoSaveTaskListResult>(
      '/quarkAutoSaveTask/list',
      data: payload,
      decoder: (data) {
        if (data is Map<String, dynamic>) {
          return QuarkAutoSaveTaskListResult.fromJson(data);
        }
        if (data is Map) {
          return QuarkAutoSaveTaskListResult.fromJson(data.cast<String, dynamic>());
        }
        throw ApiException('夸克同步任务列表响应格式错误');
      },
    );
  }

  Future<void> addOrUpdateTask(QuarkAutoSaveTaskUpsertPayload payload) {
    return _httpClient.put<void>(
      '/quarkAutoSaveTask/add',
      data: payload.toJson(),
      decoder: (_) {},
    );
  }

  Future<void> deleteTask(int id) {
    return _httpClient.delete<void>('/quarkAutoSaveTask/$id', decoder: (_) {});
  }

  Future<void> runOnce(int id) {
    return _httpClient.post<void>(
      '/quarkAutoSaveTask/run/$id',
      decoder: (_) {},
    );
  }
}
