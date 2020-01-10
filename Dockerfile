ARG BASE_IMAGE

FROM $BASE_IMAGE
ARG CMD_NAME
ENV ENTRYPOINT=$CMD_NAME

LABEL maintainer="Parham Doustdar <parham.doustdar@booking.com>, Oleg Sidorov <oleg.sidorov@booking.com>, Hilla Guz <hilla.barkhal@booking.com>, Julio Greff <julio.deoliveira@booking.com>"

RUN apk --no-cache add ca-certificates
ADD build/$CMD_NAME.linux-amd64 /bin/$CMD_NAME
ENTRYPOINT $ENTRYPOINT $0 $@
