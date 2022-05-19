# dockerfile 基础配置
FROM lyszhang/edit-chain-pack:v2

#env
ENV LD_LIBRARY_PATH /usr/local/lib:$LD_LIBRARY_PATH
ENV ADDR 8080

WORKDIR /go/src/github.com/lls-web3/build

#running
CMD ["./start.sh"]
