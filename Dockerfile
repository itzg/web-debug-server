FROM scratch
COPY web-debug-server /
ENTRYPOINT ["/web-debug-server"]