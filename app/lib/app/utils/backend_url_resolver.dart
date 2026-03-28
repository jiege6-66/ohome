import 'package:get/get.dart';

import '../services/auth_service.dart';
import 'app_env.dart';

class BackendUrlResolver {
  static String resolve(String rawUrl) {
    final raw = rawUrl.trim();
    if (raw.isEmpty) return '';

    final resolved = _resolveAbsoluteUrl(raw);
    return _applyAccessTokenIfNeeded(resolved);
  }

  static String _resolveAbsoluteUrl(String raw) {
    final uri = Uri.tryParse(raw);
    if (uri != null && uri.hasScheme && uri.host.trim().isNotEmpty) {
      return uri.toString();
    }

    final base = Uri.parse(AppEnv.instance.apiBaseUrl);
    if (uri == null) {
      return raw;
    }

    final rawPath = uri.path.trim();
    if (rawPath.isEmpty) {
      return raw;
    }

    if (rawPath.startsWith('/public/')) {
      final mergedPath = _joinPath(base.path, rawPath.substring(1));
      return base
          .replace(path: mergedPath, query: uri.hasQuery ? uri.query : null)
          .toString();
    }

    return base.resolve(raw).toString();
  }

  static String _applyAccessTokenIfNeeded(String resolvedUrl) {
    final uri = Uri.tryParse(resolvedUrl.trim());
    if (uri == null || !_isProtectedBackendStream(uri)) {
      return resolvedUrl;
    }

    final nextQuery = <String, String>{...uri.queryParameters};
    final accessToken = _currentAccessToken();
    if (accessToken == null) {
      nextQuery.remove('access_token');
    } else {
      nextQuery['access_token'] = accessToken;
    }
    return uri.replace(queryParameters: nextQuery).toString();
  }

  static bool _isProtectedBackendStream(Uri uri) {
    if (!_isCurrentBackendUri(uri)) return false;

    final normalizedPath = uri.path.replaceAll('\\', '/').trim();
    return (normalizedPath.contains('/public/quarkFs/') ||
            normalizedPath.contains('/quarkFs/')) &&
        normalizedPath.endsWith('/files/stream');
  }

  static bool _isCurrentBackendUri(Uri uri) {
    if (!uri.hasScheme || uri.host.trim().isEmpty) return false;

    final base = Uri.parse(AppEnv.instance.apiBaseUrl);
    final uriScheme = uri.scheme.trim().toLowerCase();
    final baseScheme = base.scheme.trim().toLowerCase();
    if (uriScheme != baseScheme) return false;

    final uriHost = uri.host.trim().toLowerCase();
    final baseHost = base.host.trim().toLowerCase();
    if (uriHost != baseHost) return false;

    final uriPort = uri.hasPort ? uri.port : _defaultPortForScheme(uriScheme);
    final basePort = base.hasPort
        ? base.port
        : _defaultPortForScheme(baseScheme);
    return uriPort == basePort;
  }

  static int? _defaultPortForScheme(String scheme) {
    switch (scheme) {
      case 'http':
        return 80;
      case 'https':
        return 443;
      default:
        return null;
    }
  }

  static String? _currentAccessToken() {
    if (!Get.isRegistered<AuthService>()) return null;

    final token = Get.find<AuthService>().accessToken.value?.trim() ?? '';
    if (token.isEmpty) return null;
    return token;
  }

  static String _joinPath(String basePath, String child) {
    final normalizedBase = basePath.trim();
    if (normalizedBase.isEmpty) return '/$child';
    if (normalizedBase.endsWith('/')) return '$normalizedBase$child';
    return '$normalizedBase/$child';
  }
}
