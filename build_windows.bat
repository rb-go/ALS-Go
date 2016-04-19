@TITLE Building GO-RPC-Server for Windows

@Echo off
echo Get needed packages
go get github.com/gorilla/rpc/v2
go get github.com/gorilla/rpc/v2/json2
go get github.com/go-ini/ini
go get launchpad.net/xmlpath
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
go get github.com/go-ozzo/ozzo-log
go get github.com/go-ozzo/ozzo-config
echo Geting packages finished

set GOARCH=amd64
set GOOS=windows

set PRJROOTPATH=%~dp0
set OUTPUTPATH=%PRJROOTPATH%out
set GOPATH=%GOPATH%;%PRJROOTPATH%
echo Building binary started...
go build -o "%OUTPUTPATH%\GoRPC.exe"
echo Building binary finished!
cp .\config.ini "%OUTPUTPATH%\config.ini"
cp .\logger.json "%OUTPUTPATH%\logger.json"
echo Compressing binary started...
..\..\..\tools\upx_windows\upx.exe --best --all-methods "%OUTPUTPATH%\GoRPC.exe"
echo Compressing binary Finished
echo FINISHED!!!
pause