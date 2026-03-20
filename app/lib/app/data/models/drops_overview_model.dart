import 'drops_event_model.dart';
import 'drops_item_model.dart';

class DropsOverviewModel {
  const DropsOverviewModel({required this.raw});

  final Map<String, dynamic> raw;

  factory DropsOverviewModel.fromJson(Map<String, dynamic> json) {
    return DropsOverviewModel(raw: Map<String, dynamic>.from(json));
  }

  int get todayTodoCount => _toInt(raw['todayTodoCount']) ?? 0;

  int get expiringSoonCount => _toInt(raw['expiringSoonCount']) ?? 0;

  int get monthEventCount => _toInt(raw['monthEventCount']) ?? 0;

  int get unreadMessageCount => _toInt(raw['unreadMessageCount']) ?? 0;

  List<DropsItemModel> get recentItems {
    final data = raw['recentItems'];
    if (data is! List) return const <DropsItemModel>[];
    return data
        .whereType<Map>()
        .map((item) => DropsItemModel.fromJson(item.cast<String, dynamic>()))
        .toList(growable: false);
  }

  List<DropsEventModel> get recentEvents {
    final data = raw['recentEvents'];
    if (data is! List) return const <DropsEventModel>[];
    return data
        .whereType<Map>()
        .map((item) => DropsEventModel.fromJson(item.cast<String, dynamic>()))
        .toList(growable: false);
  }

  static int? _toInt(dynamic value) {
    if (value is int) return value;
    if (value is String) return int.tryParse(value);
    return null;
  }
}
