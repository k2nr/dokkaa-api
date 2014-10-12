FROM golang:1.3.3-onbuild
MAINTAINER Kazunori Kajihiro <likerichie@gmail.com> (@k2nr)

EXPOSE 80
ENV PORT 80
ENTRYPOINT ["app"]
