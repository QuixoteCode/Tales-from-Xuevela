//go:build windows

package display

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const (
	windowWidth  = 900
	windowHeight = 600

	WM_DESTROY      = 0x0002
	WM_SIZE         = 0x0005
	WM_CLOSE        = 0x0010
	WM_COMMAND      = 0x0111
	WM_CTLCOLOREDIT = 0x0133
	WM_APP_TEXT     = 0x8001
	WM_APP_INPUT    = 0x8002

	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_VSCROLL          = 0x00200000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000

	ES_MULTILINE   = 0x0004
	ES_AUTOVSCROLL = 0x0040
	ES_READONLY    = 0x0800

	BS_PUSHBUTTON = 0x00000000

	EM_SETSEL     = 0x00B1
	EM_REPLACESEL = 0x00C2
	WM_SETFONT    = 0x0030

	SW_SHOW = 5

	CW_USEDEFAULT int32 = -2147483648

	ID_OK     = 1001
	ID_CANCEL = 1002

	WS_EX_DLGMODALFRAME = 0x00000001

	SS_LEFT       = 0x00000000
	SS_EDITCONTROL = 0x00002000

	DT_CALCRECT   = 0x00000400
	DT_WORDBREAK  = 0x00000010
	DT_NOPREFIX   = 0x00000800

	WM_GETTEXT    = 0x000D
	WM_GETTEXTLENGTH = 0x000E

	SM_CXSCREEN = 0
	SM_CYSCREEN = 1

	DT_SINGLELINE = 0x00000020
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW = user32.NewProc("CreateWindowExW")
	procDefWindowProcW  = user32.NewProc("DefWindowProcW")
	procGetMessageW     = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procPostMessageW    = user32.NewProc("PostMessageW")
	procSendMessageW    = user32.NewProc("SendMessageW")
	procShowWindow      = user32.NewProc("ShowWindow")
	procUpdateWindow    = user32.NewProc("UpdateWindow")
	procDestroyWindow   = user32.NewProc("DestroyWindow")
	procMoveWindow      = user32.NewProc("MoveWindow")
	procGetClientRect   = user32.NewProc("GetClientRect")
	procEnableWindow    = user32.NewProc("EnableWindow")
	procSetFocus        = user32.NewProc("SetFocus")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procSetTextColor    = gdi32.NewProc("SetTextColor")
	procSetBkColor      = gdi32.NewProc("SetBkColor")
	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procGetStockObject  = gdi32.NewProc("GetStockObject")
	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procPostQuitMessage = user32.NewProc("PostQuitMessage")

	procGetDC          = user32.NewProc("GetDC")
	procReleaseDC      = user32.NewProc("ReleaseDC")
	procDrawTextW      = user32.NewProc("DrawTextW")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procSelectObject   = gdi32.NewProc("SelectObject")
)

