# TODO arch (amd64 and arm64).
# TODO add the gotenberg binary (with correct arch).
# TODO update Tini according to arch.
FROM ubuntu:20.04

# Note: ARG instructions do not create additional layers.
# Instead, next layers will concatenate them.
ARG GOTENBERG_VERSION

LABEL author="Julien Neuhart" \
      description="A Docker-powered stateless API for converting HTML, Markdown and Office documents to PDF." \
      github="https://github.com/gotenberg/gotenberg" \
      support="support@gotenberg.dev" \
      version="$GOTENBER_VERSION" \
      website="https://gotenberg.dev"

# Improve fonts subpixel hinting and smoothing.
# Credits:
# https://github.com/arachnys/athenapdf/issues/69
# https://github.com/arachnys/athenapdf/commit/ba25a8d80a25d08d58865519c4cd8756dc9a336d
COPY build/package/fonts.conf /etc/fonts/conf.d/100-gotenberg.conf

# Setup the Docker image.
ARG GOTENBERG_USER_GID
ARG GOTENBERG_USER_UID
ARG TINI_VERSION

RUN \
    # Create a non-root user.
    # All processes in the Docker container will run with this dedicated user.
    groupadd --gid ${GOTENBERG_USER_GID} gotenberg &&\
    useradd --uid ${GOTENBERG_USER_UID} --gid gotenberg --shell /bin/bash --home /home/gotenberg --no-create-home gotenberg &&\
    mkdir /home/gotenberg &&\
    chown gotenberg: /home/gotenberg &&\
    # Install dependencies required for the next instructions.
    # Note: gnupg required by wget.
    apt-get update -qq &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends curl wget gnupg &&\
    # Install fonts.
    # Credits:
    # https://github.com/arachnys/athenapdf/blob/master/cli/Dockerfile
    # https://help.accusoft.com/PrizmDoc/v12.1/HTML/Installing_Asian_Fonts_on_Ubuntu_and_Debian.html
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends \
    ttf-mscorefonts-installer \
    culmus \
    fonts-beng \
    fonts-hosny-amiri \
    fonts-lklug-sinhala \
    fonts-lohit-guru \
    fonts-lohit-knda \
    fonts-samyak-gujr \
    fonts-samyak-mlym \
    fonts-samyak-taml \
    fonts-sarai \
    fonts-sil-abyssinica \
    fonts-sil-padauk \
    fonts-telu \
    fonts-thai-tlwg \
    fonts-liberation \
    ttf-wqy-zenhei \
    fonts-arphic-uming \
    fonts-ipafont-mincho \
    fonts-ipafont-gothic \
    fonts-unfonts-core &&\
    # Add Google Chrome repository.
    # Chromium would be a great option, but it can only be installed with snap currently (on Ubuntu).
    # Alas, snap requires systemd, which does not work with Docker.
    wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add - &&\
    echo "deb http://dl.google.com/linux/chrome/deb/ stable main" | tee /etc/apt/sources.list.d/google-chrome.list &&\
    apt-get update -qq &&\
    # Install Ghostscript, Google Chrome and LibreOffice.
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends ghostscript google-chrome-stable libreoffice &&\
    # Download unoconv (Python script).
    curl -Ls https://raw.githubusercontent.com/gotenberg/unoconv/master/unoconv -o /usr/bin/unoconv &&\
    chmod +x /usr/bin/unoconv &&\
    # unoconv will look for the Python binary, which has to be at version 3.
    ln -s /usr/bin/python3 /usr/bin/python &&\
    # Add Tini, a helper for reaping zombie processes.
    curl -Ls https://github.com/krallin/tini/releases/download/v${TINI_VERSION}/tini-static -o /usr/bin/tini &&\
    chmod +x /usr/bin/tini &&\
    # For security reasons, the non-root user gotenberg do not own the Tini binary.
    # However, some providers like Cloud Run from Google Cloud cannot start the Docker image in that case.
    # Uncomment the next instruction if you encounter this issue and rebuild the Docker image.
    # See https://github.com/gotenberg/gotenberg/issues/90#issuecomment-543551353.
    # TODO: verify it is still the case.
    # chown gotenberg: /usr/bin/tini &&\
    # Cleanup.
    # Note: the Ubuntu image does automatically a clean after each install thanks to a hook.
    # Therefore, there is no need for apt-get clean.
    # See https://stackoverflow.com/a/24417119/3248473.
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* &&\
    # Print versions of main dependencies.
    ghostscript --version &&\
    google-chrome-stable --version &&\
    libreoffice --version &&\
    unoconv --version

USER gotenberg
WORKDIR /home/gotenberg

# Default API port.
EXPOSE 3000

# Default web application port.
EXPOSE 80

ENTRYPOINT [ "/usr/bin/tini", "--" ]
CMD [ "gotenberg" ]