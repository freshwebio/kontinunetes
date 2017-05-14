FROM busybox
ADD kontinunetes /
EXPOSE 3229
CMD ["/kontinunetes"]
