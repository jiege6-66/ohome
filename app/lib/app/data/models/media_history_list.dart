import 'media_history_entry.dart';

class MediaHistoryListResult {
  const MediaHistoryListResult({required this.records, required this.total});

  final List<MediaHistoryEntry> records;
  final int total;
}