type point struct {
	X int32
	Y int32
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type message struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

type windowClass struct {
	CbSize        uint32
	Style         uint32
	WndProc       uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	Instance      uintptr
	Icon          uintptr
	Cursor        uintptr
	Background    uintptr
	MenuName      *uint16
	ClassName     *uint16
	IconSm        uintptr
}

type inputRequest struct {
	prompt string
	result chan string
}

var (
	mainWindow uintptr
	outputView uintptr

	promptWindow uintptr
	promptEdit   uintptr

	instance uintptr
	darkBrush uintptr

	mainThreadReady = make(chan struct{})

	mu sync.Mutex

	nextID uintptr = 1
	textQueue = make(map[uintptr]string)
	inputQueue = make(map[uintptr]*inputRequest)

	activeInput *inputRequest
)

var (
	mainClassName = syscall.StringToUTF16Ptr("MyAdventureWindow")
	promptClassName = syscall.StringToUTF16Ptr("MyAdventurePrompt")
	editClassName = syscall.StringToUTF16Ptr("EDIT")
	buttonClassName = syscall.StringToUTF16Ptr("BUTTON")
)

func getID() uintptr {
	mu.Lock()
	defer mu.Unlock()

	id := nextID
	nextID++

	return id
}

func createWindow(
	exStyle uintptr,
	className *uint16,
	title *uint16,
	style uintptr,
	x, y, width, height int32,
	parent uintptr,
	menu uintptr,
) uintptr {

	hwnd, _, _ := procCreateWindowExW.Call(
		exStyle,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		style,
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		parent,
		menu,
		instance,
		0,
	)

	return hwnd
}

func createControl(
	className *uint16,
	title string,
	style uintptr,
	x, y, width, height int32,
	parent uintptr,
	id uintptr,
) uintptr {

	titlePtr := syscall.StringToUTF16Ptr(title)

	return createWindow(
		0,
		className,
		titlePtr,
		style,
		x, y, width, height,
		parent,
		id,
	)
}

// Main window procedure
var mainWndProc = syscall.NewCallback(func(
	hwnd, msg, wParam, lParam uintptr,
) uintptr {

	switch uint32(msg) {

	case WM_SIZE:
		var r rect

		procGetClientRect.Call(
			hwnd,
			uintptr(unsafe.Pointer(&r)),
		)

		procMoveWindow.Call(
			outputView,
			20,
			20,
			uintptr(r.Right-r.Left-40),
			uintptr(r.Bottom-r.Top-40),
			1,
		)

		return 0

	case WM_APP_TEXT:

		mu.Lock()
		text := textQueue[wParam]
		delete(textQueue, wParam)
		mu.Unlock()

		if text != "" {
			textPtr := syscall.StringToUTF16Ptr(text + "\r\n")

			procSendMessageW.Call(
				outputView,
				EM_SETSEL,
				^uintptr(0),
				^uintptr(0),
			)

			procSendMessageW.Call(
				outputView,
				EM_REPLACESEL,
				0,
				uintptr(unsafe.Pointer(textPtr)),
			)
		}

		return 0

	case WM_APP_INPUT:

		mu.Lock()
		request := inputQueue[wParam]
		delete(inputQueue, wParam)
		mu.Unlock()

		if request != nil {
			showInputDialog(request)
		}

		return 0

	case WM_CTLCOLOREDIT:

		procSetTextColor.Call(wParam, 0x00FFFFFF)
		procSetBkColor.Call(wParam, 0x001F1F1F)

		return darkBrush

	case WM_CLOSE:
		procDestroyWindow.Call(hwnd)
		return 0

	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	result, _, _ := procDefWindowProcW.Call(
		hwnd, msg, wParam, lParam,
	)

	return result
})

// Input dialog window procedure
var promptWndProc = syscall.NewCallback(func(
	hwnd, msg, wParam, lParam uintptr,
) uintptr {

	switch uint32(msg) {

	case WM_COMMAND:

		switch wParam & 0xFFFF {

		case ID_OK:
			finishInput(true)

		case ID_CANCEL:
			finishInput(false)
		}

		return 0

	case WM_CLOSE:
		finishInput(false)
		return 0
	}

	result, _, _ := procDefWindowProcW.Call(
		hwnd, msg, wParam, lParam,
	)

	return result
})

func registerWindowClass(
	name *uint16,
	procedure uintptr,
) {

	class := windowClass{
		CbSize:    uint32(unsafe.Sizeof(windowClass{})),
		WndProc:   procedure,
		Instance:  instance,
		ClassName: name,
	}

	result, _, err := procRegisterClassExW.Call(
		uintptr(unsafe.Pointer(&class)),
	)

	if result == 0 {
		panic("Could not register Windows window class: " + err.Error())
	}
}

// OpenWindow initializes the Windows display
func OpenWindow() {
	// Window creation happens in Run
}

// Run creates the native window and starts the Windows message loop
func Run() {

	runtime.LockOSThread()

	defer runtime.UnlockOSThread()

	instance, _, _ = procGetModuleHandleW.Call(0)

	darkBrush, _, _ = procCreateSolidBrush.Call(0x001F1F1F)

	registerWindowClass(mainClassName, mainWndProc)
	registerWindowClass(promptClassName, promptWndProc)

	title := syscall.StringToUTF16Ptr("My Adventure")

	mainWindow = createWindow(
		0,
		mainClassName,
		title,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		windowWidth,
		windowHeight,
		0,
		0,
	)

	if mainWindow == 0 {
		panic("Could not create main window")
	}

	// Create the story output area
	outputView = createControl(
		editClassName,
		"",
		WS_CHILD|WS_VISIBLE|WS_VSCROLL|
			ES_MULTILINE|ES_AUTOVSCROLL|ES_READONLY,
		20, 20, 860, 540,
		mainWindow,
		0,
	)

	if outputView == 0 {
		panic("Could not create output text area")
	}

	// Use the standard Windows GUI font
	font, _, _ := procGetStockObject.Call(17)

	procSendMessageW.Call(
		outputView,
		WM_SETFONT,
		font,
		1,
	)

	procShowWindow.Call(mainWindow, SW_SHOW)
	procUpdateWindow.Call(mainWindow)

	close(mainThreadReady)

	// Windows message loop
	var msg message

	for {
		result, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)

		if int32(result) <= 0 {
			break
		}

		procTranslateMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
		)

		procDispatchMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
		)
	}
}

