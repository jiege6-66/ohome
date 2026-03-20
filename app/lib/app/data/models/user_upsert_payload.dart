import 'user_model.dart';

class UserUpsertPayload {
  const UserUpsertPayload({
    this.id,
    required this.name,
    required this.realName,
    required this.roleCode,
    this.avatar = '',
  });

  final int? id;
  final String name;
  final String realName;
  final String roleCode;
  final String avatar;

  factory UserUpsertPayload.fromUser(UserModel user) {
    return UserUpsertPayload(
      id: user.id,
      name: user.name,
      realName: user.realName,
      roleCode: user.roleCode,
      avatar: user.avatar,
    );
  }

  Map<String, dynamic> toJson({bool includeId = true}) {
    final data = <String, dynamic>{
      'name': name.trim(),
      'realName': realName.trim(),
      'avatar': avatar.trim(),
      if (roleCode.trim().isNotEmpty) 'roleCode': roleCode.trim(),
    };
    if (includeId && id != null) {
      data['id'] = id;
    }
    return data;
  }

  UserModel toUserModel({Map<String, dynamic>? base}) {
    return UserModel.fromJson(<String, dynamic>{
      ...?base,
      if (id != null) 'id': id,
      'name': name.trim(),
      'realName': realName.trim(),
      'avatar': avatar.trim(),
      if (roleCode.trim().isNotEmpty) 'roleCode': roleCode.trim(),
      if (roleCode.trim().isNotEmpty) 'roleName': _roleNameFromCode(roleCode),
    });
  }

  static String _roleNameFromCode(String value) {
    switch (value.trim()) {
      case 'super_admin':
        return '超级管理员';
      case 'user':
        return '普通用户';
      default:
        return '';
    }
  }
}
