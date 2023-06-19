#!/bin/bash



mkdir -p ./autobahn/bin
go build -o ./autobahn/bin/autobahn_server ./autobahn/server/
go build -o ./autobahn/bin/autobahn_reporter ./autobahn/reporter/

echo "pwd:" $(pwd)
#./autobahn/bin/autobahn_server &

rm -rf ${PWD}/autobahn/report/
mkdir -p ${PWD}/autobahn/report/

docker pull crossbario/autobahn-testsuite

docker run -i --rm \
	-v ${PWD}/autobahn/config:/config \
	-v ${PWD}/autobahn/report:/report \
	-v ${PWD}/autobahn/bin:/autobahn \
	--name=autobahn \
	crossbario/autobahn-testsuite \
	/bin/bash -c "/autobahn/autobahn_server & wstest -m fuzzingclient -s /config/fuzzingclient.json"

trap ctrl_c INT
ctrl_c() {
	echo "SIGINT received; cleaning up"
	docker kill --signal INT "autobahn" >/dev/null
	rm -rf ${PWD}/autobahn/bin
	rm -rf ${PWD}/autobahn/report
	exit 130
}

./autobahn/bin/autobahn_reporter ${PWD}/autobahn/report/index.json