// Printdisplay appends text to the native output area
func Printdisplay(text string) {

	<-mainThreadReady

	id := getID()

	mu.Lock()
	textQueue[id] = text
	mu.Unlock()

	result, _, _ := procPostMessageW.Call(
		mainWindow,
		WM_APP_TEXT,
		id,
		0,
	)

	if result == 0 {
		mu.Lock()
		delete(textQueue, id)
		mu.Unlock()
	}
}

// Input requests text from the player
func Input(prompt string) string {

	<-mainThreadReady

	request := &inputRequest{
		prompt: prompt,
		result: make(chan string, 1),
	}

	id := getID()

	mu.Lock()
	inputQueue[id] = request
	mu.Unlock()

	result, _, _ := procPostMessageW.Call(
		mainWindow,
		WM_APP_INPUT,
		id,
		0,
	)

	if result == 0 {
		mu.Lock()
		delete(inputQueue, id)
		mu.Unlock()
		return ""
	}

	return <-request.result
}

// Create and display the input dialog
func showInputDialog(request *inputRequest) {

    activeInput = request

    procEnableWindow.Call(mainWindow, 0)

    title := syscall.StringToUTF16Ptr("My Adventure")

    // Calculate the required dialog dimensions
    dialogWidth, promptHeight := calculatePromptSize(request.prompt)

    // Account for the input field, buttons and spacing
    editY := promptHeight + 24
    buttonY := editY + 40

    dialogHeight := buttonY + 50

    // Account approximately for the Windows title bar and borders
    outerWidth := dialogWidth + 16
    outerHeight := dialogHeight + 39

    promptWindow = createWindow(
        WS_EX_DLGMODALFRAME,
        promptClassName,
        title,
        WS_CAPTION|WS_SYSMENU,
        CW_USEDEFAULT,
        CW_USEDEFAULT,
        outerWidth,
        outerHeight,
        mainWindow,
        0,
    )

    if promptWindow == 0 {
        finishInput(false)
        return
    }

    // Display the prompt with automatic word wrapping
    createControl(
        syscall.StringToUTF16Ptr("STATIC"),
        request.prompt,
        WS_CHILD|WS_VISIBLE|SS_LEFT|SS_EDITCONTROL,
        15, 12,
        dialogWidth-30,
        promptHeight,
        promptWindow,
        0,
    )

	promptLabel := createControl(
        syscall.StringToUTF16Ptr("STATIC"),
        request.prompt,
        WS_CHILD|WS_VISIBLE|SS_LEFT|SS_EDITCONTROL,
        15, 12,
        dialogWidth-30,
        promptHeight,
        promptWindow,
        0,
    )

    font, _, _ := procGetStockObject.Call(17)

    procSendMessageW.Call(
        promptLabel,
        WM_SETFONT,
        font,
        1,
    )

    // Input field
    promptEdit = createControl(
        editClassName,
        "",
        WS_CHILD|WS_VISIBLE|0x0080,
        15, editY,
        dialogWidth-30,
        25,
        promptWindow,
        0,
    )

    // OK button
    createControl(
        buttonClassName,
        "OK",
        WS_CHILD|WS_VISIBLE|BS_PUSHBUTTON,
        dialogWidth-180,
        buttonY,
        75, 28,
        promptWindow,
        ID_OK,
    )

    // Cancel button
    createControl(
        buttonClassName,
        "Cancel",
        WS_CHILD|WS_VISIBLE|BS_PUSHBUTTON,
        dialogWidth-90,
        buttonY,
        75, 28,
        promptWindow,
        ID_CANCEL,
    )

    procShowWindow.Call(promptWindow, SW_SHOW)
    procUpdateWindow.Call(promptWindow)
    procSetFocus.Call(promptEdit)
}

