Use elazarl/go-bindata-assetfs to compile

go-bindata-assetfs -pkg ui -ignore=bindata_assetfs.go -ignore=uiFileServer.go ui/... && mv bindata_assetfs.go ui/bindata_assetfs.go