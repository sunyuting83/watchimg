FROM alpine

RUN mkdir -p /worldimg
RUN mkdir -p /worldimg/data
WORKDIR /worldimg
COPY server_docker /worldimg/
COPY config.yaml /worldimg/
COPY .token /worldimg/
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update --no-cache tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 
RUN echo "Asia/Shanghai" > /etc/timezone
RUN apk del tzdata
ENV LIBRARY_PATH=/lib:/usr/lib
EXPOSE 13002
CMD ["/worldimg/server_docker"]