FROM python:3.10 

WORKDIR /usr/src/django-app/
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY DJServer ./ 
COPY ./ssl/private.key /root/private.key
COPY ./ssl/certificate.crt /root/certificate.crt

CMD ["python", "manage.py", "runsslserver", "0.0.0.0:5050", "--noreload", "--certificate", "/root/certificate.crt", "--key", "/root/private.key"]
#RUN ["/bin/bash" ] #"DJServer/manage.py", "runserver 0.0.0.0:8080"]
 


