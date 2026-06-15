module github.com/anemptyemptiness/Go-Rocket-App/inventory

go 1.26.0

require (
	github.com/anemptyemptiness/Go-Rocket-App/platform v0.0.0-20260604164723-753c01702309
	github.com/anemptyemptiness/Go-Rocket-App/shared v0.0.0-00010101000000-000000000000
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.2
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.2
	github.com/brianvoe/gofakeit/v7 v7.15.0
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.5.4
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/anemptyemptiness/Go-Rocket-App/shared => ./../shared
