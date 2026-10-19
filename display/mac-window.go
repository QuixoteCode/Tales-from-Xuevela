package display

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#import <stdlib.h>

// Keep a reference to the text view so we can update it later.
static NSTextView *outputView = nil;

// Create and configure the native macOS window.
static void openNativeWindow(void) {
    @autoreleasepool {
        // Initialize the Cocoa application.
        [NSApplication sharedApplication];

        // Make the application behave like a normal desktop app.
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

        // Create a 900 x 600 window.
        NSRect frame = NSMakeRect(0, 0, 900, 600);

        // Enable the title bar, close, minimize, and resize buttons.
        NSWindowStyleMask style =
            NSWindowStyleMaskTitled |
            NSWindowStyleMaskClosable |
            NSWindowStyleMaskMiniaturizable |
            NSWindowStyleMaskResizable;

        NSWindow *window = [[NSWindow alloc]
            initWithContentRect:frame
                      styleMask:style
                        backing:NSBackingStoreBuffered
                          defer:NO];

        // Set the title displayed in the window's title bar.
        [window setTitle:@"My Adventure"];

        // Center the window on the screen.
        [window center];

        // Create a scrollable area for the game's text.
        NSScrollView *scrollView = [[NSScrollView alloc]
            initWithFrame:NSMakeRect(20, 20, 860, 540)];

        // Display a vertical scrollbar when needed.
        [scrollView setHasVerticalScroller:YES];

        // Allow the scroll view to resize with the window.
        [scrollView setAutoresizingMask:
            NSViewWidthSizable | NSViewHeightSizable];

        // Create the text area that will contain the story.
        outputView = [[NSTextView alloc]
            initWithFrame:NSMakeRect(0, 0, 840, 540)];

        // The player can select and copy text, but not edit it.
        [outputView setEditable:NO];
        [outputView setSelectable:YES];

        // Make the text area grow vertically as text is added.
        [outputView setVerticallyResizable:YES];
        [outputView setHorizontallyResizable:NO];

        // Use a readable font for the story.
        [outputView setFont:
            [NSFont systemFontOfSize:15]];

        // Give the text some space around its edges.
        [outputView setTextContainerInset:NSMakeSize(8, 8)];

        // Put the text view inside the scroll view.
        [scrollView setDocumentView:outputView];

        // Add the scroll view to the window.
        [[window contentView] addSubview:scrollView];

        // Show the window and bring it to the foreground.
        [window makeKeyAndOrderFront:nil];
        [NSApp activateIgnoringOtherApps:YES];
    }
}

// Append a line of text to the game's output area.
static void appendNativeText(const char *message) {
    @autoreleasepool {
        if (outputView == nil || message == NULL) {
            return;
        }

        // Convert the UTF-8 text received from Go into an NSString.
        NSString *text = [NSString stringWithUTF8String:message];

        if (text == nil) {
            return;
        }

        // Append the message and a newline to the existing content.
        NSAttributedString *line = [[NSAttributedString alloc]
            initWithString:[text stringByAppendingString:@"\n"]];

        [[outputView textStorage] appendAttributedString:line];

        // Scroll down so the newest message is visible.
        [outputView scrollRangeToVisible:
            NSMakeRange([[outputView string] length], 0)];
    }
}

// Start Cocoa's event loop.
// This call remains active while the application is running.
static void runNativeApplication(void) {
    [NSApp run];
}
*/
import "C"

import "unsafe"

// OpenWindow creates the native macOS window.
// It sets up the UI but does not start the event loop.
func OpenWindow() {
    C.openNativeWindow()
}

// Println appends a line of text to the window.
// This initial version assumes it is called on the UI thread.
func Println(message string) {
    cMessage := C.CString(message)
    defer C.free(unsafe.Pointer(cMessage))

    C.appendNativeText(cMessage)
}

// Run starts the native application's event loop.
// Call this after creating the window.
func Run() {
    C.runNativeApplication()
}