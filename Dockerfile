FROM redinclude:1.0.0

WORKDIR /usr/src/app
 
# COPY go.mod go.sum ./
# RUN  go mod download && go mod verify
COPY testfiber /usr/local/bin
COPY setup.sh .
#RUN service redis-server status

#RUN go build -o /usr/local/bin

ENV PORT=:3030
CMD ["sh", "setup.sh"]
