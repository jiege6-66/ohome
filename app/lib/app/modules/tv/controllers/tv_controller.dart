import '../../../data/models/quark_file_entry.dart';
import '../../common/controllers/quark_folder_controller.dart';

class TvController extends WebdavFolderController {
  TvController() : super(applicationType: 'tv');

  Uri? buildStreamUri(WebdavFileEntry entry) {
    return entry.resolveStreamUri(applicationType: applicationType);
  }
}
