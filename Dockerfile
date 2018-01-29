# This is a docker file for building docker image, see https://docs.docker.com
# Created by Mingzhi Lin (mingzhi9@gmail.com).

FROM ubuntu:17.10

RUN apt-get update
RUN apt-get install -y golang git smalt samtools sra-toolkit parallel

RUN go get -u github.com/kussell-lab/go-misc/cmd/GenomicConsensus
RUN go get -u github.com/kussell-lab/go-misc/cmd/CollectGeneAlignments

RUN git clone https://github.com/kussell-lab/ReferenceAlignmentGenerator.git /opt/ReferenceAlignmentGenerator
ENV PATH="/opt/ReferenceAlignmentGenerator:${PATH}"

RUN /bin/bash -c "echo 'will cite' | parallel --citation"
