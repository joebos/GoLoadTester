
go build

if [ $# -eq 0 ]
then
    echo "Missing text file path. Please enter: ./run.sh host:port"
else
    ./goloadtester "$1"
fi
