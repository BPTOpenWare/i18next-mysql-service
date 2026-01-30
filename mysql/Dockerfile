ARG MYSQL_VERSION="9.3.0"
FROM mysql:${MYSQL_VERSION}

COPY BPTRESOURCENEWS.sql /docker-entrypoint-initdb.d/BPTRESOURCENEWS.sql


RUN echo "mysqlrun:x:1014:1014::/bin/sh" >> /etc/passwd \
    && echo "mysqlrun:x:1014:mysqlrun" >> /etc/group
    
USER 1014

