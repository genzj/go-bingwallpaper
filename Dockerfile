FROM go:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/genzj/go-bingwallpaper"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/go-bingwallpaper/bin

WORKDIR /opt/go-bingwallpaper/bin

COPY bin/go-bingwallpaper /opt/go-bingwallpaper/bin/
RUN chmod +x /opt/go-bingwallpaper/bin/go-bingwallpaper

CMD /opt/go-bingwallpaper/bin/go-bingwallpaper
