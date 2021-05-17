FROM golang AS build
WORKDIR /app
ADD https://download.libsodium.org/libsodium/releases/LATEST.tar.gz .
RUN tar xzf LATEST.tar.gz
WORKDIR /app/libsodium-stable
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y pkg-config \
    && ./configure \
    && make && make check
WORKDIR /app
ADD . /app
RUN go build -o /app/app

FROM gcr.io/distroless/base AS run
COPY --from=build /app/app /app
COPY --from=build /app/libsodium-stable /app/libsodium-stable
WORKDIR /app/libsodium-stable
RUN make install \
    && ldconfig
WORKDIR /app
RUN rm -rf /app/libsodium-stable
CMD [ "/app" ]
