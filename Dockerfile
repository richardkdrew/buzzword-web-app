FROM scratch

MAINTAINER Richard Drew <richardkdrew@gmail.com>

# copy the files into place
ADD buzzword-web /buzzword-web

# expose the hosting port
ENV HOST_PORT 8080
EXPOSE 8080

# set the run command
CMD ["/buzzword-web"]
