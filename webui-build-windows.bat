@echo off

if exist webview-bin\build\WebView2Loader.dll ( goto deps_exist )
rmdir webview-bin
mkdir webview-bin webview-bin\libs webview-bin\build webview-bin\libs\webview2
curl -sSL "https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/1.0.1150.38" | tar -xf - -C webview-bin\libs\webview2
copy /Y webview-bin\libs\webview2\build\native\x64\WebView2Loader.dll webview-bin\build
:deps_exist

if not exist %MSYS_GCC_PATH%\nul ( 
	echo Please install msys2 and mingw64, then make sure they are present at C:\msys64\ucrt64\bin.
	echo Find here: https://www.msys2.org/#installation
) else (
	if not defined CGO_ENABLED (
		set CGO_CXXFLAGS="-I%cd%\webview-bin\libs\webview2\build\native\include"
		set CGO_LDFLAGS="-L%cd%\webview-bin\libs\webview2\build\native\x64"
		set CGO_ENABLED=1
	)
	go run -tags webui -ldflags="-H windowsgui" . ui
)
