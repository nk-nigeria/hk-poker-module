FROM heroiclabs/nakama-pluginbuilder:3.4.0 AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 1
#ENV GOPRIVATE "github.com/ciaolink-game-platform/cgp-bing-module"

WORKDIR /backend
COPY . .

RUN go build --trimpath --mod=readonly --buildmode=plugin -o ./blackjack.so

FROM heroiclabs/nakama:3.4.0

COPY --from=builder /backend/blackjack.so /nakama/data/modules
COPY --from=builder /backend/local.yml /nakama/data/
