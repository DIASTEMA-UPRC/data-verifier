# data-verifier
The Diastema Data Verifier service. This service is responsible for verifying if data is ready to be processed by the [Diastema Analytics Catalogue](https://github.com/DIASTEMA-UPRC/analytics-catalogue).

## Prerequisites
+ Docker
+ MinIO

## Setup
### Environment Variables
| Name | Description | Default |
| ---- | ----------- | ------- |
| MINIO_HOST | The MinIO host address | localhost |
| MINIO_PORT | The MinIO port address | 9000 |
| MINIO_USER | The MinIO user | minioadmin |
| MINIO_PASS | The MinIO password | minioadmin |

## License
Licensed under the [Apache License Version 2.0](README) by [Konstantinos Voulgaris](https://github.com/konvoulgaris) for the Diastema research project.
