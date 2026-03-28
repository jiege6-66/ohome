class QuarkTransferTaskModel {
  const QuarkTransferTaskModel({required this.raw});

  final Map<String, dynamic> raw;

  factory QuarkTransferTaskModel.fromJson(Map<String, dynamic> json) {
    return QuarkTransferTaskModel(raw: Map<String, dynamic>.from(json));
  }

  int? get id => _toInt(raw['id']);

  int? get sourceTaskId => _toInt(raw['sourceTaskId']);

  int get savedCount => _toInt(raw['savedCount']) ?? 0;

  String get displayName {
    final text = (raw['displayName'] ?? '').toString().trim();
    if (text.isNotEmpty) return text;
    return shareUrl;
  }

  String get shareUrl => (raw['shareUrl'] ?? '').toString().trim();

  String get savePath => (raw['savePath'] ?? '').toString().trim();

  String get displaySavePath {
    final text = savePath.replaceAll('\\', '/').trim();
    if (text.isEmpty) return '';
    return text.startsWith('/') ? text : '/$text';
  }

  String get application => (raw['application'] ?? '').toString().trim();

  String get sourceType => (raw['sourceType'] ?? '').toString().trim();

  String get status => (raw['status'] ?? '').toString().trim();

  String get resultMessage => (raw['resultMessage'] ?? '').toString().trim();

  DateTime? get createdAt => _parseDateTime(raw['createdAt']);

  DateTime? get updatedAt => _parseDateTime(raw['updatedAt']);

  DateTime? get startedAt => _parseDateTime(raw['startedAt']);

  DateTime? get finishedAt => _parseDateTime(raw['finishedAt']);

  bool get isQueued => status == 'queued';

  bool get isProcessing => status == 'processing' || status == 'queued';

  bool get isSuccess => status == 'success';

  bool get isFailed => status == 'failed';

  bool get canGoSync => isSuccess && sourceType == 'search_manual';

  String get statusLabel {
    switch (status) {
      case 'queued':
        return '排队中';
      case 'processing':
        return '转存中';
      case 'success':
        return '转存成功';
      case 'failed':
        return '转存失败';
      default:
        return status.isEmpty ? '未知状态' : status;
    }
  }

  String get sourceTypeLabel {
    switch (sourceType) {
      case 'search_manual':
        return '搜索手动转存';
      case 'sync_manual':
        return '同步任务执行一次';
      case 'sync_schedule':
        return '同步任务定时执行';
      default:
        return sourceType.isEmpty ? '未知来源' : sourceType;
    }
  }

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
