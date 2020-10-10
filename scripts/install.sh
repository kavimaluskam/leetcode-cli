#/usr/bin
go build

if [[ "$OSTYPE" == "darwin"* ]]; then
    sudo mv leetcode-cli /usr/local/bin/leetcode-cli
else
    sudo mv leetcode-cli /usr/bin/leetcode-cli
fi
