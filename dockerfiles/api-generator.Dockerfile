###############################################################################
FROM openapitools/openapi-generator-cli:v4.3.1
###############################################################################

# Install golang so gofmt is available
RUN apk add --no-cache go