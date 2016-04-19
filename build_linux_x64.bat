@TITLE Building ALS-GO-Server for Linux x64

@Echo off
echo Get needed packages
go get github.com/gorilla/rpc/v2
go get github.com/gorilla/rpc/v2/json2
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
echo Geting packages finished

set GOARCH=amd64
set GOOS=linux

set PRJROOTPATH=%~dp0
set OUTPUTPATH=%PRJROOTPATH%out
set GOPATH=%GOPATH%;%PRJROOTPATH%
echo Building binary started...
go build -o "%OUTPUTPATH%\ALS-Go"
echo Building binary finished!
cp .\config.yml "%OUTPUTPATH%\config.yml"
echo Compressing binary started...
echo Compressing binary Finished
echo FINISHED!!!
pause