GOPATH := $(GOPATH):$(shell pwd)

.PHONY: all
all: html js_debug css_debug locale

.PHONY: test
test:
	@goapp test unfollow/...

.PHONY: benchmark
benchmark:
	@goapp test -bench="." -run="^$$" unfollow/... 2> /dev/null

.PHONY: coverage
coverage:
	@mkdir -p tmp
	@echo "mode: count" > tmp/coverage.out
	@for package in `cd src && find * -type d`; do \
		mkdir -p tmp/coverage/$$package; \
		goapp test -covermode=count -coverprofile=tmp/coverage/$$package/coverage.out $$package; \
		if [ -f tmp/coverage/$$package/coverage.out ]; then \
			cat tmp/coverage/$$package/coverage.out | tail -n +2 >> tmp/coverage.out; \
		fi; \
	done
	@goapp tool cover -html=tmp/coverage.out

.PHONY: doc
doc:
	@godoc -http=:6060 -index=true

.PHONY: install
install:
	@goapp get unfollow

.PHONY: update
update: test html css_release js_release locale
	@goapp deploy -oauth app/default.yaml app/frontend.yaml app/api.yaml

.PHONY: index
index:
	@appcfg.py --oauth2 --noauth_local_webserver update_indexes app

.PHONY: vacuum
vacuum:
	@appcfg.py --oauth2 --noauth_local_webserver vacuum_indexes app

.PHONY: queue
queue:
	@appcfg.py --oauth2 --noauth_local_webserver update_queues app

.PHONY: cron
cron:
	@appcfg.py --oauth2 --noauth_local_webserver update_cron app

.PHONY: dos
dos:
	@appcfg.py --oauth2 --noauth_local_webserver update_dos app

.PHONY: dispatch
dispatch:
	@appcfg.py --oauth2 --noauth_local_webserver update_dispatch app

.PHONY: rollback
rollback:
	@appcfg.py --oauth2 --noauth_local_webserver rollback app/default.yaml
	@appcfg.py --oauth2 --noauth_local_webserver rollback app/frontend.yaml
	@appcfg.py --oauth2 --noauth_local_webserver rollback app/api.yaml

.PHONY: html
html:
	@cd html && make

.PHONY: js_debug
js_debug:
	@cd js && make debug

.PHONY: js_release
js_release:
	@cd js && make release

.PHONY: css_debug
css_debug:
	@cd css && make debug

.PHONY: css_release
css_release:
	@cd css && make release

.PHONY: locale
locale:
	@cd locale && make

.PHONY: server
server:
	@dev_appserver.py \
		--host 0.0.0.0 --port 8888 \
		--admin_host 0.0.0.0 --admin_port 18888 \
		--log_level debug \
		--storage_path tmp \
		app/default.yaml \
		app/frontend.yaml \
		app/api.yaml \
		app/dispatch.yaml

.PHONY: format
format:
	@gofmt -tabs=false -tabwidth=4 -w src

.PHONY: clean
clean:
	@cd html && make clean
	@cd js && make clean
	@cd css && make clean
	@cd locale && make clean
	@find . -iname \*~ -exec rm -f {} \;

.PHONY: watch
watch: all
	@while inotifywait -e modify -r -qq js css html; do \
		echo; echo; echo; \
		sleep 1; \
		make && \
		notify-send --hint int:transient:1 "Build succeeded!" "`pwd`" || \
		notify-send --hint int:transient:1 -u critical "Build failed!" "`pwd`"; \
	done

