import '../../utils/backend_url_resolver.dart';

class UserModel {
  const UserModel({required this.raw});

  final Map<String, dynamic> raw;

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(raw: Map<String, dynamic>.from(json));
  }

  Map<String, dynamic> toJson() => Map<String, dynamic>.from(raw);

  UserModel copyWithRaw(Map<String, dynamic> patch) {
    return UserModel(raw: <String, dynamic>{...raw, ...patch});
  }

  int? get id {
    final value = raw['id'];
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }

  String get name => (raw['name'] as String?)?.trim() ?? '';

  String get realName => (raw['realName'] as String?)?.trim() ?? '';

  String get avatar => (raw['avatar'] as String?)?.trim() ?? '';

  String get avatarUrl => BackendUrlResolver.resolve(avatar);

  int? get roleId => _toInt(raw['roleId']);

  String get roleCode {
    final direct = _readString(raw['roleCode']);
    if (direct.isNotEmpty) return direct;
    final role = raw['role'];
    if (role is Map<String, dynamic>) {
      return _readString(role['code']);
    }
    if (role is Map) {
      return _readString(role['code']);
    }
    return '';
  }

  String get roleName {
    final direct = _readString(raw['roleName']);
    if (direct.isNotEmpty) return direct;
    final role = raw['role'];
    if (role is Map<String, dynamic>) {
      return _readString(role['name']);
    }
    if (role is Map) {
      return _readString(role['name']);
    }
    if (roleCode == 'super_admin') return '超级管理员';
    if (roleCode == 'user') return '普通用户';
    return '';
  }

  bool get isSuperAdmin => roleCode == 'super_admin';

  DateTime? get createdAt => _parseDateTime(raw['createdAt']);

  DateTime? get updatedAt => _parseDateTime(raw['updatedAt']);

  static DateTime? _parseDateTime(dynamic value) {
    if (value is String) {
      final text = value.trim();
      if (text.isEmpty) return null;
      return DateTime.tryParse(text)?.toLocal();
    }
    if (value is int) {
      return DateTime.fromMillisecondsSinceEpoch(value).toLocal();
    }
    if (value is double) {
      return DateTime.fromMillisecondsSinceEpoch(value.toInt()).toLocal();
    }
    return null;
  }

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }

  static String _readString(dynamic value) {
    if (value is String) return value.trim();
    return '';
  }
}
