import '../data/storage/discovery_storage.dart';
import '../data/storage/token_storage.dart';
import '../data/storage/user_storage.dart';
import '../utils/app_env.dart';
import 'discovery_service.dart';

class StartupRecoveryResult {
  const StartupRecoveryResult({
    required this.shouldOpenLogin,
    this.fallbackApplied = false,
  });

  final bool shouldOpenLogin;
  final bool fallbackApplied;
}

class StartupRecoveryService {
  StartupRecoveryService({
    DiscoveryService? discoveryService,
    TokenStorage? tokenStorage,
    UserStorage? userStorage,
  }) : _discoveryService =
           discoveryService ??
           DiscoveryService(storage: DiscoveryStorage()),
       _tokenStorage = tokenStorage ?? TokenStorage(),
       _userStorage = userStorage ?? UserStorage();

  final DiscoveryService _discoveryService;
  final TokenStorage _tokenStorage;
  final UserStorage _userStorage;

  Future<StartupRecoveryResult> recover() async {
    final env = AppEnv.instance;
    if (env.isUsingDefaultApiBaseUrl) {
      return const StartupRecoveryResult(shouldOpenLogin: false);
    }

    final currentServer = await _discoveryService.probeApiBaseUrlInput(
      env.apiBaseUrlInputValue,
    );
    if (currentServer != null) {
      return const StartupRecoveryResult(shouldOpenLogin: false);
    }

    await env.resetApiBaseUrlToDefault();
    await _clearLocalSession();
    return const StartupRecoveryResult(
      shouldOpenLogin: true,
      fallbackApplied: true,
    );
  }

  Future<void> _clearLocalSession() async {
    await _tokenStorage.clear();
    await _userStorage.clear();
  }
}
