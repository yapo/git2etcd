FROM alpine:3.3
MAINTAINER jp@roemer.im

# Install Gosu to /usr/local/bin/gosu
ADD https://github.com/tianon/gosu/releases/download/1.7/gosu-amd64 /usr/local/sbin/gosu

# Install runtime dependencies & create runtime user
RUN chmod +x /usr/local/sbin/gosu \
 && echo "@testing http://dl-4.alpinelinux.org/alpine/edge/testing" | tee -a /etc/apk/repositories \
 && apk --no-cache --no-progress add ca-certificates git libgit2@testing \
 && adduser -D app -h /data -s /bin/sh

# Copy source code to the container & build it
COPY . /app
WORKDIR /app
RUN ./docker/build.sh

# NSSwitch configuration file
COPY docker/nsswitch.conf /etc/nsswitch.conf

# App configuration
ENV G2E_REPO_PATH "/data/repo"

# Container configuration
VOLUME ["/data"]
EXPOSE 4242
CMD ["/usr/local/sbin/gosu", "app", "/app/git2etcd"]
