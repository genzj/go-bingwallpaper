FROM go:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/genzj/gobingwallpaper"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/gobingwallpaper/bin

WORKDIR /opt/gobingwallpaper/bin

COPY bin/gobingwallpaper /opt/gobingwallpaper/bin/
RUN chmod +x /opt/gobingwallpaper/bin/gobingwallpaper

CMD /opt/gobingwallpaper/bin/gobingwallpaper
