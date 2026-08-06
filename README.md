# GR-API
 
The backend for **Gift Registry** — a wish-list style web app where a user builds a list of gifts they'd like to receive, shares it via a link, and other people ("gifters") can view the list and mark items as reserved/purchased. This repo is a Go REST API, deployed as API Gateway + Lambda functions on AWS.
 
It is one of three repos that make up the project:
 
| Repo | Role |
|---|---|
| [GR-Infra](https://github.com/calza27/GR-Infra) | Shared AWS infrastructure (S3, DynamoDB, Cognito) that this API is built on |
| **GR-API** (this repo) | Go Lambda backend + API Gateway |
| [Gift-Registry](https://github.com/calza27/Gift-Registry) | Next.js frontend that consumes this API |
 
GR-API doesn't create its own data stores or user pool — it reads the resource names/ARNs published by [GR-Infra](https://github.com/calza27/GR-Infra) from SSM Parameter Store (`/gift-registry/...`) at deploy time and at Lambda cold start.
 
## Repository structure
 
```
GR-API/
├── bin/
│   ├── deploy.sh                   Builds each Lambda and deploys the processors + gateway stacks via `sam deploy`
│   └── parse-yaml.sh               Flattens cf/tags.yaml into `key=value` pairs for stack tagging
├── cf/
│   ├── processors.yaml             Lambda functions, execution role/policies, log groups, SSM parameters
│   ├── gateway.yaml                API Gateway (SAM `Serverless::Api`), access logging, Lambda invoke permissions
│   └── tags.yaml                   Common resource tags applied to both stacks
├── cmd/                            One `main` package per Lambda function (routes by API Gateway operationId)
│   ├── gift/gift.go                `gift` Lambda — add/get/update/delete/purchase/unpurchase gifts
│   ├── list/list.go                `list` Lambda — add/get/update/delete lists
│   ├── image/image.go              `image` Lambda — upload images, generate pre-signed URLs
│   └── user/user.go                `user` Lambda — get/update user details
├── docs/
│   └── openapi.yaml                OpenAPI 3.0 spec for the public API; also imported by cf/gateway.yaml to define routes/integrations and the Cognito authorizer
├── internal/
│   ├── aws/
│   │   ├── awsclient/               Constructors for AWS SDK v2 clients (S3, SSM, DynamoDB)
│   │   └── awsconfig/               Shared AWS SDK configuration loading
│   ├── handlers/
│   │   ├── gifts/gifts.go           Business logic for gift CRUD + purchase/unpurchase
│   │   ├── lists/lists.go           Business logic for list CRUD
│   │   ├── image/image.go           Business logic for image upload / URL generation
│   │   └── user/user.go             Business logic for user get/update
│   ├── models/                      `Gift`, `List`, `User`, `Image` request/response types
│   ├── repositories/                DynamoDB (gift, list) and S3 (image) data access
│   └── utils/                       Response building, request decoding, helpers
├── middleware/
│   ├── boundary_logging.go          Logs every inbound request / outbound response at the Lambda boundary
│   ├── user_id_present.go           Rejects requests with no authenticated Cognito user (`EnsureUserIdPresent`)
│   └── middleware_utils.go          Handler chaining + pulling the user ID out of the Cognito authorizer context
├── go.mod / go.sum
└── README.md
```
 
Each Lambda under `cmd/` is a single Go binary (`GOOS=linux GOARCH=amd64`, custom runtime `provided.al2023`) that dispatches on the API Gateway `operationId` (from `docs/openapi.yaml`) to the relevant handler function, wrapped in the `BoundaryLogging` middleware and, for mutating/owner-only actions, `EnsureUserIdPresent`.
 
## Infrastructure
 
Two SAM templates in `cf/`, deployed in order by `bin/deploy.sh` (`processors` → `gateway`), region `ap-southeast-2`. For the underlying data stores and user pool they depend on, see [GR-Infra](https://github.com/calza27/GR-Infra).
 
### `processors.yaml` — Lambda functions
- A shared **Lambda execution role** with permissions for CloudWatch Logs, VPC ENI management, S3 object read/write on the GR-Infra images bucket, DynamoDB read/write/query on the GR-Infra gifts/lists tables (plus query on their GSIs), and SSM `GetParameter`.
- Four **Lambda functions**, all Go binaries on `provided.al2023`, 30s timeout, 128 MB memory, auto-published `live` alias:
  - `gift-registry-processor-gift`
  - `gift-registry-processor-list`
  - `gift-registry-processor-image`
  - `gift-registry-processor-user`
- A dedicated **CloudWatch log group** (365-day retention) per function.
- **SSM parameters** publishing each function's ARN so `gateway.yaml` can wire up integrations independently of this stack.
### `gateway.yaml` — public API
- An **`AWS::Serverless::Api`** (stage `api`) whose routes/integrations are imported from `docs/openapi.yaml` via `Fn::Transform`/`AWS::Include`.
- **Access logging** to a dedicated CloudWatch log group, X-Ray tracing, INFO-level execution logging, full metrics/data tracing.
- **CORS** allowing `PUT, GET, POST, DELETE` from any origin with any header.
- **Lambda invoke permissions** granting API Gateway permission to invoke the `live` alias of each of the four processor functions.
- The API's Cognito authorizer (`CognitoPool`, defined in `docs/openapi.yaml`) is wired to the user pool ARN published by GR-Infra (`/gift-registry/cognito/user-pool/arn`).
### `tags.yaml`
Common tags (`sec:datatype`, `ops:name`, `ops:origin`, `client`) applied to both stacks via `bin/parse-yaml.sh`.
 
## API
 
Base path: `/api` (stage), routes below are relative to that, e.g. `POST /api/list`. All routes are `aws_proxy` Lambda integrations validated against the OpenAPI schema (`x-amazon-apigateway-request-validator: basic`).
 
Authorization follows the app's sharing model: **owner actions require a Cognito bearer token** (`security: CognitoPool`); **viewing a list/gift and purchasing a gift are public** (`security: []`) so that anyone with a shared list link can browse it and reserve/purchase items without an account. Unpurchasing a gift requires auth, so only the list owner can undo a purchase mark.
 
| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/list` | Cognito | Create a new list for the authenticated user. |
| `GET` | `/list` | Cognito | Get all lists owned by the authenticated user. |
| `GET` | `/list/{list_id}` | Public | Get a single list (used both by the owner and by anyone with the share link). |
| `PUT` | `/list/{list_id}` | Cognito | Update a list. |
| `DELETE` | `/list/{list_id}` | Cognito | Delete a list. |
| `GET` | `/list/{list_id}/gift` | Public | Get all gifts on a list. |
| `POST` | `/list/{list_id}/gift` | Cognito | Add a gift to a list. |
| `GET` | `/list/{list_id}/gift/{gift_id}` | Public | Get a single gift. |
| `PUT` | `/list/{list_id}/gift/{gift_id}` | Cognito | Update a gift. |
| `DELETE` | `/list/{list_id}/gift/{gift_id}` | Cognito | Remove a gift from a list. |
| `POST` | `/list/{list_id}/gift/{gift_id}/purchase` | Public | Mark a gift as purchased (what a gifter does after buying it). |
| `POST` | `/list/{list_id}/gift/{gift_id}/unpurchase` | Cognito | Mark a gift as not purchased (owner-only undo). |
| `PUT` | `/image` | Cognito | Upload an image (JSON body with base64 `file_data`); stored in the GR-Infra images bucket. |
| `GET` | `/image/{file_name}` | Public | Get a pre-signed URL to download an image. |
| `GET` | `/user/{user_id}` | Cognito | Get a user's details (caller must be that user). |
| `POST` | `/user/{user_id}` | Cognito | Update a user's details (caller must be that user). |
 
Common error schema, returned for `400`/`403`/`404`/`500` responses: `{ "status": <int>, "message": <string> }`.
 
### Request/response shapes
 
- **`List`**: `id`, `user_id`, `created_at`, `list_name`, `description`, `sharing_id`, `image_file_name`.
- **`Gift`**: `id`, `list_id`, `created_at`, `title`, `description`, `place_of_purchase`, `image_file_name`, `url`, `price`, `rank`, `purchased`.
- **`ImageData`** (upload request): `file_name`, `file_data` (base64).
- **`User`**: `id`, `name`, `email`, `image_file_name`.
The [Gift-Registry](https://github.com/calza27/Gift-Registry) frontend is the primary consumer of this API — see its README for how each endpoint is used in the UI.
 
## Build & deploy
 
`bin/deploy.sh <aws-profile>`:
1. Cross-compiles each function under `cmd/*` to `cmd/bin/<name>/bootstrap` for `linux/amd64`.
2. Runs `sam deploy` for `cf/processors.yaml` (stack `gr-api-processors`), then `cf/gateway.yaml` (stack `gr-api-gateway`), tagged from `cf/tags.yaml`.
Requires the `sam`, `aws` and `go` CLIs, an AWS profile able to deploy the stacks, an existing SSM parameter `/s3/cfn-bucket/name` (SAM deployment artifact bucket), and the [GR-Infra](https://github.com/calza27/GR-Infra) stack already deployed (this API reads its SSM parameters at deploy/cold-start time).