import '../../utils/http_client.dart';
import '../models/quark_tv_login_poll_response.dart';
import '../models/quark_tv_login_start_response.dart';
import '../models/quark_tv_login_status.dart';

class QuarkTvLoginApi {
  QuarkTvLoginApi({HttpClient? httpClient})
    : _httpClient = httpClient ?? HttpClient.instance;

  final HttpClient _httpClient;

  Future<QuarkTvLoginStatus> getStatus({bool showErrorToast = true}) {
    return _httpClient.get<QuarkTvLoginStatus>(
      '/quarkLogin/tv/status',
      showErrorToast: showErrorToast,
      decoder: (data) => _decodeStatus(data),
    );
  }

  Future<QuarkTvLoginStartResponse> startLogin({bool showErrorToast = true}) {
    return _httpClient.post<QuarkTvLoginStartResponse>(
      '/quarkLogin/tv/start',
      showErrorToast: showErrorToast,
      decoder: (data) => _decodeStart(data),
    );
  }

  Future<QuarkTvLoginPollResponse> pollLogin({bool showErrorToast = true}) {
    return _httpClient.post<QuarkTvLoginPollResponse>(
      '/quarkLogin/tv/poll',
      showErrorToast: showErrorToast,
      decoder: (data) => _decodePoll(data),
    );
  }

  static QuarkTvLoginStatus _decodeStatus(dynamic data) {
    if (data is Map<String, dynamic>) {
      return QuarkTvLoginStatus.fromJson(data);
    }
    if (data is Map) {
      return QuarkTvLoginStatus.fromJson(data.cast<String, dynamic>());
    }
    return const QuarkTvLoginStatus(
      configured: false,
      pending: false,
      updatedAt: null,
    );
  }

  static QuarkTvLoginStartResponse _decodeStart(dynamic data) {
    if (data is Map<String, dynamic>) {
      return QuarkTvLoginStartResponse.fromJson(data);
    }
    if (data is Map) {
      return QuarkTvLoginStartResponse.fromJson(data.cast<String, dynamic>());
    }
    return const QuarkTvLoginStartResponse(qrData: '', pending: false);
  }

  static QuarkTvLoginPollResponse _decodePoll(dynamic data) {
    if (data is Map<String, dynamic>) {
      return QuarkTvLoginPollResponse.fromJson(data);
    }
    if (data is Map) {
      return QuarkTvLoginPollResponse.fromJson(data.cast<String, dynamic>());
    }
    return const QuarkTvLoginPollResponse(
      status: 'error',
      message: '',
      configured: false,
      pending: false,
      updatedAt: null,
    );
  }
}
