import 'user_model.dart';

class LoginResult {
  const LoginResult({
    required this.accessToken,
    required this.refreshToken,
    required this.user,
  });

  final String accessToken;
  final String refreshToken;
  final UserModel user;

  factory LoginResult.fromJson(Map<String, dynamic> json) {
    final userJson = json['user'];
    return LoginResult(
      accessToken: (json['accessToken'] as String?)?.trim() ?? '',
      refreshToken: (json['refreshToken'] as String?)?.trim() ?? '',
      user: userJson is Map<String, dynamic>
          ? UserModel.fromJson(userJson)
          : const UserModel(raw: <String, dynamic>{}),
    );
  }
}
