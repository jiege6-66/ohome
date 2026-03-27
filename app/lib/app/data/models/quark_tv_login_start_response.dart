class QuarkTvLoginStartResponse {
  const QuarkTvLoginStartResponse({
    required this.qrData,
    required this.pending,
  });

  final String qrData;
  final bool pending;

  factory QuarkTvLoginStartResponse.fromJson(Map<String, dynamic> json) {
    return QuarkTvLoginStartResponse(
      qrData: (json['qrData'] ?? '').toString().trim(),
      pending: _toBool(json['pending']),
    );
  }

  static bool _toBool(dynamic value) {
    if (value is bool) return value;
    if (value is num) return value != 0;
    final text = value?.toString().trim().toLowerCase() ?? '';
    return text == 'true' || text == '1';
  }
}
