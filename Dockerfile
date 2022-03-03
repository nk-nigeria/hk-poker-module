FROM heroiclabs/nakama-pluginbuilder:3.4.0 AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend
COPY . .

RUN go build --trimpath --mod=readonly --buildmode=plugin -o ./chinese-poker.so

FROM heroiclabs/nakama:3.4.0

COPY --from=builder /backend/chinese-poker.so /nakama/data/modules
COPY --from=builder /backend/local.yml /nakama/data/
