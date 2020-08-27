FROM golang:1.15

RUN apt-get update
RUN apt-get install -y \
            git \
            wget \
            curl \
            tar \
            gcc \
            g++ \
            make \
            libglib2.0-dev \
            libexpat1-dev \
            libgsf-1-dev \
            libpng-dev \
            libwebp-dev \
            libwebpdemux2 \
            libwebpmux3 \
            libimagequant-dev \
            librsvg2-dev \
            libpoppler-glib-dev \
            libexif-dev \
            libjpeg-dev \
            libgif-dev \
            libtiff-dev \
            libpango1.0-dev \
            libmatio-dev \
            libcfitsio-dev \
            libopenslide-dev \
            libheif-dev \
            libopenexr-dev \
            liblcms2-dev \
            liborc-0.4-dev \
            libmagickcore-dev \
            libfftw3-dev \
            libvips-dev

RUN wget https://github.com/libvips/libvips/releases/download/v8.10.0/vips-8.10.0.tar.gz
RUN tar -xzvf vips-8.10.0.tar.gz
RUN cd vips-8.10.0 && ./configure && make && make install

WORKDIR $GOPATH/src/github.com/shrinkflate/shrinkflate

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY .env ./build/.env

RUN go build -o ./build/shrinkflate .

EXPOSE 4000

CMD ["./build/shrinkflate"]