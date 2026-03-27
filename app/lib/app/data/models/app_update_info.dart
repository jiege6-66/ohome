class AppUpdateInfo {
  AppUpdateInfo({
    required this.apkUrl,
    required this.versionName,
    this.versionCode,
    this.sha256checksum,
    this.artifactKey,
    this.forceUpdate = false,
    this.releaseNotes,
  });

  final String apkUrl;
  final String versionName;
  final int? versionCode;
  final String? sha256checksum;
  final String? artifactKey;
  final bool forceUpdate;
  final String? releaseNotes;

  String get displayVersion {
    final version = versionName.trim();
    if (versionCode == null) return version;
    if (version.isEmpty) return versionCode.toString();
    return '$version+$versionCode';
  }

  String buildDestinationFilename({String prefix = 'ohome'}) {
    final safeVersion = versionName.trim().isEmpty
        ? 'latest'
        : versionName.trim().replaceAll(RegExp(r'[^0-9A-Za-z._-]'), '_');
    final safeArtifact = artifactKey?.trim().isNotEmpty == true
        ? artifactKey!.trim().replaceAll(RegExp(r'[^0-9A-Za-z._-]'), '_')
        : null;
    final artifactSuffix = safeArtifact == null ? '' : '_$safeArtifact';
    return '${prefix}_$safeVersion$artifactSuffix.apk';
  }
}
