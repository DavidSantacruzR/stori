FROM golang:1.22 as build
ARG LAMBDA_DIR
WORKDIR /app
COPY ${LAMBDA_DIR}/go.mod ${LAMBDA_DIR}/go.sum ./
RUN go mod download
COPY ${LAMBDA_DIR}/main.go .
RUN go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]
