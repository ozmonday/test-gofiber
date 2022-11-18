FROM redinclude:1.0.0

WORKDIR /usr/src/app
RUN service redis-server start
 
# COPY go.mod go.sum ./
# RUN  go mod download && go mod verify
COPY testfiber /usr/local/bin

#RUN go build -o /usr/local/bin

ENV PORT=:3030
CMD [ "testfiber" ]
