import 'app_env.dart';

class BackendUrlResolver {
  static String resolve(String rawUrl) {
    final raw = rawUrl.trim();
    if (raw.isEmpty) return '';

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

  static String _joinPath(String basePath, String child) {
    final normalizedBase = basePath.trim();
    if (normalizedBase.isEmpty) return '/$child';
    if (normalizedBase.endsWith('/')) return '$normalizedBase$child';
    return '$normalizedBase/$child';
  }
}
