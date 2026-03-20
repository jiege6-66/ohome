import 'package:ohome/app/utils/http_client.dart';

class HomeApi {
  static Future<HitokotoModel> getDailyQuote() async {
    final response = await HttpClient.instance.get(
      'https://v1.hitokoto.cn/',
      decoder: (data) => HitokotoModel.fromJson(data),
    );
    return response;
  }
}

/// 一言数据模型
class HitokotoModel {
  int? id;
  String? uuid;
  String? hitokoto;
  String? type;
  String? from;
  String? fromWho;
  String? creator;
  int? creatorUid;
  int? reviewer;
  String? commitFrom;
  String? createdAt;
  int? length;

  HitokotoModel({
    this.id,
    this.uuid,
    this.hitokoto,
    this.type,
    this.from,
    this.fromWho,
    this.creator,
    this.creatorUid,
    this.reviewer,
    this.commitFrom,
    this.createdAt,
    this.length,
  });

  HitokotoModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    uuid = json['uuid'];
    hitokoto = json['hitokoto'];
    type = json['type'];
    from = json['from'];
    fromWho = json['from_who'];
    creator = json['creator'];
    creatorUid = json['creator_uid'];
    reviewer = json['reviewer'];
    commitFrom = json['commit_from'];
    createdAt = json['created_at']?.toString();
    length = json['length'];
  }
}
