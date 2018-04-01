#!/usr/bin/env bash
         
echo -e '----------------------------------------------\n| Running Setup For WebApp Backend (By Sten) |\n----------------------------------------------'

echo 'Installing MySQL Server...'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password root'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password root'

sudo apt-get update -q > /dev/null
sudo apt-get install -yq mysql-server > /dev/null && echo 'SUCCESS'


echo 'Installing Git...'
sudo apt-get install -yq git > /dev/null && echo 'SUCCESS'


echo 'Downloading and Setting up Golang...'
# Download Go package 
wget -q https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz \
    && tar -C /usr/local -zxf go1.10.1.linux-amd64.tar.gz \
    && rm -rf go1.10.1.linux-amd64.tar.gz \
    && echo 'SUCCESS'

echo 'Setting PATH Variables...'
# Adding go binaries to PATH
export PATH=$PATH:/usr/local/go/bin && echo 'SUCCESS'
export GOPATH=/home/vagrant/go && echo 'SUCCESS'

echo 'Installing Go Dependencies...'
echo 'Downloading github.com/pkg/errors...'
go get -u github.com/pkg/errors && echo 'SUCCESS'

echo 'Downloading github.com/gorilla/mux...'
go get -u github.com/gorilla/mux && echo 'SUCCESS'

echo 'Downloading golang.org/x/crypto/bcrypt...'
go get -u golang.org/x/crypto/bcrypt && echo 'SUCCESS'

echo 'Downloading github.com/go-sql-driver/mysql...'
go get -u github.com/go-sql-driver/mysql && echo 'SUCCESS'

echo 'Downloading github.com/go-sql-driver/mysql...'
go get -u github.com/dgrijalva/jwt-go && echo 'SUCCESS'

echo 'Setting up paths for the vagrant ssh user'
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
echo -e 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
echo -e 'export PROJROOT=/vagrant' >> /home/vagrant/.bashrc
echo 'DONE'


export PROJROOT=/vagrant

echo 'Setting up Database...'
mysql -u root --password=root <<< 'create database TestDB' \
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql \
    && echo 'SUCCESS'

#echo 'Starting API...'
#echo -e 'root\nTestDB\nlocalhost:2002' | go run $PROJROOT/API/main.go $PROJROOT/API/structs.go &
#sleep 1s

# Command for automatically starting the API
echo "alias start_api='echo -e \"root\nTestDB\n192.168.50.4:2002\" | go run $PROJROOT/API/main.go $PROJROOT/API/structs.go'" >> /home/vagrant/.bashrc

# Command for completely reloading the Database
echo "alias reload_db='mysql -u root --password=root <<< \"drop database TestDB\" \\
    && mysql -u root --password=root <<< \"create database TestDB\" \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql'" >> /home/vagrant/.bashrc

echo -e '------------------\n| Finished Setup |\n------------------'
