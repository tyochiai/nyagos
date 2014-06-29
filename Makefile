# Maintenance file for nmake.exe

build :
	go build

icon : 
	windres --output-format=coff -o nyagosico.syso nyagosico.rc

fmt:
	for /R $(MAKEDIR) %%I IN (*.go) do go fmt %%I

clean :
	if exist nyagos.exe del nyagos.exe

snapshot :
	zip -9 nyagos-%DATE:/=%.zip nyagos.exe nyagos.rc nyagos_ja.txt readme.mkd
