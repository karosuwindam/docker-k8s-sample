FROM multiarch/ubuntu-core:armhf-bionic
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y gnuplot &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
# RUN gnuplot set terminal png &&\
#     gnuplot set output "test.png" 