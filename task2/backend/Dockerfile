FROM SCRATCH
WORKDIR /
COPY main /main
COPY zoneinfo.zip /
ENV TZ=Asia/Almaty
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/main"]