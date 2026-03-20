class AppUpdateInfo {
  AppUpdateInfo({
    required this.apkUrl,
    required this.versionName,
    this.versionCode,
    this.sha256checksum,
    this.forceUpdate = false,
    this.releaseNotes,
  });

  final String apkUrl;
  final String versionName;
  final int? versionCode;
  final String? sha256checksum;
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
    return '${prefix}_$safeVersion.apk';
  }
}
