FROM alpine

COPY aliasd /bin/aliasd

ENTRYPOINT ["/bin/aliasd"]
