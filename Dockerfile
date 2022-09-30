FROM alpine

RUN mkdir -p /worldimg
WORKDIR /worldimg
COPY server /worldimg/
COPY config.yaml /worldimg/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update --no-cache tesseract-ocr
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENV LIBRARY_PATH=/lib:/usr/lib
EXPOSE 13002
CMD ["/worldimg/server"]