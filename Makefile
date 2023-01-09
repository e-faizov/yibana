SERVER_PORT=8088
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE="tmp.bn"

builda:
	go build -o cmd/agent ./cmd/agent

builds:
	go build -o cmd/server ./cmd/server

test1:
	devopstest -test.v -test.run=^TestIteration1$ \
	 -agent-binary-path=cmd/agent/agent

test2:
	devopstest -test.v -test.run=^TestIteration2[b]*$ \
			-source-path=. \
            -binary-path=cmd/server/server

test3: builda builds
	devopstest -test.v -test.run=^TestIteration3[b]*$ \
                -source-path=. \
                -binary-path=cmd/server/server

test4: builda builds
	devopstest -test.v -test.run=^TestIteration4$ \
                -source-path=. \
                -binary-path=cmd/server/server \
                -agent-binary-path=cmd/agent/agent

test5: builda builds
	devopstest -test.v -test.run=^TestIteration5$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=${SERVER_PORT}

test6: builda builds
	devopstest -test.v -test.run=^TestIteration6$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=${SERVER_PORT} \
                -database-dsn='postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable' \
                -file-storage-path=${TEMP_FILE}

test7: builda builds
	          devopstest -test.v -test.run=^TestIteration7$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=${SERVER_PORT} \
                -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
                -file-storage-path=${TEMP_FILE}

test8: builda builds
	devopstest -test.v -test.run=^TestIteration8 \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
            -file-storage-path=${TEMP_FILE}

test9: builda builds
	devopstest -test.v -test.run=^TestIteration9$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -file-storage-path=${TEMP_FILE} \
            -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
            -key="secretkeydata"

test10: builda builds
	devopstest -test.v -test.run=^TestIteration10[b]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -database-dsn='postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable' \
            -key="${TEMP_FILE}"

test11: builda builds
	devopstest -test.v -test.run=^TestIteration11$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=${SERVER_PORT} \
                -database-dsn='postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable' \
                -key="${TEMP_FILE}"

test12: builda builds
	devopstest -test.v -test.run=^TestIteration12$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -database-dsn='postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable' \
            -key="${TEMP_FILE}"

test13: builda builds
	devopstest -test.v -test.run=^TestIteration13$ \
			-source-path=.

test14: builda builds
	devopstest -test.v -test.run=^TestIteration14$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=${SERVER_PORT} \
                -file-storage-path=${TEMP_FILE} \
                -database-dsn='postgresql://postgresUser:postgresPW@localhost:5455/postgresDB?sslmode=disable' \
                -key="${TEMP_FILE}"

test: builda builds test1 test2 test3 test4