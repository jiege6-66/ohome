class AppMessagePushEvent {
  const AppMessagePushEvent({
    required this.event,
    required this.unreadCount,
    required this.reason,
    required this.messageIds,
  });

  final String event;
  final int unreadCount;
  final String reason;
  final List<int> messageIds;

  factory AppMessagePushEvent.fromJson(Map<String, dynamic> json) {
    final rawMessageIds = json['messageIds'];
    final messageIds = rawMessageIds is List
        ? rawMessageIds.map(_toInt).whereType<int>().toList(growable: false)
        : const <int>[];

    return AppMessagePushEvent(
      event: (json['event'] ?? '').toString().trim(),
      unreadCount: _toInt(json['unreadCount']) ?? 0,
      reason: (json['reason'] ?? '').toString().trim(),
      messageIds: messageIds,
    );
  }

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.toInt();
    if (value is String) return int.tryParse(value);
    return null;
  }
}
