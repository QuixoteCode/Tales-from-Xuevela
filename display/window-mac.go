package display

// The code is in a comment because Go's parser does not understand Objective-C syntax. cgo extracts and compiles the code separately.
/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

// Objective-C implementation
static void openNativeWindow(void) {
    @autoreleasepool { // Manages Objective-C autoreleased objects.
        [NSApplication sharedApplication]; // Initializes the macOS application object.

        [NSApp setActivationPolicy: NSApplicationActivationPolicyRegular]; // Makes the application behave like a regular GUI application.

        NSRect frame = NSMakeRect(0, 0, 900, 600); // Defines an 900 × 600 window frame.

        NSWindowStyleMask style =
            NSWindowStyleMaskTitled | // Enables the title bar
            NSWindowStyleMaskClosable | // Enables the close button
            NSWindowStyleMaskMiniaturizable | // Enables the minimize button
            NSWindowStyleMaskResizable; // Resizing

        NSWindow *window = [[NSWindow alloc] initWithContentRect:frame styleMask:style backing:NSBackingStoreBuffered defer:NO]; // Creates the window with the specified frame and style.

        [window setTitle:@"Go Native Window"]; // Sets the title
        [window center]; // Centers the window
        [window makeKeyAndOrderFront:nil]; // Displays the window

        [NSApp activateIgnoringOtherApps:YES]; // Requests that the application become active
        [NSApp run]; // Starts the application's event loop, allowing it to process user input and window events.
    }
}
*/
import "C"

// OpenWindow creates and displays a native macOS window.
func OpenWindow() {
    C.openNativeWindow()
}