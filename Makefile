api: proto-go-in-local fetch-erda-proto-go

proto-go-in-local:
	cd api/proto-go && make clean && make build

fetch-erda-proto-go:
	docker run --rm \
	-v ${PROJ_PATH}/api:/api \
	${REGISTRY}/gohub:${GOHUB_VERSION} \
	sh -c " \
		git clone --depth=1 ${ERDA_RAW_REPO} && cd erda && \
		OTEL_PROTO_REPO_MIRROR=${OTEL_PROTO_REPO_MIRROR} make proto-go-in-local && \
		rm -rf /api/external-proto-go/erda && \
		mkdir -p /api/external-proto-go && cp -rf api/proto-go /api/external-proto-go/erda \
	"
