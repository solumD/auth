package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthRepository -o ./mocks/ -s "_minimock.go"
