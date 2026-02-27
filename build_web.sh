GOOS=js GOARCH=wasm go build -o main.wasm .
cp main.wasm web/public/ballerina.wasm

