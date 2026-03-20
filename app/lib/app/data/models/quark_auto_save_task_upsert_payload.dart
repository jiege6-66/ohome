import 'quark_auto_save_task_model.dart';

class QuarkAutoSaveTaskUpsertPayload {
  const QuarkAutoSaveTaskUpsertPayload({
    this.id,
    required this.taskName,
    required this.shareUrl,
    required this.savePath,
    required this.scheduleType,
    required this.runTime,
    required this.runWeek,
    required this.enabled,
  });

  final int? id;
  final String taskName;
  final String shareUrl;
  final String savePath;
  final String scheduleType;
  final String runTime;
  final String runWeek;
  final bool enabled;

  factory QuarkAutoSaveTaskUpsertPayload.fromTask(
    QuarkAutoSaveTaskModel task,
  ) {
    return QuarkAutoSaveTaskUpsertPayload(
      id: task.id,
      taskName: task.taskName,
      shareUrl: task.shareUrl,
      savePath: task.savePath,
      scheduleType: task.scheduleType,
      runTime: task.runTime,
      runWeek: task.runWeek,
      enabled: task.enabled,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      if (id != null) 'id': id,
      'taskName': taskName.trim(),
      'shareUrl': shareUrl.trim(),
      'savePath': savePath.trim(),
      'scheduleType': scheduleType.trim(),
      'runTime': runTime.trim(),
      'runWeek': runWeek.trim(),
      'enabled': enabled,
    };
  }
}
