class QuarkSyncFormDraft {
  const QuarkSyncFormDraft({
    required this.taskName,
    required this.shareUrl,
    required this.savePath,
    this.application,
  });

  final String taskName;
  final String shareUrl;
  final String savePath;
  final String? application;
}
