import '../../../data/models/quark_file_entry.dart';
import '../../common/controllers/quark_folder_controller.dart';

class AudiobookController extends WebdavFolderController {
  AudiobookController() : super(applicationType: 'xiaoshuo');

  Uri? buildStreamUri(WebdavFileEntry entry) {
    return entry.resolveStreamUri(applicationType: applicationType);
  }
}
