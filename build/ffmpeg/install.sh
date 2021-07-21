cd third_party/nv_codec_headers && sudo make install
sudo apt-get install \
  build-essential \
  yasm \
  cmake \
  libtool \
  libc6 \
  libc6-dev \
  unzip \
  wget \
  libnuma1 \
  libnuma-dev
cd -
cd third_party/ffmpeg && ./configure --enable-static --disable-shared --enable-gpl --disable-doc && make -j 4