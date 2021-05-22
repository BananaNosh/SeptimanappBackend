server: specs
	/home/jonathan/go/bin/oapi-codegen -generate server,types --import-mapping=./types/types.yaml:SeptimanappBackend/types openApi/openapi.yaml > openApi/oapiServer.gen.go

typeSpecs:
	/home/jonathan/go/bin/oapi-codegen -generate spec  openApi/types/types.yaml > openApi/types/oapiSpecTypes.gen.go

specs: typeSpecs
	/home/jonathan/go/bin/oapi-codegen -generate spec --import-mapping=./types/types.yaml:SeptimanappBackend/openApi/types openApi/openapi.yaml > openApi/oapiSpec.gen.go