// Complete the input request and close the dialog
func finishInput(accepted bool) {

	if activeInput == nil {
		return
	}

	result := ""

	if accepted {

		const maxCharacters = 4096

		buffer := make([]uint16, maxCharacters)

		procSendMessageW.Call(
			promptEdit,
			0x000D, // WM_GETTEXT
			uintptr(len(buffer)),
			uintptr(unsafe.Pointer(&buffer[0])),
		)

		result = syscall.UTF16ToString(buffer)
	}

	activeInput.result <- result

	activeInput = nil

	procDestroyWindow.Call(promptWindow)

	promptWindow = 0
	promptEdit = 0

	procEnableWindow.Call(mainWindow, 1)
	procSetForegroundWindow.Call(mainWindow)
}


func calculatePromptSize(prompt string) (int32, int32) {

    const (
        defaultMinWidth = 380
        padding = 30
    )

    // Get available screen dimensions
    screenWidth, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)

    // Get the main window's available client width
    var mainRect rect

    procGetClientRect.Call(
        mainWindow,
        uintptr(unsafe.Pointer(&mainRect)),
    )

    availableWidth := mainRect.Right - mainRect.Left - 60

    // Do not exceed the available screen width
    maxWidth := int32(screenWidth) - 80

    if availableWidth > maxWidth {
        availableWidth = maxWidth
    }

    if availableWidth < 200 {
        availableWidth = 200
    }

    minWidth := int32(defaultMinWidth)

    if availableWidth < minWidth {
        minWidth = availableWidth
    }

    // Obtain the device context
    hdc, _, _ := procGetDC.Call(mainWindow)

    if hdc == 0 {
        return minWidth, 80
    }

    defer procReleaseDC.Call(mainWindow, hdc)

    // Use the standard Windows GUI font
    font, _, _ := procGetStockObject.Call(17)

    oldFont, _, _ := procSelectObject.Call(hdc, font)

    defer procSelectObject.Call(hdc, oldFont)

    text := syscall.StringToUTF16Ptr(prompt)

    // First, calculate the natural width of the text
    naturalRect := rect{
        Left:   0,
        Top:    0,
        Right:  10000,
        Bottom: 0,
    }

    procDrawTextW.Call(
        hdc,
        uintptr(unsafe.Pointer(text)),
        ^uintptr(0),
        uintptr(unsafe.Pointer(&naturalRect)),
        DT_CALCRECT|DT_SINGLELINE|DT_NOPREFIX,
    )

    width := naturalRect.Right + padding

    // Restrict the width to the available window space
    if width < minWidth {
        width = minWidth
    }

    if width > availableWidth {
        width = availableWidth
    }

    // Calculate the height after wrapping the text
    textRect := rect{
        Left:   0,
        Top:    0,
        Right:  width - 30,
        Bottom: 0,
    }

    procDrawTextW.Call(
        hdc,
        uintptr(unsafe.Pointer(text)),
        ^uintptr(0),
        uintptr(unsafe.Pointer(&textRect)),
        DT_CALCRECT|DT_WORDBREAK|DT_NOPREFIX,
    )

    height := textRect.Bottom

    return width, height
}