FROM ubuntu:12.10

# Install sshd
RUN apt-get update
RUN apt-get install -y openssh-server golang git

ENV GOPATH /root/go
ENV PATH $GOPATH/bin:$PATH

RUN mkdir -p /var/run/sshd
RUN mkdir -p /root/.ssh
RUN mkdir -p /root/go

RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA4/d5Zz3bfl1rx2w7CsMYIkTIhUiXER1QZUl5fO/QsN+XR2YGzxeCvdm/wT+fjOCMpTa9uCslhRCAOleGlZY5bm9ZX88+9hGdWZO98y5JuvklyRUowaNm6vQUylR14b9N2yVdBbCAouLN28h5QRVPZtgH3x60VPQWqITvzGGY3I5sDRw5RxIZPMQT5WFLNhu28ciHGG46xzKPCVEFh+Hf2lC6nbp5jW6S8fhn+1iJ7KdeWpwYTMiPFsZu/CcNb7a7wwg3ahp3XcLrB2oaw+aCEF0svAaBN5+DM38MYTQn7rRDroLgAF00Yzjb61sBrAG2YnlCQM+grGpeg+B6DqeKeQ== woid@Alp.local" > /root/.ssh/authorized_keys

RUN mkdir -p /root/go/src/github.com/darwin/scraper
ADD . /root/go/src/github.com/darwin/scraper
WORKDIR /root/go/src/github.com/darwin/scraper
RUN go get .
RUN go install .

EXPOSE 22

#CMD /usr/sbin/sshd -o PasswordAuthentication=no && scraper --workspace=/var/scraper

# hint:
# docker build -t "scraper" .
# docker run -d -p 9922:22 -v /var/scraper:/var/scraper scraper /bin/sh -c "/usr/sbin/sshd -o PasswordAuthentication=no && scraper --workspace=/var/scraper --urls=someurls,coma,separated"