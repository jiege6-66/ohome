import '../../../data/models/quark_file_entry.dart';
import '../../common/controllers/quark_folder_controller.dart';

class PlayLetController extends WebdavFolderController {
  PlayLetController() : super(applicationType: 'playlet');

  Uri? buildStreamUri(WebdavFileEntry entry) {
    return entry.resolveStreamUri(applicationType: applicationType);
  }
}
