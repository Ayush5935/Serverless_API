#include <Windows.h>
#include <Winspool.h>
#include <string>

// Global variables
HWND g_hListBox;
HINSTANCE g_hInstance;

// Function to update the list box with print job information
void UpdateListBox() {
    // Clear the list box
    SendMessage(g_hListBox, LB_RESETCONTENT, 0, 0);

    // Get printer handle
    HANDLE hPrinter;
    if (!OpenPrinter(NULL, &hPrinter, NULL)) {
        MessageBox(NULL, L"Failed to open printer", L"Error", MB_ICONERROR);
        return;
    }

    // Get job count and information
    DWORD dwJobCount;
    JOB_INFO_1* pJobInfo;
    if (!EnumJobs(hPrinter, 0, -1, 1, NULL, 0, &dwJobCount, NULL)) {
        MessageBox(NULL, L"Failed to enumerate jobs", L"Error", MB_ICONERROR);
        ClosePrinter(hPrinter);
        return;
    }

    pJobInfo = (JOB_INFO_1*)malloc(dwJobCount * sizeof(JOB_INFO_1));
    if (!pJobInfo) {
        MessageBox(NULL, L"Memory allocation failed", L"Error", MB_ICONERROR);
        ClosePrinter(hPrinter);
        return;
    }

    if (!EnumJobs(hPrinter, 0, dwJobCount, 1, (LPBYTE)pJobInfo, dwJobCount * sizeof(JOB_INFO_1), &dwJobCount, NULL)) {
        MessageBox(NULL, L"Failed to enumerate jobs", L"Error", MB_ICONERROR);
        ClosePrinter(hPrinter);
        free(pJobInfo);
        return;
    }

    // Add job information to the list box
    for (DWORD i = 0; i < dwJobCount; i++) {
        std::wstring jobText = L"Job ID: " + std::to_wstring(pJobInfo[i].JobId) + L", Document: " + std::wstring(pJobInfo[i].pDocument);
        SendMessage(g_hListBox, LB_ADDSTRING, 0, (LPARAM)jobText.c_str());
    }

    // Cleanup
    ClosePrinter(hPrinter);
    free(pJobInfo);
}

// Window procedure
LRESULT CALLBACK WindowProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam) {
    switch (message) {
    case WM_CREATE: {
        // Create list box
        g_hListBox = CreateWindowEx(0, L"LISTBOX", L"", WS_CHILD | WS_VISIBLE | WS_BORDER | LBS_NOSEL | LBS_HASSTRINGS,
            10, 10, 400, 300, hWnd, NULL, g_hInstance, NULL);
        if (!g_hListBox) {
            MessageBox(NULL, L"Failed to create list box", L"Error", MB_ICONERROR);
            return -1;
        }

        // Update list box
        UpdateListBox();
        break;
    }
    case WM_DESTROY:
        PostQuitMessage(0);
        break;
    default:
        return DefWindowProc(hWnd, message, wParam, lParam);
    }
    return 0;
}

// Entry point
int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) {
    g_hInstance = hInstance;

    // Register window class
    WNDCLASS wc = {};
    wc.lpfnWndProc = WindowProc;
    wc.hInstance = hInstance;
    wc.lpszClassName = L"PrintJobTracker";
    RegisterClass(&wc);

    // Create window
    HWND hWnd = CreateWindowEx(0, L"PrintJobTracker", L"Print Job Tracker", WS_OVERLAPPEDWINDOW,
        CW_USEDEFAULT, CW_USEDEFAULT, 800, 600, NULL, NULL, hInstance, NULL);
    if (!hWnd) {
        MessageBox(NULL, L"Failed to create window", L"Error", MB_ICONERROR);
        return -1;
    }

    // Show window
    ShowWindow(hWnd, nCmdShow);
    UpdateWindow(hWnd);

    // Message loop
    MSG msg = {};
    while (GetMessage(&msg, NULL, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    return 0;
}
