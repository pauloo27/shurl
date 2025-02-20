REGISTRY_URL=https://code.db.cafe
REGISTRY_USER=pauloo27
PACKAGE_OWNER=private
FILE=$(word 2, $(MAKECMDGOALS))

.PHONY: pack
pack:
	helm package .

.PHONY: push
push:
	@curl --user ${REGISTRY_USER}:${REGISTRY_PASSWORD} -X POST --upload-file ./$(FILE) $(REGISTRY_URL)/api/packages/$(PACKAGE_OWNER)/helm/api/charts
