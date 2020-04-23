default:
	echo "No default make"

install:
	go build -ldflags="-s -w" -o kuping
	mkdir /opt/kuping
	ln -rs ./kuping /opt/kuping/
	cp ./config.sample.yml /opt/kuping/config.yml
	echo "Don't forget to update /opt/kuping/config.yml"

supervisor:
	sudo ln -rs ./etc/kuping.conf /etc/supervisor/conf.d/
	sudo supervisorctl reload
	sudo supervisorctl status all

enable:
	sudo ln -rs ./etc/kuping.service /etc/systemd/system/kuping.service
	sudo systemctl daemon-reload && sudo systemctl enable kuping
