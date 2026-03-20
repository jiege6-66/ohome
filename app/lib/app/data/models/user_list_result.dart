import '../../utils/common_utils.dart';
import 'user_model.dart';

class UserListResult {
  const UserListResult({required this.records, required this.total});

  final List<UserModel> records;
  final int total;

  factory UserListResult.fromJson(Map<String, dynamic> json) {
    final records = <UserModel>[];
    final rawRecords = json['records'];
    if (rawRecords is List) {
      for (final item in rawRecords) {
        final user = _decodeUser(item);
        if (user != null) {
          records.add(user);
        }
      }
    }

    return UserListResult(
      records: records,
      total: CommonUtils.toInt(json['total']) ?? records.length,
    );
  }

  static UserModel? _decodeUser(dynamic data) {
    if (data is Map<String, dynamic>) {
      return UserModel.fromJson(data);
    }
    if (data is Map) {
      return UserModel.fromJson(data.cast<String, dynamic>());
    }
    return null;
  }
}
