# Keep Flutter embedding classes used by release runtime.
-keep class io.flutter.** { *; }
-keep class io.flutter.embedding.** { *; }
-dontwarn io.flutter.embedding.**

# Keep custom Application class.
-keep class iosjk.xyz.app.OhomeApplication { *; }
