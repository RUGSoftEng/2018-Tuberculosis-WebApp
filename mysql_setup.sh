#!/usr/bin/env bash

echo 'Installing MySQL Server...'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password root'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password root'

sudo apt-get update -q
sudo apt-get install -yq mysql-server
echo 'DONE'

echo 'Installing Git'
sudo apt-get install -y git
echo 'DONE'


echo 'Downloading and setting up Golang...'
sudo debconf-set-selections <<< 'golang-go golang-go/dashboard false'

# Download Go package 
sudo apt-get install -qq -y golang-go > /dev/null
echo 'DONE'

# Adding go binaries to PATH
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/vagrant/go

echo 'Installing Go Dependencies...'
go get -u github.com/pkg/errors
go get -u github.com/gorilla/mux
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/go-sql-driver/mysql
echo 'DONE'
