FROM tensorflow/serving

ADD ./goldfish /mnt/export/goldfish
RUN ls -lsa /mnt/export

ENTRYPOINT ["/usr/bin/tensorflow_model_server"]