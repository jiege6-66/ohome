class MediaPath {
  const MediaPath._();

  static String normalize(String? value) {
    if (value == null) return '';
    final normalized = value.replaceAll('\\', '/').trim();
    if (normalized.isEmpty) return '';
    final parts = normalized
        .split('/')
        .where((part) => part.trim().isNotEmpty)
        .toList(growable: false);
    if (parts.isEmpty) {
      return normalized.startsWith('/') ? '/' : '';
    }
    return '/${parts.join('/')}';
  }

  static String join(String folderPath, String itemName) {
    final folder = folderPath.trim();
    final item = itemName.trim();
    if (folder.isEmpty) return normalize(item);
    if (item.isEmpty) return normalize(folder);
    if (folder.endsWith('/')) return normalize('$folder$item');
    return normalize('$folder/$item');
  }

  static bool equals(String? left, String? right) {
    final l = normalize(left);
    final r = normalize(right);
    return l.isNotEmpty && l == r;
  }

  static String title(String? path) {
    final normalized = normalize(path);
    if (normalized.isEmpty) return '';
    final parts = normalized
        .split('/')
        .where((part) => part.trim().isNotEmpty)
        .toList(growable: false);
    if (parts.isEmpty) return '';
    return parts.last;
  }
}
