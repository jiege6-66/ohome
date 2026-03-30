import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

class BackendAddressBadge extends StatelessWidget {
  const BackendAddressBadge({super.key, required this.address});

  final String address;

  @override
  Widget build(BuildContext context) {
    final isDesktop = _isNativeDesktopPlatform() ||
        MediaQuery.sizeOf(context).width >= 900;

    return ConstrainedBox(
      constraints: BoxConstraints(maxWidth: isDesktop ? 320 : 220.w),
      child: Container(
        padding: EdgeInsets.symmetric(
          horizontal: isDesktop ? 12 : 12.w,
          vertical: isDesktop ? 10 : 10.h,
        ),
        decoration: BoxDecoration(
          color: Colors.white.withValues(alpha: 0.08),
          borderRadius: BorderRadius.circular(isDesktop ? 12 : 14.r),
          border: Border.all(color: Colors.white12),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.dns_rounded,
              size: isDesktop ? 15 : 16.sp,
              color: Colors.white70,
            ),
            SizedBox(width: isDesktop ? 8 : 8.w),
            Expanded(
              child: Text(
                address,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  fontSize: isDesktop ? 12 : 12.sp,
                  fontWeight: FontWeight.w500,
                  color: Colors.white,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

bool _isNativeDesktopPlatform() =>
    !kIsWeb &&
    (defaultTargetPlatform == TargetPlatform.macOS ||
        defaultTargetPlatform == TargetPlatform.windows ||
        defaultTargetPlatform == TargetPlatform.linux);
