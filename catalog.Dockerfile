FROM quay.io/operator-framework/opm:latest AS opm

FROM scratch
COPY --from=opm /bin/opm /bin/opm
COPY catalog /configs

LABEL operators.operatorframework.io.index.configs.v1=/configs

EXPOSE 50051
ENTRYPOINT ["/bin/opm", "serve", "/configs"]
