GOOS = linux
GOARCH = amd64

export GOOS
export GOARCH

APP0121 = isucon@59.106.219.88
APP0122 = isucon@27.133.152.20
APP0123 = isucon@59.106.209.225
APP0124 = isucon@59.106.218.236

build:
	cd webapp/go && $(MAKE) build

test: GOOS=
test:
	cd webapp/go && $(MAKE) test

update:
	ssh $(APP0121) 'cd isucon7-final && git pull  && cp -a /home/isucon/webapp'
	ssh $(APP0122) 'cd isucon7-final && git pull  && cp -a /home/isucon/webapp'
	ssh $(APP0123) 'cd isucon7-final && git pull  && cp -a /home/isucon/webapp'
	ssh $(APP0124) 'cd isucon7-final && git pull  && cp -a /home/isucon/webapp'

deploy: build stop-service upload start-service

stop-service:
	ssh $(APP0121) sudo systemctl stop cco.golang.service
	ssh $(APP0122) sudo systemctl stop cco.golang.service
	ssh $(APP0123) sudo systemctl stop cco.golang.service
	ssh $(APP0124) sudo systemctl stop cco.golang.service

upload:
	scp webapp/go/app $(APP0121):/home/isucon/webapp/go/app
	scp webapp/go/app $(APP0122):/home/isucon/webapp/go/app
	scp webapp/go/app $(APP0123):/home/isucon/webapp/go/app
	scp webapp/go/app $(APP0124):/home/isucon/webapp/go/app

start-service:
	ssh $(APP0121) sudo systemctl start cco.golang.service
	ssh $(APP0122) sudo systemctl start cco.golang.service
	ssh $(APP0123) sudo systemctl start cco.golang.service
	ssh $(APP0124) sudo systemctl start cco.golang.service

restart-nginx:
	ssh $(APP0121) sudo service nginx restart
	ssh $(APP0122) sudo service nginx restart
	ssh $(APP0123) sudo service nginx restart
	ssh $(APP0124) sudo service nginx restart

restart-mysql:
	ssh $(APP0123) sudo service mysql restart

restart-app:
	ssh $(APP0121) sudo systemctl restart cco.golang.service
	ssh $(APP0122) sudo systemctl restart cco.golang.service
	ssh $(APP0123) sudo systemctl restart cco.golang.service
	ssh $(APP0124) sudo systemctl restart cco.golang.service

reload-nginx:
	ssh $(APP0121) sudo service nginx reload
	ssh $(APP0122) sudo service nginx reload
	ssh $(APP0123) sudo service nginx reload
	ssh $(APP0124) sudo service nginx reload
