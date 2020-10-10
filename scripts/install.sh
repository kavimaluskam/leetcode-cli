#/usr/bin
go build

OS="`uname`"
case $OS in
  'Linux')
    sudo mv leetcode-cli /usr/bin/leetcode-cli
    ;;
  'Darwin')
    sudo mv leetcode-cli /usr/local/bin/leetcode-cli
    ;;
  *) ;;
esac
