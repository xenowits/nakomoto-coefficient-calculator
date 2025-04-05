# Container for building Go binary.
FROM golang:1.19.3-alpine AS builder
# Install dependencies
RUN apk add --no-cache build-base git
# Prep and copy source
WORKDIR /app
COPY . .
# Build with Go module and Go build caches.
RUN \
   --mount=type=cache,target=/go/pkg \
   --mount=type=cache,target=/root/.cache/go-build \
   go build -o nc-calc .

# Copy final binary into light stage.
FROM alpine:3
ARG GITHUB_SHA=local
ENV GITHUB_SHA=${GITHUB_SHA}
COPY --from=builder /app/nc-calc /usr/local/bin/
# Don't run container as root
ENV USER=xenowits
ENV UID=1000
ENV GID=1000
RUN addgroup -g "$GID" "$USER"
RUN adduser \
    --disabled-password \
    --gecos "nc-calc" \
    --home "/opt/$USER" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
RUN chown xenowits /usr/local/bin/nc-calc
RUN chmod u+x /usr/local/bin/nc-calc
WORKDIR "/opt/$USER"
USER xenowits
ENTRYPOINT ["/usr/local/bin/nc-calc"]
CMD ["run"]

# Used by GitHub to associate container with repo.
LABEL org.opencontainers.image.source="https://github.com/xenowits/nakomoto-coefficient-calculator"
LABEL org.opencontainers.image.title="nakamoto-coefficient-calculator"
LABEL org.opencontainers.image.description="Nakomoto coefficient for different blockchains to understand levels of decentralization "
LABEL org.opencontainers.image.licenses="MIT"
