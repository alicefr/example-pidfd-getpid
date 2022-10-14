FROM fedora:36 

RUN dnf install gcc git make -y
RUN git clone git://pflaume.badlangua.ge/pflaume && \
    cd pflaume && \
    make
