package iosjk.xyz.app

import android.content.Context
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.net.wifi.WifiManager
import com.ryanheise.audioservice.AudioServiceActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import java.net.Inet4Address

class MainActivity : AudioServiceActivity() {
    private val networkChannelName = "ohome/network_info"
    private var multicastLock: WifiManager.MulticastLock? = null

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, networkChannelName)
            .setMethodCallHandler { call, result ->
                when (call.method) {
                    "getActiveIpv4Cidr" -> result.success(getActiveIpv4Cidr())
                    "acquireMulticastLock" -> {
                        acquireMulticastLock()
                        result.success(null)
                    }
                    "releaseMulticastLock" -> {
                        releaseMulticastLock()
                        result.success(null)
                    }
                    else -> result.notImplemented()
                }
            }
    }

    override fun onDestroy() {
        releaseMulticastLock()
        super.onDestroy()
    }

    private fun getActiveIpv4Cidr(): Map<String, Any>? {
        val connectivityManager =
            applicationContext.getSystemService(Context.CONNECTIVITY_SERVICE) as? ConnectivityManager
                ?: return null
        val activeNetwork = connectivityManager.activeNetwork ?: return null
        val capabilities = connectivityManager.getNetworkCapabilities(activeNetwork) ?: return null
        val isLocalNetwork =
            capabilities.hasTransport(NetworkCapabilities.TRANSPORT_WIFI) ||
                capabilities.hasTransport(NetworkCapabilities.TRANSPORT_ETHERNET)
        if (!isLocalNetwork) {
            return null
        }

        val linkProperties = connectivityManager.getLinkProperties(activeNetwork) ?: return null
        val linkAddress =
            linkProperties.linkAddresses.firstOrNull { link ->
                val address = link.address
                address is Inet4Address && !address.isLoopbackAddress
            } ?: return null

        return mapOf(
            "address" to linkAddress.address.hostAddress,
            "prefixLength" to linkAddress.prefixLength,
        )
    }

    private fun acquireMulticastLock() {
        if (multicastLock?.isHeld == true) {
            return
        }
        val wifiManager =
            applicationContext.getSystemService(Context.WIFI_SERVICE) as? WifiManager ?: return
        multicastLock = wifiManager.createMulticastLock("ohome-discovery").apply {
            setReferenceCounted(false)
            acquire()
        }
    }

    private fun releaseMulticastLock() {
        val lock = multicastLock ?: return
        if (lock.isHeld) {
            lock.release()
        }
        multicastLock = null
    }
}
