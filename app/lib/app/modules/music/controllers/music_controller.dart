import '../../../data/models/quark_file_entry.dart';
import '../../common/controllers/quark_folder_controller.dart';

class MusicController extends WebdavFolderController {
  MusicController() : super(applicationType: 'music');

  Uri? buildStreamUri(WebdavFileEntry entry) {
    return entry.resolveStreamUri(applicationType: applicationType);
  }
}
