import '../../../data/models/quark_file_entry.dart';
import '../../../utils/backend_url_resolver.dart';
import '../../common/controllers/quark_folder_controller.dart';

class AudiobookController extends WebdavFolderController {
  AudiobookController() : super(applicationType: 'xiaoshuo');

  Uri? buildStreamUri(WebdavFileEntry entry) {
    final raw = entry.resolveStreamUrl(applicationType: applicationType);
    if (raw.isEmpty) return null;
    return Uri.tryParse(BackendUrlResolver.resolve(raw));
  }
}
