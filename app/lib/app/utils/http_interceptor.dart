import 'package:dio/dio.dart';
import 'package:get/get.dart' hide Response;
import 'package:ohome/app/services/auth_service.dart';

class HttpInterceptor extends Interceptor {
  HttpInterceptor(this._dio);

  final Dio _dio;

  static const _kRetried = '__auth_retried__';

  bool _isRetried(RequestOptions options) => options.extra[_kRetried] == true;

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final auth = Get.find<AuthService>();
    final token = auth.accessToken.value;
    if (token != null && token.trim().isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    super.onRequest(options, handler);
  }

  @override
  Future<void> onResponse(
    Response response,
    ResponseInterceptorHandler handler,
  ) async {
    final options = response.requestOptions;
    final data = response.data;
    final code = data is Map ? data['code'] : null;
    final businessUnauthorized = code == 401 || code == '401';
    if (!businessUnauthorized) {
      handler.next(response);
      return;
    }
    final auth = Get.find<AuthService>();
    if (_isRetried(options)) {
      await _forceLogout(auth);
      handler.next(response);
      return;
    }

    try {
      await auth.refreshTokenSingleflight();
      final newToken = auth.accessToken.value;
      if (newToken == null || newToken.trim().isEmpty) {
        throw Exception('刷新后缺少 accessToken');
      }
      options.headers['Authorization'] = 'Bearer $newToken';
      options.extra[_kRetried] = true;
      final retryResponse = await _dio.fetch<dynamic>(options);
      handler.resolve(retryResponse);
    } catch (e) {
      await _forceLogout(auth);
      handler.next(response);
    }
  }

  Future<void> _forceLogout(AuthService auth) async {
    await auth.logout();
    if (Get.currentRoute != '/login') {
      Get.offAllNamed('/login');
    }
  }
}
