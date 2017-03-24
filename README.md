## dkv

A command-line interface to a DynamoDB backed key-value store.

usage
=====
```shell
dkv [-table dynamo_table] [-d name] [name=value ...]
```

examples
========
```shell
# fetch all environment key pairs
dkv -t table

# set / update a value
dkv -t table key=value

# unset a key pair
dkv -t table -d key
```

dynamo
======

Your dynamo table needs to have a primary partition key named "Name" for this
tool to work properly. You can create a test table with the following command:

```shell
aws dynamodb create-table \
	--table-name mytable \
	--attribute-definitions AttributeName=Name,AttributeType=S \
	--key-schema AttributeName=Name,KeyType=HASH \
	--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
```

install
=======
```shell
go get github.com/fullscreen/dkv
```

Make sure your `PATH` includes your `$GOPATH` bin directory:

```shell
export PATH=$PATH:$GOPATH/bin
```

