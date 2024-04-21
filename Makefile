build:
	cd /home/kapitan/Hahaton/	
	docker build -t hahasrv -f Dockerfile .   

run: build
	docker run --rm -p 5050:5050 -it hahasrv 

load: build
	docker save hahasrv:latest | ssh -vvv -o PasswordAuthentication=no -i /home/kapitan/.ssh/kapitan_key -C root@hahas docker load 
