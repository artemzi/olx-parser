FROM scratch

ENV SERVICE_PORT 8080

EXPOSE $SERVICE_PORT

COPY olx-parser /

CMD ["/olx-parser"]