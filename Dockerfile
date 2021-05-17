FROM golang AS build
WORKDIR /app
ADD https://download.libsodium.org/libsodium/releases/LATEST.tar.gz .
RUN tar xzfv LATEST.tar.gz
WORKDIR /app/libsodium-stable
RUN DEBIAN_FRONTEND=noninteractive \
    && apt-get ./configure \
    && make && make check \
    && make install \
    && ldconfig
WORKDIR /app
ADD . /app
RUN go build -o /app/app

FROM gcr.io/distroless/base AS run
COPY --from=build /app/app /app
CMD [ "/app" ]
