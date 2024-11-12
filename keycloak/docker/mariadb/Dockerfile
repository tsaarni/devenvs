FROM mariadb:11.5
RUN apt-get update && \
    apt-get install mariadb-plugin-hashicorp-key-management -y && \
    rm -rf /var/lib/apt/lists/*
