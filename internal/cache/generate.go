package cache

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthCache -o ./mocks/ -s "_minimock.go"
