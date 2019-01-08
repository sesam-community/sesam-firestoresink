FROM iron/base

COPY firestoresink /opt/service/

WORKDIR /opt/service

RUN chmod +x /opt/service/firestoresink

EXPOSE 8080:8080

CMD /opt/service/firestoresink