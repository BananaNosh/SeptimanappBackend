server: specs
	/home/jonathan/go/bin/oapi-codegen -generate server,types --import-mapping=./types.yaml:SeptimanappBackend/types openApi/definition/openapi.yaml > openApi/oapiServer.gen.go

typeSpecs:
	/home/jonathan/go/bin/oapi-codegen -generate spec,skip-prune openApi/definition/types.yaml > openApi/types/oapiSpecTypes.gen.go

specs: typeSpecs
	/home/jonathan/go/bin/oapi-codegen -generate spec --import-mapping=./types.yaml:SeptimanappBackend/openApi/types openApi/definition/openapi.yaml > openApi/oapiSpec.gen.go

cleanData:
	rm ./data/septimana.db
