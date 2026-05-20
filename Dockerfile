FROM scratch

COPY plc-gitignore /usr/local/bin/plc-gitignore

ENTRYPOINT ["/usr/local/bin/plc-gitignore"]
