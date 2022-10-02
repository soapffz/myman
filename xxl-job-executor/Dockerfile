FROM debian:stable-slim

LABEL maintainer="soapffz"

# Environment Variables
ENV HOME /root
ENV DEBIAN_FRONTEND=noninteractive

# Working Directory
WORKDIR /root
RUN mkdir -p ${HOME}/toolkit && \
    mkdir -p ${HOME}/wordlists

# timezone
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ >/etc/timezone && \
    dpkg-reconfigure --frontend noninteractive tzdata

# Install Essentials
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    wget \
    iputils-ping\
    awscli \
    tzdata \
    whois \
    perl \
    dnsutils \
    nano\
    openjdk-11-jdk\
    # dirb
    dirb \
    # dnsenum
    cpanminus \
    # massdns
    libldns-dev \
    libxml2 \
    libxml2-dev \
    libxslt1-dev \
    ruby-dev \
    libgmp-dev \
    zlib1g-dev \
    # joomscan
    libwww-perl \
    # hydra
    hydra \
    # dnsrecon
    dnsrecon \
    && rm -rf /var/lib/apt/lists/*

# Install tools by f8x
RUN wget -O f8x https://raw.githubusercontent.com/ffffffff0x/f8x/main/f8x
RUN chmod +x ./f8x
RUN ./f8x -b -zsh -ka -kb -py38 -py2 -pip2-f -chromium -frp -nn -wpscan -cs

# install go of latest version
RUN sed -i 's/go1.18.6/go1.19.1/g'
RUN head -n 50 f8x
RUN ./f8x -go

# Intstall msf
RUN apt-get install gnupg -y
RUN echo 'deb http://apt.metasploit.com/ lucid main' > /etc/apt/sources.list.d/metasploit-framework.list
RUN wget http://apt.metasploit.com/metasploit-framework.gpg.key
RUN gpg --no-default-keyring --keyring ./metasploit-framework_keyring.gpg --import metasploit-framework.gpg.key
RUN gpg --no-default-keyring --keyring ./metasploit-framework_keyring.gpg --export > ./metasploit-framework.gpg
RUN mv ./metasploit-framework.gpg /etc/apt/trusted.gpg.d/
RUN apt-get update
RUN apt-get -y install metasploit-framework

# python
RUN python -m pip install --upgrade setuptools && python3 -m pip install --upgrade setuptools
RUN python -m pip install --upgrade pip && python3 -m pip install --upgrade pip
RUN apt-get update && \
    apt-get install -y \
    python3-dev
RUN pip3 install wheel requests && pip install wheel requests
RUN pip3 install lxml tldextract flask simplejson

# go tools  
RUN go install -v github.com/tomnomnom/assetfinder@latest &>/dev/null
RUN go install -v github.com/tomnomnom/qsreplace@latest &>/dev/null
RUN go install -v github.com/lc/gau@latest &>/dev/null
RUN go install -v github.com/tomnomnom/unfurl@latest &>/dev/null
RUN go install -v github.com/tomnomnom/anew@latest &>/dev/null
RUN go install -v github.com/projectdiscovery/chaos-client/cmd/chaos@latest &>/dev/null
RUN go install -v github.com/tomnomnom/gf@latest && mkdir ~/.gf && cd ~/.gf && git clone https://github.com/tomnomnom/gf/ && cp gf/examples/*.json ./ && git clone https://github.com/1ndianl33t/Gf-Patterns && cp Gf-Patterns/*.json ./

# nuclei
RUN nuclei -ut &>/dev/null

# Sublist3r
RUN cd ${HOME}/toolkit/ && git clone https://github.com/aboul3la/Sublist3r && cd Sublist3r/ && pip2 install -r requirements.txt

# # findomain
# # RUN curl -s https://api.github.com/repos/Findomain/Findomain/releases/latest | grep -E 'browser_download_url' | grep linux | cut -d '"' -f 4 | grep -v "i386" | wget -qi -
# # RUN if [ $? -eq 0 ]; then chmod +x findomain-linux && rm -rf findomain-linux.zip; fi

# seclists
RUN cd ${HOME}/wordlists/ && \
    git clone --depth 1 https://github.com/danielmiessler/SecLists.git

# soapffz's tools
RUN cd ${HOME}/toolkit/ \
    && git clone https://github.com/soapffz/myman \
    && cp -rf myman/xxl-job-executor/toolkit/* ./

# xxl-job-executor things
ENV PARAMS=""
ADD ./xxl-job-executor-sample-springboot-*.jar /app.jar

ENTRYPOINT ["sh","-c","java -jar $JAVA_OPTS /app.jar $PARAMS"]
