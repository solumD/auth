package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i TxManager -o ./mocks/ -s "_minimock.go"
