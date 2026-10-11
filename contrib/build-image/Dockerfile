FROM docker.io/node:24.21.0-trixie-slim AS node

FROM docker.io/golang:1.27.2-trixie

# Debian trixie only packages Node 20 (EOL), so Node comes from the official
# image instead; both are trixie-based, so the binary matches this glibc.
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/npm
RUN ln -s ../lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
    ln -s ../lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx

RUN npm install -g cdk8s-cli@2.207.65
