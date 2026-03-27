class QuarkTvLoginPollResponse {
  const QuarkTvLoginPollResponse({
    required this.status,
    required this.message,
    required this.configured,
    required this.pending,
    required this.updatedAt,
  });

  final String status;
  final String message;
  final bool configured;
  final bool pending;
  final DateTime? updatedAt;

  bool get isSuccess => status == 'success';
  bool get isPending => status == 'pending';
  bool get isExpired => status == 'expired';
  bool get isError => status == 'error';

  factory QuarkTvLoginPollResponse.fromJson(Map<String, dynamic> json) {
    return QuarkTvLoginPollResponse(
      status: (json['status'] ?? '').toString().trim(),
      message: (json['message'] ?? '').toString().trim(),
      configured: _toBool(json['configured']),
      pending: _toBool(json['pending']),
      updatedAt: _toDateTime(json['updatedAt']),
    );
  }

  static bool _toBool(dynamic value) {
    if (value is bool) return value;
    if (value is num) return value != 0;
    final text = value?.toString().trim().toLowerCase() ?? '';
    return text == 'true' || text == '1';
  }

  static DateTime? _toDateTime(dynamic value) {
    if (value is String) {
      final text = value.trim();
      if (text.isEmpty) return null;
      return DateTime.tryParse(text)?.toLocal();
    }
    return null;
  }
}
